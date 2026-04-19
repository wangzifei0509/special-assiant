package embedding

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Service Embedding 服务
type Service struct {
	apiKey  string
	baseURL string
	model   string
	client  *http.Client
}

// Config Embedding 配置
type Config struct {
	APIKey  string
	BaseURL string
	Model   string // 默认 text-embedding-v3
}

// NewService 创建 Embedding 服务
func NewService(cfg *Config) *Service {
	if cfg.Model == "" {
		cfg.Model = "text-embedding-v3"
	}
	if cfg.BaseURL == "" {
		cfg.BaseURL = "https://dashscope.aliyuncs.com/compatible-mode/v1"
	}

	return &Service{
		apiKey:  cfg.APIKey,
		baseURL: cfg.BaseURL,
		model:   cfg.Model,
		client:  &http.Client{},
	}
}

// embeddingRequest OpenAI兼容格式的embedding请求
type embeddingRequest struct {
	Model string   `json:"model"`
	Input []string `json:"input"`
}

// embeddingResponse OpenAI兼容格式的embedding响应
type embeddingResponse struct {
	Object string `json:"object"`
	Data   []struct {
		Object    string    `json:"object"`
		Index     int       `json:"index"`
		Embedding []float32 `json:"embedding"`
	} `json:"data"`
	Model string `json:"model"`
	Error *struct {
		Message string `json:"message"`
		Type    string `json:"type"`
		Code    string `json:"code"`
	} `json:"error,omitempty"`
}

// Embedding 生成文本嵌入向量
func (s *Service) Embedding(ctx context.Context, texts []string) ([][]float32, error) {
	reqBody := embeddingRequest{
		Model: s.model,
		Input: texts,
	}

	data, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", s.baseURL+"/embeddings", bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.apiKey)

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var embResp embeddingResponse
	if err := json.Unmarshal(body, &embResp); err != nil {
		return nil, err
	}

	if embResp.Error != nil {
		return nil, fmt.Errorf("embedding API error: %s", embResp.Error.Message)
	}

	result := make([][]float32, len(embResp.Data))
	for i, d := range embResp.Data {
		result[i] = d.Embedding
	}

	return result, nil
}

// EmbeddingSingle 生成单个文本的嵌入向量
func (s *Service) EmbeddingSingle(ctx context.Context, text string) ([]float32, error) {
	vectors, err := s.Embedding(ctx, []string{text})
	if err != nil {
		return nil, err
	}
	if len(vectors) == 0 {
		return nil, fmt.Errorf("embedding 返回空结果")
	}
	return vectors[0], nil
}