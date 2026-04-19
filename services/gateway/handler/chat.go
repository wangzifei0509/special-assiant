package handler

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/sashabaranov/go-openai"

	"github.com/chunfengshili/sa/pkg/llm"
	"github.com/chunfengshili/sa/pkg/response"
)

// ChatHandler 聊天处理器
type ChatHandler struct {
	client *llm.Client
}

// NewChatHandler 创建聊天处理器
func NewChatHandler(client *llm.Client) *ChatHandler {
	return &ChatHandler{client: client}
}

// ChatRequest 聊天请求
type ChatRequest struct {
	Message string `json:"message" binding:"required"`
}

// ChatResponse 聊天响应
type ChatResponse struct {
	Reply string `json:"reply"`
}

// Chat 单次聊天
func (h *ChatHandler) Chat(c *gin.Context) {
	var req ChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "无效的请求参数")
		return
	}

	resp, err := h.client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: h.client.Model(),
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: req.Message,
				},
			},
		},
	)
	if err != nil {
		response.InternalError(c, "AI 服务调用失败: "+err.Error())
		return
	}

	if len(resp.Choices) == 0 {
		response.InternalError(c, "AI 服务返回空响应")
		return
	}

	response.Success(c, ChatResponse{
		Reply: resp.Choices[0].Message.Content,
	})
}

// ChatStream 流式聊天
func (h *ChatHandler) ChatStream(c *gin.Context) {
	var req ChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "无效的请求参数")
		return
	}

	// 设置 SSE 头
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")

	stream, err := h.client.CreateChatCompletionStream(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: h.client.Model(),
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: req.Message,
				},
			},
		},
	)
	if err != nil {
		response.InternalError(c, "AI 服务调用失败")
		return
	}
	defer stream.Close()

	for {
		resp, err := stream.Recv()
		if err != nil {
			break
		}
		if len(resp.Choices) > 0 {
			_, _ = c.Writer.WriteString("data: " + resp.Choices[0].Delta.Content + "\n\n")
			c.Writer.Flush()
		}
	}

	_, _ = c.Writer.WriteString("data: [DONE]\n\n")
}