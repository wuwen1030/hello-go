package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type APITestSuite struct {
	suite.Suite
	baseURL string
}

func TestAPISuite(t *testing.T) {
	suite.Run(t, new(APITestSuite))
}

func (s *APITestSuite) SetupSuite() {
	s.baseURL = "http://localhost:8080/api/v1"
	// 等待服务器启动
	time.Sleep(2 * time.Second)
}

func (s *APITestSuite) TestHealthCheck() {
	resp := s.makeRequest("GET", "/health", nil)
	defer resp.Body.Close()

	var result map[string]interface{}
	err := json.NewDecoder(resp.Body).Decode(&result)
	assert.NoError(s.T(), err)

	assert.Equal(s.T(), float64(0), result["code"])
	assert.Equal(s.T(), "success", result["message"])
	assert.NotNil(s.T(), result["data"].(map[string]interface{})["time"])
}

func (s *APITestSuite) TestArticleCRUD() {
	// 1. 创建文章
	article := s.createArticle()
	articleID := int(article["id"].(float64))

	// 2. 获取文章
	s.getArticle(articleID)

	// 3. 更新文章
	s.updateArticle(articleID)

	// 4. 获取文章列表
	s.listArticles()

	// 5. 删除文章
	s.deleteArticle(articleID)

	// 6. 验证删除
	s.verifyArticleDeleted(articleID)
}

func (s *APITestSuite) createArticle() map[string]interface{} {
	body := map[string]interface{}{
		"title":   "测试文章",
		"content": "这是一篇测试文章的内容",
	}

	resp := s.makeRequest("POST", "/articles", body)
	defer resp.Body.Close()

	var result map[string]interface{}
	err := json.NewDecoder(resp.Body).Decode(&result)
	assert.NoError(s.T(), err)

	assert.Equal(s.T(), float64(0), result["code"])
	data := result["data"].(map[string]interface{})
	assert.NotNil(s.T(), data["id"])
	assert.Equal(s.T(), body["title"], data["title"])
	assert.Equal(s.T(), body["content"], data["content"])

	return data
}

func (s *APITestSuite) getArticle(id int) {
	resp := s.makeRequest("GET", fmt.Sprintf("/articles/%d", id), nil)
	defer resp.Body.Close()

	var result map[string]interface{}
	err := json.NewDecoder(resp.Body).Decode(&result)
	assert.NoError(s.T(), err)

	assert.Equal(s.T(), float64(0), result["code"])
	data := result["data"].(map[string]interface{})
	assert.Equal(s.T(), float64(id), data["id"])
}

func (s *APITestSuite) updateArticle(id int) {
	body := map[string]interface{}{
		"title":   "更新后的标题",
		"content": "更新后的内容",
		"status":  2,
	}

	resp := s.makeRequest("PUT", fmt.Sprintf("/articles/%d", id), body)
	defer resp.Body.Close()

	var result map[string]interface{}
	err := json.NewDecoder(resp.Body).Decode(&result)
	assert.NoError(s.T(), err)

	assert.Equal(s.T(), float64(0), result["code"])
	data := result["data"].(map[string]interface{})
	assert.Equal(s.T(), body["title"], data["title"])
	assert.Equal(s.T(), body["content"], data["content"])
	assert.Equal(s.T(), float64(2), data["status"])
}

func (s *APITestSuite) listArticles() {
	resp := s.makeRequest("GET", "/articles?page=1&page_size=10", nil)
	defer resp.Body.Close()

	var result map[string]interface{}
	err := json.NewDecoder(resp.Body).Decode(&result)
	assert.NoError(s.T(), err)

	assert.Equal(s.T(), float64(0), result["code"])
	data := result["data"].(map[string]interface{})
	assert.NotNil(s.T(), data["items"])
	assert.NotNil(s.T(), data["total"])
}

func (s *APITestSuite) deleteArticle(id int) {
	resp := s.makeRequest("DELETE", fmt.Sprintf("/articles/%d", id), nil)
	defer resp.Body.Close()

	var result map[string]interface{}
	err := json.NewDecoder(resp.Body).Decode(&result)
	assert.NoError(s.T(), err)

	assert.Equal(s.T(), float64(0), result["code"])
}

func (s *APITestSuite) verifyArticleDeleted(id int) {
	resp := s.makeRequest("GET", fmt.Sprintf("/articles/%d", id), nil)
	defer resp.Body.Close()

	var result map[string]interface{}
	err := json.NewDecoder(resp.Body).Decode(&result)
	assert.NoError(s.T(), err)

	assert.NotEqual(s.T(), float64(0), result["code"])
}

func (s *APITestSuite) makeRequest(method, path string, body interface{}) *http.Response {
	var req *http.Request
	var err error

	if body != nil {
		jsonBody, _ := json.Marshal(body)
		req, err = http.NewRequest(method, s.baseURL+path, bytes.NewBuffer(jsonBody))
	} else {
		req, err = http.NewRequest(method, s.baseURL+path, nil)
	}

	assert.NoError(s.T(), err)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	assert.NoError(s.T(), err)

	return resp
}
