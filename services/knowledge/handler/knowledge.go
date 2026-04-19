package handler

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/chunfengshili/sa/services/knowledge/embedding"
	"github.com/chunfengshili/sa/services/knowledge/qdrant"
	"github.com/chunfengshili/sa/pkg/response"
)

// KnowledgeHandler 知识库处理器
type KnowledgeHandler struct {
	embeddingSvc *embedding.Service
	qdrant      *qdrant.Client
	collection   string
}

// NewKnowledgeHandler 创建知识库处理器
func NewKnowledgeHandler(emb *embedding.Service, qd *qdrant.Client, collection string) *KnowledgeHandler {
	return &KnowledgeHandler{
		embeddingSvc: emb,
		qdrant:      qd,
		collection:   collection,
	}
}

// AddDocumentRequest 添加文档请求
type AddDocumentRequest struct {
	Content string                 `json:"content" binding:"required"`
	Meta    map[string]interface{} `json:"meta"`
}

// AddDocument 添加文档
func (h *KnowledgeHandler) AddDocument(c *gin.Context) {
	var req AddDocumentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "无效的请求参数")
		return
	}

	// 生成 embedding
	vector, err := h.embeddingSvc.EmbeddingSingle(c.Request.Context(), req.Content)
	if err != nil {
		response.InternalError(c, "生成向量失败: "+err.Error())
		return
	}

	// 生成 ID
	id := uuid.New().String()

	// 存储到 Qdrant
	point := qdrant.Point{
		ID:      id,
		Vector:  vector,
		Payload: map[string]interface{}{"content": req.Content, "meta": req.Meta},
	}
	err = h.qdrant.UpsertVectors(c.Request.Context(), h.collection, []qdrant.Point{point})
	if err != nil {
		response.InternalError(c, "存储向量失败: "+err.Error())
		return
	}

	response.Success(c, gin.H{
		"id":      id,
		"message": "文档添加成功",
	})
}

// SearchRequest 搜索请求
type SearchRequest struct {
	Query string `json:"query" binding:"required"`
	Limit int    `json:"limit"`
}

// SearchResult 搜索结果
type SearchResult struct {
	ID      string                 `json:"id"`
	Score   float32                `json:"score"`
	Content string                 `json:"content"`
	Meta    map[string]interface{} `json:"meta"`
}

// Search 搜索相似文档
func (h *KnowledgeHandler) Search(c *gin.Context) {
	var req SearchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "无效的请求参数")
		return
	}

	if req.Limit == 0 {
		req.Limit = 10
	}

	// 生成查询向量
	vector, err := h.embeddingSvc.EmbeddingSingle(c.Request.Context(), req.Query)
	if err != nil {
		response.InternalError(c, "生成向量失败: "+err.Error())
		return
	}

	// 搜索
	results, err := h.qdrant.SearchVectors(c.Request.Context(), h.collection, vector, req.Limit)
	if err != nil {
		response.InternalError(c, "搜索失败: "+err.Error())
		return
	}

	// 转换结果
	searchResults := make([]SearchResult, len(results))
	for i, r := range results {
		content := ""
		meta := make(map[string]interface{})
		if r.Payload != nil {
			if c, ok := r.Payload["content"].(string); ok {
				content = c
			}
			if m, ok := r.Payload["meta"].(map[string]interface{}); ok {
				meta = m
			}
		}
		searchResults[i] = SearchResult{
			ID:      r.ID,
			Score:   r.Score,
			Content: content,
			Meta:    meta,
		}
	}

	response.Success(c, searchResults)
}

// CreateCollectionRequest 创建集合请求
type CreateCollectionRequest struct {
	Name       string `json:"name" binding:"required"`
	VectorSize int    `json:"vector_size"`
}

// CreateCollection 创建集合
func (h *KnowledgeHandler) CreateCollection(c *gin.Context) {
	var req CreateCollectionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "无效的请求参数")
		return
	}

	if req.VectorSize == 0 {
		req.VectorSize = 1024 // text-embedding-v3 默认维度
	}

	err := h.qdrant.CreateCollection(context.Background(), req.Name, req.VectorSize)
	if err != nil {
		response.InternalError(c, "创建集合失败: "+err.Error())
		return
	}

	response.Success(c, gin.H{
		"name":       req.Name,
		"vectorSize": req.VectorSize,
	})
}