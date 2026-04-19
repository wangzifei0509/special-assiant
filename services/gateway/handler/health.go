package handler

import (
	"github.com/gin-gonic/gin"

	"github.com/chunfengshili/sa/pkg/response"
)

// HealthHandler 健康检查处理器
type HealthHandler struct{}

// NewHealthHandler 创建健康检查处理器
func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

// Check 健康检查
func (h *HealthHandler) Check(c *gin.Context) {
	response.Success(c, gin.H{
		"status": "ok",
	})
}