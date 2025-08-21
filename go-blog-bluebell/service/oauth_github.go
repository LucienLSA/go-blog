package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/LucienLSA/go-blog/pkg/jwt"
	"github.com/LucienLSA/go-blog/pkg/snowflake"
	"github.com/LucienLSA/go-blog/repository/db/dao/mysql"
	redisCache "github.com/LucienLSA/go-blog/repository/db/dao/redis"
	"github.com/LucienLSA/go-blog/repository/db/models"
	"github.com/LucienLSA/go-blog/settings"
	"go.uber.org/zap"
)

type OAuthGitHubService struct{}

var oauthGitHubService *OAuthGitHubService

func GetOAuthGitHubService() *OAuthGitHubService {
	if oauthGitHubService == nil {
		oauthGitHubService = &OAuthGitHubService{}
	}
	return oauthGitHubService
}

// GenerateState creates a random state and stores it in Redis
func (s *OAuthGitHubService) GenerateState(ctx context.Context, ttl time.Duration) (string, error) {
	// 生成 16 字节随机字符串（通过rand.Read），经hex.EncodeToString转为 32 位字符串作为state
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	state := hex.EncodeToString(b)
	// 通过redisCache.SaveOAuthState存储到 Redis
	if err := redisCache.SaveOAuthState(state, ttl); err != nil {
		return "", err
	}
	return state, nil
}

// BuildAuthorizeURL builds GitHub authorize URL with configured clientID and redirectURL
func (s *OAuthGitHubService) BuildAuthorizeURL(state string) string {
	// 从配置中读取 GitHub 的ClientID、RedirectURL、Scopes（权限范围，如user:email获取邮箱）
	// ，拼接为 GitHub 授权页面的 URL。
	cfg := settings.Conf.GitHubConfig
	scopes := strings.Join(cfg.Scopes, " ")
	q := url.Values{}
	q.Set("client_id", cfg.ClientID)
	q.Set("redirect_uri", cfg.RedirectURL)
	q.Set("scope", scopes)
	q.Set("state", state)
	q.Set("allow_signup", "true")
	// 用户访问该 URL 后，会跳转至 GitHub 的登录授权页，
	// 确认授权后 GitHub 会携带code和state回调到redirect_uri。
	return "https://github.com/login/oauth/authorize?" + q.Encode()
}

// 处理 GitHub 回调
// HandleCallback exchanges code for token, fetches user, binds/creates local user, and returns JWT and redirect URL
func (s *OAuthGitHubService) HandleCallback(ctx context.Context, code, state string) (token string, redirect string, err error) {
	// 检查回调中的state是否存在于 Redis（且未过期），无效则跳转至错误页，防止非法请求
	valid, err := redisCache.ConsumeOAuthState(state)
	if err != nil || !valid {
		return "", settings.Conf.OAuthRedirectConfig.ErrorRedirect, errors.New("invalid state")
	}
	// 调用exchangeCodeForToken，向 GitHub 的 token 端点发送请求（携带client_id、client_secret、code等），
	// 获取access_token（用于后续调用 GitHub API）
	accessToken, err := s.exchangeCodeForToken(ctx, code)
	if err != nil {
		return "", settings.Conf.OAuthRedirectConfig.ErrorRedirect, err
	}
	// 调用fetchGitHubUser，使用access_token调用 GitHub 的用户 API 和邮箱 API
	// ，获取用户 ID、登录名、头像、主邮箱（优先取已验证的主邮箱）
	ghUser, email, err := s.fetchGitHubUser(ctx, accessToken)
	if err != nil {
		return "", settings.Conf.OAuthRedirectConfig.ErrorRedirect, err
	}

	//先通过邮箱查询本地用户（若 GitHub 返回了邮箱），
	// 存在则直接绑定（更新 GitHub 相关字段：github_id、github_login等）。
	userDao := mysql.NewUserDao(ctx)
	var user *models.User
	// Prefer by email first if present
	if email != "" {
		if existed, exist, _ := userDao.ExistUserEmail(email); exist {
			user = existed
		}
	}
	if user == nil {
		// Try by github login (custom field)
		// No direct DAO; we fall back to create new user
	}

	if user == nil {
		// 若不存在，通过 GitHub 登录名创建新用户，调用uniqueUsername
		// 确保用户名唯一（避免重复），并插入数据库。
		newUser := &models.User{
			UserID:        snowflake.GenID(),
			UserName:      s.uniqueUsername(ctx, ghUser.Login),
			Email:         email,
			Avatar:        ghUser.AvatarURL,
			GitHubID:      ghUser.ID,
			GitHubLogin:   ghUser.Login,
			GitHubAvatar:  ghUser.AvatarURL,
			OAuthProvider: "github",
		}
		if err := userDao.InsertUser(newUser); err != nil {
			zap.L().Error("insert user failed", zap.Error(err))
			return "", settings.Conf.OAuthRedirectConfig.ErrorRedirect, err
		}
		user = newUser
	} else {
		// Update binding fields if missing
		updates := map[string]interface{}{
			"github_id":      ghUser.ID,
			"github_login":   ghUser.Login,
			"github_avatar":  ghUser.AvatarURL,
			"oauth_provider": "github",
		}
		_ = userDao.UpdateUser(user.UserID, updates)
	}

	// 通过jwt.GenToken生成系统内的 JWT 令牌，
	// 存储到 Redis（redisCache.StorgeUserToken），最终返回令牌和成功跳转 URL
	jwtToken, jerr := jwt.GenToken(user.UserID, user.UserName)
	if jerr != nil {
		return "", settings.Conf.OAuthRedirectConfig.ErrorRedirect, jerr
	}
	user.Token = jwtToken
	_ = redisCache.StorgeUserToken(jwtToken, user.UserName)

	return jwtToken, settings.Conf.OAuthRedirectConfig.SuccessRedirect, nil
}

