package qdrant

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Client Qdrant HTTP 客户端
type Client struct {
	baseURL string
	apiKey  string
	client  *http.Client
}

// Config Qdrant 配置
type Config struct {
	Host    string
	Port    int
	APIKey  string
	UseHTTP bool
}

// NewClient 创建 Qdrant 客户端
func NewClient(cfg *Config) (*Client, error) {
	if cfg.Host == "" {
		cfg.Host = "localhost"
	}
	if cfg.Port == 0 {
		cfg.Port = 6333 // HTTP API 默认端口
	}

	return &Client{
		baseURL: fmt.Sprintf("http://%s:%d", cfg.Host, cfg.Port),
		apiKey:  cfg.APIKey,
		client:  &http.Client{},
	}, nil
}

// doRequest 执行 HTTP 请求
func (c *Client) doRequest(ctx context.Context, method, path string, body interface{}) ([]byte, error) {
	var reqBody io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		reqBody = bytes.NewReader(data)
	}

	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, reqBody)
	if err != nil {
		return nil, err
	}

	if c.apiKey != "" {
		req.Header.Set("api-key", c.apiKey)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("Qdrant error: %s", string(data))
	}

	return data, nil
}

// CollectionInfo 集合信息
type CollectionInfo struct {
	Name string `json:"name"`
}

// VectorParams 向量参数
type VectorParams struct {
	Size     uint64 `json:"size"`
	Distance string `json:"distance"` // Cosine, Euclid, Dot
}

// CreateCollectionRequest 创建集合请求
type CreateCollectionRequest struct {
	Vectors VectorParams `json:"vectors"`
}

// CreateCollection 创建集合
func (c *Client) CreateCollection(ctx context.Context, name string, vectorSize int) error {
	body := CreateCollectionRequest{
		Vectors: VectorParams{
			Size:     uint64(vectorSize),
			Distance: "Cosine",
		},
	}
	_, err := c.doRequest(ctx, "PUT", "/collections/"+name, body)
	return err
}

// DeleteCollection 删除集合
func (c *Client) DeleteCollection(ctx context.Context, name string) error {
	_, err := c.doRequest(ctx, "DELETE", "/collections/"+name, nil)
	return err
}

// Point 点结构
type Point struct {
	ID      string                 `json:"id"`
	Vector  []float32              `json:"vector"`
	Payload map[string]interface{} `json:"payload,omitempty"`
}

// UpsertPointsRequest 批量插入请求
type UpsertPointsRequest struct {
	Points []Point `json:"points"`
}

// UpsertVectors 插入/更新向量
func (c *Client) UpsertVectors(ctx context.Context, collection string, points []Point) error {
	body := UpsertPointsRequest{Points: points}
	_, err := c.doRequest(ctx, "PUT", "/collections/"+collection+"/points", body)
	return err
}

// SearchRequest 搜索请求
type SearchRequest struct {
	Vector []float32 `json:"vector"`
	Limit  int       `json:"limit"`
}

// SearchResult 搜索结果
type SearchResult struct {
	ID     string                 `json:"id"`
	Score  float32                `json:"score"`
	Vector []float32              `json:"vector,omitempty"`
	Payload map[string]interface{} `json:"payload,omitempty"`
}

// SearchResponse 搜索响应
type SearchResponse struct {
	Result []SearchResult `json:"result"`
}

// SearchVectors 搜索相似向量
func (c *Client) SearchVectors(ctx context.Context, collection string, vector []float32, limit int) ([]SearchResult, error) {
	body := SearchRequest{
		Vector: vector,
		Limit:  limit,
	}
	data, err := c.doRequest(ctx, "POST", "/collections/"+collection+"/points/search", body)
	if err != nil {
		return nil, err
	}

	var resp SearchResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, err
	}

	return resp.Result, nil
}

// Close 关闭连接
func (c *Client) Close() error {
	return nil
}