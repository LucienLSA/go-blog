package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	redisCache "github.com/LucienLSA/go-blog/repository/db/dao/redis"
	"github.com/LucienLSA/go-blog/types"
	"go.uber.org/zap"
)

// GitHubService GitHub服务
type GitHubService struct{}

// GetGitHubTrendingData 获取GitHub热点数据
func (s *GitHubService) GetGitHubTrendingData(ctx context.Context, language, since string) (*types.GitHubTrendingResponse, error) {
	// 首先尝试从Redis获取数据
	data, err := redisCache.GetGitHubTrendingData(language, since)
	if err == nil && len(data) > 0 {
		// 数据存在且未过期，直接返回
		return &types.GitHubTrendingResponse{
			Language: language,
			Since:    since,
			Data:     data,
			Total:    len(data),
		}, nil
	}

	// 如果Redis中没有数据或已过期，从GitHub API获取
	zap.L().Info("fetching github trending data from API",
		zap.String("language", language),
		zap.String("since", since))

	trendingData, err := s.fetchFromGitHubAPI(language, since)
	if err != nil {
		zap.L().Error("fetch github trending data failed", zap.Error(err))
		return nil, err
	}

	// 保存到Redis
	err = redisCache.SaveGitHubTrendingData(language, since, trendingData)
	if err != nil {
		zap.L().Error("save github trending data to redis failed", zap.Error(err))
		// 即使保存失败也返回数据
	}

	return &types.GitHubTrendingResponse{
		Language: language,
		Since:    since,
		Data:     trendingData,
		Total:    len(trendingData),
	}, nil
}

// GetAllGitHubTrendingData 获取所有GitHub热点数据
func (s *GitHubService) GetAllGitHubTrendingData(ctx context.Context) (map[string][]types.GitHubTrendingData, error) {
	return redisCache.GetAllGitHubTrendingData()
}

// fetchFromGitHubAPI 从GitHub API获取热点数据
func (s *GitHubService) fetchFromGitHubAPI(language, since string) ([]types.GitHubTrendingData, error) {
	// 尝试多个API端点
	apis := []string{
		"https://github-trending-api.vercel.app/repositories",
		"https://github-trending-api.now.sh/repositories",
		"https://gh-trending-api.herokuapp.com/repositories",
	}

	var lastErr error
	for _, baseURL := range apis {
		// 构建API URL
		url := baseURL
		params := []string{}

		if language != "" {
			params = append(params, "language="+language)
		}
		if since != "" {
			params = append(params, "since="+since)
		}

		if len(params) > 0 {
			url += "?" + strings.Join(params, "&")
		}

		zap.L().Info("trying github trending API", zap.String("url", url))

		// 创建HTTP客户端
		client := &http.Client{
			Timeout: 10 * time.Second, // 减少超时时间
		}

		// 发送请求
		resp, err := client.Get(url)
		if err != nil {
			lastErr = fmt.Errorf("failed to fetch from %s: %w", baseURL, err)
			zap.L().Warn("API request failed, trying next endpoint",
				zap.String("url", url),
				zap.Error(err))
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			lastErr = fmt.Errorf("API %s returned status: %d", baseURL, resp.StatusCode)
			zap.L().Warn("API returned non-OK status",
				zap.String("url", url),
				zap.Int("status", resp.StatusCode))
			continue
		}

		// 读取响应体
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			lastErr = fmt.Errorf("failed to read response from %s: %w", baseURL, err)
			continue
		}

		// 解析JSON响应
		var rawData []map[string]interface{}
		err = json.Unmarshal(body, &rawData)
		if err != nil {
			lastErr = fmt.Errorf("failed to unmarshal JSON from %s: %w", baseURL, err)
			continue
		}

		// 转换为我们的数据结构
		var trendingData []types.GitHubTrendingData
		for _, item := range rawData {
			// 提取ID
			id, _ := strconv.ParseInt(fmt.Sprintf("%v", item["id"]), 10, 64)

			// 提取stars和forks
			stars := 0
			if starsStr, ok := item["stars"].(string); ok {
				stars, _ = strconv.Atoi(strings.ReplaceAll(starsStr, ",", ""))
			} else if starsInt, ok := item["stars"].(float64); ok {
				stars = int(starsInt)
			}

			forks := 0
			if forksStr, ok := item["forks"].(string); ok {
				forks, _ = strconv.Atoi(strings.ReplaceAll(forksStr, ",", ""))
			} else if forksInt, ok := item["forks"].(float64); ok {
				forks = int(forksInt)
			}

			// 构建数据
			data := types.GitHubTrendingData{
				ID:          id,
				Name:        fmt.Sprintf("%v", item["name"]),
				FullName:    fmt.Sprintf("%v", item["fullName"]),
				Description: fmt.Sprintf("%v", item["description"]),
				Language:    fmt.Sprintf("%v", item["language"]),
				Stars:       stars,
				Forks:       forks,
				URL:         fmt.Sprintf("%v", item["url"]),
				CreatedAt:   fmt.Sprintf("%v", item["createdAt"]),
				UpdatedAt:   fmt.Sprintf("%v", item["updatedAt"]),
			}

			trendingData = append(trendingData, data)
		}

		zap.L().Info("successfully fetched github trending data",
			zap.String("language", language),
			zap.String("since", since),
			zap.Int("count", len(trendingData)),
			zap.String("api", baseURL))

		return trendingData, nil
	}

	// 如果所有API都失败了，返回模拟数据
	zap.L().Warn("all GitHub APIs failed, returning mock data", zap.Error(lastErr))
	return s.getMockData(language, since), nil
}

