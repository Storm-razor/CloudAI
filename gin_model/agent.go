package ginmodel

import (
	"github.com/cloudwego/eino/schema"
	"github.com/wwwzy/CloudAI/model"
)

// CreateAgentRequest 创建Agent请求
// 用于创建新的智能体配置
type CreateAgentRequest struct {
	Name        string `json:"name" binding:"required"` // 智能体名称
	Description string `json:"description"`             // 智能体描述
}

// UpdateAgentRequest 更新Agent请求
// 用于全量更新智能体的各项配置
type UpdateAgentRequest struct {
	ID          string                `json:"id" binding:"required"` // 智能体唯一ID
	Name        string                `json:"name"`                  // 新名称
	Description string                `json:"description"`           // 新描述
	LLMConfig   model.LLMConfig       `json:"llm_config"`            // 大模型配置(温度/模型名等)
	MCP         model.MCPConfig       `json:"mcp"`                   // MCP服务配置
	Tools       model.ToolsConfig     `json:"tools"`                 // 关联工具配置
	Prompt      string                `json:"prompt"`                // 系统提示词(System Prompt)
	Knowledge   model.KnowledgeConfig `json:"knowledge"`             // 关联知识库配置
}

// PageAgentRequest 分页查询Agent请求
type PageAgentRequest struct {
	Page int `form:"page,default=1"` // 页码(默认1)
	Size int `form:"size,default=10"`// 每页数量(默认10)
}

// UserMessage 用户消息结构体
// 用于执行Agent时传递上下文
type UserMessage struct {
	Query   string            `json:"query" binding:"required"` // 当前用户问题
	History []*schema.Message `json:"history"`                  // 历史对话记录(用于多轮对话)
}

// ExecuteAgentRequest 执行Agent请求
// 直接调试/运行Agent，不经过会话系统
type ExecuteAgentRequest struct {
	ID      string      `json:"id" binding:"required"`       // 请求ID
	AgentID string      `json:"agent_id" binding:"required"` // 目标Agent ID
	Message UserMessage `json:"message" binding:"required"`  // 输入消息
}
