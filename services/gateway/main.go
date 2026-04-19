package main

import (
	"fmt"
	"os"

	"github.com/chunfengshili/sa/services/gateway/app"
	"github.com/chunfengshili/sa/pkg/config"
	"github.com/chunfengshili/sa/pkg/llm"
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

	// 创建 LLM 客户端
	client := llm.NewClient(&llm.Config{
		APIKey:  cfg.OpenAI.APIKey,
		BaseURL: cfg.OpenAI.BaseURL,
		Model:   cfg.OpenAI.Model,
	})

	// 创建路由
	router := app.NewRouter(client, cfg.Server.Mode)

	// 启动服务
	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	fmt.Printf("🚀 Gateway 服务启动于 %s\n", addr)
	if err := router.Run(addr); err != nil {
		fmt.Fprintf(os.Stderr, "服务启动失败: %v\n", err)
		os.Exit(1)
	}
}