package controller

import (
	"errors"
	"io"
	"log"
	"time"

	"github.com/gin-contrib/sse"
	"github.com/gin-gonic/gin"
	ginmodel "github.com/wwwzy/CloudAI/gin_model"
	"github.com/wwwzy/CloudAI/internal/service"
	"github.com/wwwzy/CloudAI/internal/utils"
	"github.com/wwwzy/CloudAI/pkgs/errcode"
	"github.com/wwwzy/CloudAI/pkgs/response"
)

type ConversationController struct {
	svc service.ConversationService
}

func NewConversationController(svc service.ConversationService) *ConversationController {
	return &ConversationController{svc: svc}
}

// CreateConversation 创建新会话
func (c *ConversationController) CreateConversation(ctx *gin.Context) {
	userID, err := utils.GetUserIDFromContext(ctx)
	if err != nil {
		response.UnauthorizedError(ctx, errcode.UnauthorizedError, "Failed to get user")
		return
	}

	var req ginmodel.CreateConvRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ParamError(ctx, errcode.ParamBindError, "Parameter error: "+err.Error())
		return
	}

	// 创建会话
	convID, err := c.svc.CreateConversation(ctx.Request.Context(), userID, req.AgentID)
	if err != nil {
		log.Printf("[Conversation Create] Error creating conversation: %v\n", err)
		response.InternalError(ctx, errcode.InternalServerError, "Failed to create conversation")
		return
	}

	response.SuccessWithMessage(ctx, "Conversation created successfully", gin.H{"conv_id": convID})
}

// ListConversations 获取用户所有会话
func (c *ConversationController) ListConversations(ctx *gin.Context) {
	userID, err := utils.GetUserIDFromContext(ctx)
	if err != nil {
		response.UnauthorizedError(ctx, errcode.UnauthorizedError, "Failed to get user")
		return
	}

	// 分页参数
	page := utils.StringToInt(ctx.DefaultQuery("page", "1"))
	size := utils.StringToInt(ctx.DefaultQuery("size", "10"))

	// 获取会话列表
	convs, count, err := c.svc.ListConversations(ctx.Request.Context(), userID, page, size)
	if err != nil {
		response.InternalError(ctx, errcode.InternalServerError, "Failed to list conversations: "+err.Error())
		return
	}

	// 返回分页数据
	response.PageSuccess(ctx, convs, count)
}

// ListAgentConversations 获取特定Agent的会话
func (c *ConversationController) ListAgentConversations(ctx *gin.Context) {
	userID, err := utils.GetUserIDFromContext(ctx)
	if err != nil {
		response.UnauthorizedError(ctx, errcode.UnauthorizedError, "Failed to get user")
		return
	}

	// 获取AgentID
	agentID := ctx.Query("agent_id")
	if agentID == "" {
		response.ParamError(ctx, errcode.ParamBindError, "Agent ID is required")
		return
	}

	// 分页参数
	page := utils.StringToInt(ctx.DefaultQuery("page", "1"))
	size := utils.StringToInt(ctx.DefaultQuery("size", "10"))

	// 获取会话列表
	convs, count, err := c.svc.ListAgentConversations(ctx.Request.Context(), userID, agentID, page, size)
	if err != nil {
		response.InternalError(ctx, errcode.InternalServerError, "Failed to list agent conversations: "+err.Error())
		return
	}

	// 返回分页数据
	response.PageSuccess(ctx, convs, count)
}

// GetConversationHistory 获取会话历史消息
func (c *ConversationController) GetConversationHistory(ctx *gin.Context) {
	// 获取会话ID
	convID := ctx.Query("conv_id")
	if convID == "" {
		response.ParamError(ctx, errcode.ParamBindError, "Conversation ID is required")
		return
	}

	// 限制参数
	limit := utils.StringToInt(ctx.DefaultQuery("limit", "50"))

	// 获取历史消息
	msgs, err := c.svc.GetConversationHistory(ctx.Request.Context(), convID, limit)
	if err != nil {
		response.InternalError(ctx, errcode.InternalServerError, "Failed to get conversation history: "+err.Error())
		return
	}

	// 返回历史消息
	response.SuccessWithMessage(ctx, "Conversation history retrieved successfully", gin.H{"messages": msgs})
}

// DeleteConversation 删除会话
func (c *ConversationController) DeleteConversation(ctx *gin.Context) {
	// 获取会话ID
	convID := ctx.Query("conv_id")
	if convID == "" {
		response.ParamError(ctx, errcode.ParamBindError, "Conversation ID is required")
		return
	}

	// 删除会话
	err := c.svc.DeleteConversation(ctx.Request.Context(), convID)
	if err != nil {
		log.Printf("[Conversation Delete] Error deleting conversation: %v\n", err)
		response.InternalError(ctx, errcode.InternalServerError, "Failed to delete conversation: "+err.Error())
		return
	}

	// 返回成功消息
	response.SuccessWithMessage(ctx, "Conversation deleted successfully", nil)
}

