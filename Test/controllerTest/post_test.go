package controllerTest

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/LucienLSA/go-blog/controller"
	"github.com/LucienLSA/go-blog/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestCreatePostHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	// r := routers.SetupRouter(settings.Conf.Mode)
	r := gin.Default()
	url := "/api/v1/post"
	r.POST(url, controller.CreatePostHandler)
	// body := `{
	// 	"community_id": 1,
	// 	"title": "test",
	// 	"content": "This is a test"
	// }`
	body := `{
		"community_id": 1,
		"title": "test",
	}`
	// 发送测试请求
	req, _ := http.NewRequest(http.MethodPost, url, bytes.NewReader([]byte(body)))
	// 测试响应
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// 判断响应码
	assert.Equal(t, 200, w.Code)

	// 判断响应内容是否完整 body
	res := new(response.ResponseData)
	if err := json.Unmarshal(w.Body.Bytes(), res); err != nil {
		t.Fatalf("json.Unmarshal w.Body failed, err:%v\n", err)
	}
	assert.Equal(t, res.Code, response.CodeInvalidParam)

	// 判断响应内容是不是预期需要登录的错误
	// 方法1： 判断响应的内容是不是包含指定的字符串
	// assert.Contains(t, w.Body.String(), "需要登录")

	// 方法2：将响应的内容反序列化到res，再判断字段与预期是否一致
	// res = new(response.ResponseData)
	// if err := json.Unmarshal(w.Body.Bytes(), res); err != nil {
	// 	t.Fatalf("json.Unmarshal w.Body failed, err:%v\n", err)
	// }
	// assert.Equal(t, int(res.Code), response.CodeNeedLogin)

}
