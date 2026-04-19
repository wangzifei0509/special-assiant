package app

import (
	"github.com/gin-gonic/gin"

	"github.com/chunfengshili/sa/services/gateway/handler"
	"github.com/chunfengshili/sa/services/gateway/middleware"
	"github.com/chunfengshili/sa/pkg/llm"
)

// Router 网关路由
type Router struct {
	engine *gin.Engine
}

// NewRouter 创建路由
func NewRouter(client *llm.Client, mode string) *Router {
	gin.SetMode(mode)
	engine := gin.New()

	// 中间件
	engine.Use(middleware.Cors())
	engine.Use(middleware.Logger())
	engine.Use(gin.Recovery())

	// 处理器
	healthHandler := handler.NewHealthHandler()
	chatHandler := handler.NewChatHandler(client)

	// 公开路由
	engine.GET("/health", healthHandler.Check)

	// API v1
	v1 := engine.Group("/api/v1")
	{
		v1.POST("/chat", chatHandler.Chat)
		v1.POST("/chat/stream", chatHandler.ChatStream)
	}

	// 需要认证的路由
	auth := engine.Group("/api/v1")
	auth.Use(middleware.Auth())
	{
		// TODO: 添加需要认证的路由
	}

	return &Router{engine: engine}
}

// Run 启动服务
func (r *Router) Run(addr string) error {
	return r.engine.Run(addr)
}