// getMockData 获取模拟数据（当API不可用时使用）
func (s *GitHubService) getMockData(language, since string) []types.GitHubTrendingData {
	mockData := []types.GitHubTrendingData{
		{
			ID:          1,
			Name:        "awesome-go-project",
			FullName:    "developer/awesome-go-project",
			Description: "这是一个用Go语言开发的优秀项目",
			Language:    "Go",
			Stars:       1500,
			Forks:       120,
			URL:         "https://github.com/developer/awesome-go-project",
			CreatedAt:   "2024-01-01T00:00:00Z",
			UpdatedAt:   "2024-01-15T00:00:00Z",
		},
		{
			ID:          2,
			Name:        "go-web-framework",
			FullName:    "team/go-web-framework",
			Description: "轻量级的Go Web框架",
			Language:    "Go",
			Stars:       800,
			Forks:       50,
			URL:         "https://github.com/team/go-web-framework",
			CreatedAt:   "2024-01-05T00:00:00Z",
			UpdatedAt:   "2024-01-20T00:00:00Z",
		},
		{
			ID:          3,
			Name:        "go-microservice",
			FullName:    "architect/go-microservice",
			Description: "基于Go的微服务架构示例",
			Language:    "Go",
			Stars:       2000,
			Forks:       300,
			URL:         "https://github.com/architect/go-microservice",
			CreatedAt:   "2024-01-10T00:00:00Z",
			UpdatedAt:   "2024-01-25T00:00:00Z",
		},
	}

	zap.L().Info("returning mock data",
		zap.String("language", language),
		zap.String("since", since),
		zap.Int("count", len(mockData)))

	return mockData
}

// RefreshGitHubTrendingData 刷新GitHub热点数据（定时任务使用）
func (s *GitHubService) RefreshGitHubTrendingData() {
	zap.L().Info("starting github trending data refresh task")

	// 定义要获取的语言和时间范围
	languages := []string{"", "javascript", "python", "go", "java", "rust"}
	sinceOptions := []string{"daily", "weekly", "monthly"}

	for _, language := range languages {
		for _, since := range sinceOptions {
			// 检查数据是否过期
			if redisCache.IsGitHubTrendingDataExpired(language, since) {
				zap.L().Info("refreshing expired github trending data",
					zap.String("language", language),
					zap.String("since", since))

				// 获取新数据
				trendingData, err := s.fetchFromGitHubAPI(language, since)
				if err != nil {
					zap.L().Error("failed to refresh github trending data",
						zap.String("language", language),
						zap.String("since", since),
						zap.Error(err))
					continue
				}

				// 保存到Redis
				err = redisCache.SaveGitHubTrendingData(language, since, trendingData)
				if err != nil {
					zap.L().Error("failed to save refreshed github trending data",
						zap.String("language", language),
						zap.String("since", since),
						zap.Error(err))
				}
			}
		}
	}

	zap.L().Info("github trending data refresh task completed")
}

// 全局GitHub服务实例
var githubService *GitHubService

// GetGitHubService 获取GitHub服务实例
func GetGitHubService() *GitHubService {
	if githubService == nil {
		githubService = &GitHubService{}
	}
	return githubService
}