// 向 GitHub 的https://github.com/login/oauth/access_token发送 POST 请求，
// 携带授权码code等参数，解析返回的 JSON 获取access_token。
func (s *OAuthGitHubService) exchangeCodeForToken(ctx context.Context, code string) (string, error) {
	cfg := settings.Conf.GitHubConfig
	form := url.Values{}
	form.Set("client_id", cfg.ClientID)
	form.Set("client_secret", cfg.ClientSecret)
	form.Set("code", code)
	form.Set("redirect_uri", cfg.RedirectURL)

	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, "https://github.com/login/oauth/access_token", strings.NewReader(form.Encode()))
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("token exchange failed: %s", string(body))
	}
	var out struct {
		AccessToken string `json:"access_token"`
		Scope       string `json:"scope"`
		TokenType   string `json:"token_type"`
	}
	if err := json.Unmarshal(body, &out); err != nil {
		return "", err
	}
	if out.AccessToken == "" {
		return "", errors.New("empty access token")
	}
	return out.AccessToken, nil
}

type ghUserResponse struct {
	ID        int64  `json:"id"`
	Login     string `json:"login"`
	Name      string `json:"name"`
	AvatarURL string `json:"avatar_url"`
}

type ghEmail struct {
	Email    string `json:"email"`
	Primary  bool   `json:"primary"`
	Verified bool   `json:"verified"`
}

func (s *OAuthGitHubService) fetchGitHubUser(ctx context.Context, accessToken string) (*ghUserResponse, string, error) {
	// 调用 GitHub 用户 API（https://api.github.com/user）获取用户基本信息（ID、登录名、头像）
	userReq, _ := http.NewRequestWithContext(ctx, http.MethodGet, "https://api.github.com/user", nil)
	userReq.Header.Set("Authorization", "Bearer "+accessToken)
	userReq.Header.Set("Accept", "application/vnd.github+json")
	userResp, err := http.DefaultClient.Do(userReq)
	if err != nil {
		return nil, "", err
	}
	defer userResp.Body.Close()
	ub, _ := io.ReadAll(userResp.Body)
	if userResp.StatusCode != http.StatusOK {
		return nil, "", fmt.Errorf("github user api failed: %s", string(ub))
	}
	var gh ghUserResponse
	if err := json.Unmarshal(ub, &gh); err != nil {
		return nil, "", err
	}

	// 调用邮箱 API（https://api.github.com/user/emails）获取邮箱列表，
	// 优先选择 “主要且已验证” 的邮箱
	emailReq, _ := http.NewRequestWithContext(ctx, http.MethodGet, "https://api.github.com/user/emails", nil)
	emailReq.Header.Set("Authorization", "Bearer "+accessToken)
	emailReq.Header.Set("Accept", "application/vnd.github+json")
	emailResp, err := http.DefaultClient.Do(emailReq)
	if err != nil {
		return &gh, "", nil // no hard fail, email optional
	}
	defer emailResp.Body.Close()
	eb, _ := io.ReadAll(emailResp.Body)
	var emails []ghEmail
	_ = json.Unmarshal(eb, &emails)
	primary := ""
	for _, e := range emails {
		if e.Primary && e.Verified {
			primary = e.Email
			break
		}
	}
	return &gh, primary, nil
}

// uniqueUsername ensures no duplicates when creating new users
func (s *OAuthGitHubService) uniqueUsername(ctx context.Context, base string) string {
	// 创建新用户时，确保用户名不重复 —— 以 GitHub 登录名为基础
	// ，若已存在则添加时间戳后缀（最多重试 5 次），最终生成唯一用户名
	if base == "" {
		base = "github_user"
	}
	userDao := mysql.NewUserDao(ctx)
	name := base
	for i := 0; i < 5; i++ {
		_, exist, err := userDao.CheckUserExist(name)
		if err == nil && !exist {
			return name
		}
		name = fmt.Sprintf("%s_%d", base, time.Now().Unix()%100000)
	}
	return fmt.Sprintf("%s_%d", base, time.Now().UnixNano()%1000000)
}