// DebugStreamAgent 调试模式，不保存历史
func (c *ConversationController) DebugStreamAgent(ctx *gin.Context) {
	userID, err := utils.GetUserIDFromContext(ctx)
	if err != nil {
		response.UnauthorizedError(ctx, errcode.UnauthorizedError, "Failed to get user")
		return
	}

	var req ginmodel.DebugRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ParamError(ctx, errcode.ParamBindError, "Parameter error: "+err.Error())
		return
	}

	// 调用debug模式流式处理
	sr, err := c.svc.DebugStreamAgent(ctx.Request.Context(), userID, req.AgentID, req.Message)
	if err != nil {
		response.InternalError(ctx, errcode.InternalServerError, "Agent execution failed: "+err.Error())
		return
	}

	// 设置SSE响应头
	ctx.Writer.Header().Set("Content-Type", "text/event-stream")
	ctx.Writer.Header().Set("Cache-Control", "no-cache")
	ctx.Writer.Header().Set("Connection", "keep-alive")
	ctx.Writer.Header().Set("Transfer-Encoding", "chunked")

	// 传输流
	sessionID := utils.GenerateUUID()
	done := make(chan struct{})
	defer func() {
		sr.Close()
		close(done)
		log.Printf("[Debug Stream] Finish Stream with ID: %s\n", sessionID)
	}()

	type streamResult struct {
		content string
		err     error
	}

	recvCh := make(chan streamResult, 1)
	go func() {
		for {
			msg, err := sr.Recv()
			result := streamResult{err: err}
			if err == nil {
				if msg == nil {
					result.err = errors.New("nil message received")
				} else {
					result.content = msg.Content
				}
			}
			select {
			case <-ctx.Request.Context().Done():
				close(recvCh)
				return
			case recvCh <- result:
				if result.err != nil {
					close(recvCh)
					return
				}
			}
		}
	}()

	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()

	// 流式响应
	ctx.Stream(func(w io.Writer) bool {
		select {
		case <-ctx.Request.Context().Done():
			log.Printf("[Debug Stream] Context done for session ID: %s\n", sessionID)
			return false
		case <-done:
			return false
		case result, ok := <-recvCh:
			if !ok {
				return false
			}
			if errors.Is(result.err, io.EOF) {
				log.Printf("[Debug Stream] EOF received for session ID: %s\n", sessionID)
				return false
			}
			if result.err != nil {
				log.Printf("[Debug Stream] Error receiving message: %v\n", result.err)
				return false
			}

			if err := sse.Encode(w, sse.Event{
				Data: []byte(result.content),
			}); err != nil {
				log.Printf("[Debug Stream] Error sending message: %v\n", err)
				return false
			}

			ctx.Writer.Flush()
			return true
		case <-ticker.C:
			if _, err := w.Write([]byte(": ping\n\n")); err != nil {
				log.Printf("[Debug Stream] Error sending heartbeat: %v\n", err)
				return false
			}
			ctx.Writer.Flush()
			return true
		}
	})
}

// StreamConversation 会话模式，保存历史
func (c *ConversationController) StreamConversation(ctx *gin.Context) {
	userID, err := utils.GetUserIDFromContext(ctx)
	if err != nil {
		response.UnauthorizedError(ctx, errcode.UnauthorizedError, "Failed to get user")
		return
	}

	var req ginmodel.ConvRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ParamError(ctx, errcode.ParamBindError, "Parameter error: "+err.Error())
		return
	}

	if req.ConvID == "" {
		convID, err := c.svc.CreateConversation(ctx.Request.Context(), userID, req.AgentID)
		if err != nil {
			log.Printf("[Conversation Create] Error creating conversation: %v\n", err)
			response.InternalError(ctx, errcode.InternalServerError, "Failed to create conversation")
			return
		}
		req.ConvID = convID
	}

	// 调用会话模式流式处理
	sr, err := c.svc.StreamAgentWithConversation(ctx.Request.Context(), userID, req.AgentID, req.ConvID, req.Message)
	if err != nil {
		log.Printf("[Conversation Stream] Error running agent: %v\n", err)
		response.InternalError(ctx, errcode.InternalServerError, "Agent execution failed")
		return
	}

	// 设置SSE响应头
	ctx.Writer.Header().Set("Content-Type", "text/event-stream")
	ctx.Writer.Header().Set("Cache-Control", "no-cache")
	ctx.Writer.Header().Set("Connection", "keep-alive")
	ctx.Writer.Header().Set("Transfer-Encoding", "chunked")

	// 传输流
	done := make(chan struct{})
	defer func() {
		sr.Close()
		close(done)
		log.Printf("[Conversation Stream] Finish Stream with ConvID: %s\n", req.ConvID)
	}()

	type streamResult struct {
		content string
		err     error
	}

	recvCh := make(chan streamResult, 1)
	go func() {
		for {
			msg, err := sr.Recv()
			result := streamResult{err: err}
			if err == nil {
				if msg == nil {
					result.err = errors.New("nil message received")
				} else {
					result.content = msg.Content
				}
			}
			select {
			case <-ctx.Request.Context().Done():
				close(recvCh)
				return
			case recvCh <- result:
				if result.err != nil {
					close(recvCh)
					return
				}
			}
		}
	}()

	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()

	// 流式响应
	ctx.Stream(func(w io.Writer) bool {
		select {
		case <-ctx.Request.Context().Done():
			log.Printf("[Conversation Stream] Context done for ConvID: %s\n", req.ConvID)
			return false
		case <-done:
			return false
		case result, ok := <-recvCh:
			if !ok {
				return false
			}
			if errors.Is(result.err, io.EOF) {
				log.Printf("[Conversation Stream] EOF received for ConvID: %s\n", req.ConvID)
				return false
			}
			if result.err != nil {
				log.Printf("[Conversation Stream] Error receiving message: %v\n", result.err)
				return false
			}
			if err := sse.Encode(w, sse.Event{
				Data: []byte(result.content),
			}); err != nil {
				log.Printf("[Conversation Stream] Error sending message: %v\n", err)
				return false
			}
			ctx.Writer.Flush()
			return true
		case <-ticker.C:
			if _, err := w.Write([]byte(": ping\n\n")); err != nil {
				log.Printf("[Conversation Stream] Error sending heartbeat: %v\n", err)
				return false
			}
			ctx.Writer.Flush()
			return true
		}
	})
}
