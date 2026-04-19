package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

// Config 全局配置结构
type Config struct {
	OpenAI    OpenAIConfig    `yaml:"openai"`
	Qdrant    QdrantConfig    `yaml:"qdrant"`
	Server    ServerConfig    `yaml:"server"`
	Embedding EmbeddingConfig `yaml:"embedding"`
}

// OpenAIConfig OpenAI/DashScope API 配置
type OpenAIConfig struct {
	APIKey  string `yaml:"api_key"`
	BaseURL string `yaml:"base_url"`
	Model   string `yaml:"model"`
}

// QdrantConfig Qdrant 向量数据库配置
type QdrantConfig struct {
	Host    string `yaml:"host"`
	Port    int    `yaml:"port"`
	APIKey  string `yaml:"api_key"`
	UseHTTP bool   `yaml:"use_http"`
}

// ServerConfig 服务配置
type ServerConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
	Mode string `yaml:"mode"` // debug, release, test
}

// EmbeddingConfig Embedding 服务配置
type EmbeddingConfig struct {
	Model string `yaml:"model"` // text-embedding-v3 等
}

// Load 从文件加载配置
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	// 环境变量覆盖
	if key := os.Getenv("OPENAI_API_KEY"); key != "" {
		cfg.OpenAI.APIKey = key
	}
	if key := os.Getenv("QDRANT_API_KEY"); key != "" {
		cfg.Qdrant.APIKey = key
	}

	return &cfg, nil
}

// LoadDefault 尝试加载默认配置文件
func LoadDefault() (*Config, error) {
	paths := []string{
		"config.yaml",
		"configs/config.yaml",
		"../config.yaml",
		"../../configs/config.yaml",
	}

	for _, p := range paths {
		if cfg, err := Load(p); err == nil {
			return cfg, nil
		}
	}

	// 返回默认配置
	return &Config{
		OpenAI: OpenAIConfig{
			Model: "qwen-plus",
		},
		Qdrant: QdrantConfig{
			Host: "localhost",
			Port: 6333,
		},
		Server: ServerConfig{
			Host: "0.0.0.0",
			Port: 8080,
			Mode: "debug",
		},
		Embedding: EmbeddingConfig{
			Model: "text-embedding-v3",
		},
	}, nil
}