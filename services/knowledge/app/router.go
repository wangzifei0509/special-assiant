package app

import (
	"github.com/gin-gonic/gin"

	"github.com/chunfengshili/sa/services/knowledge/embedding"
	"github.com/chunfengshili/sa/services/knowledge/handler"
	"github.com/chunfengshili/sa/services/knowledge/qdrant"
)

// Router 知识库服务路由
type Router struct {
	engine *gin.Engine
}

// NewRouter 创建路由
func NewRouter(emb *embedding.Service, qd *qdrant.Client, collection, mode string) *Router {
	gin.SetMode(mode)
	engine := gin.New()

	// 中间件
	engine.Use(gin.Logger())
	engine.Use(gin.Recovery())

	// CORS
	engine.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// 处理器
	knowledgeHandler := handler.NewKnowledgeHandler(emb, qd, collection)

	// 健康检查
	engine.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// API v1
	v1 := engine.Group("/api/v1")
	{
		// 知识库管理
		v1.POST("/collections", knowledgeHandler.CreateCollection)

		// 文档操作
		v1.POST("/documents", knowledgeHandler.AddDocument)
		v1.POST("/documents/search", knowledgeHandler.Search)
	}

	return &Router{engine: engine}
}

// Run 启动服务
func (r *Router) Run(addr string) error {
	return r.engine.Run(addr)
}