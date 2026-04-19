package llm

import (
	"github.com/sashabaranov/go-openai"
)

// Client LLM 客户端封装
type Client struct {
	*openai.Client
	config *Config
}

// Config LLM 配置
type Config struct {
	APIKey  string
	BaseURL string
	Model   string
}

// NewClient 创建新的 LLM 客户端
func NewClient(cfg *Config) *Client {
	config := openai.DefaultConfig(cfg.APIKey)
	if cfg.BaseURL != "" {
		config.BaseURL = cfg.BaseURL
	}

	return &Client{
		Client: openai.NewClientWithConfig(config),
		config: cfg,
	}
}

// Model 返回配置的模型名称
func (c *Client) Model() string {
	return c.config.Model
}