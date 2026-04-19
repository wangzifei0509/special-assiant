package main

import (
	"fmt"
	"os"

	"github.com/chunfengshili/sa/services/knowledge/app"
	"github.com/chunfengshili/sa/services/knowledge/embedding"
	"github.com/chunfengshili/sa/services/knowledge/qdrant"
	"github.com/chunfengshili/sa/pkg/config"
)

func main() {
	// 加载配置
	cfg, err := config.LoadDefault()
	if err != nil {
		fmt.Fprintf(os.Stderr, "加载配置失败: %v\n", err)
		os.Exit(1)
	}

	// 检查 API Key
	if cfg.OpenAI.APIKey == "" {
		fmt.Fprintln(os.Stderr, "请设置 OPENAI_API_KEY 环境变量或在配置文件中配置 api_key")
		os.Exit(1)
	}

	// 创建 Embedding 服务
	embeddingSvc := embedding.NewService(&embedding.Config{
		APIKey:  cfg.OpenAI.APIKey,
		BaseURL: cfg.OpenAI.BaseURL,
		Model:   cfg.Embedding.Model,
	})

	// 创建 Qdrant 客户端
	qdrantClient, err := qdrant.NewClient(&qdrant.Config{
		Host:    cfg.Qdrant.Host,
		Port:    cfg.Qdrant.Port,
		APIKey:  cfg.Qdrant.APIKey,
		UseHTTP: cfg.Qdrant.UseHTTP,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "连接 Qdrant 失败: %v\n", err)
		os.Exit(1)
	}
	defer qdrantClient.Close()

	// 默认集合名
	collection := "knowledge_base"

	// 创建路由
	router := app.NewRouter(embeddingSvc, qdrantClient, collection, cfg.Server.Mode)

	// 启动服务
	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	fmt.Printf("🚀 Knowledge 服务启动于 %s\n", addr)
	if err := router.Run(addr); err != nil {
		fmt.Fprintf(os.Stderr, "服务启动失败: %v\n", err)
		os.Exit(1)
	}
}