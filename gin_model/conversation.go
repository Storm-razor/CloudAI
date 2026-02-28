package ginmodel

// DebugRequest Agent 调试请求
// 用于快速验证 Agent 逻辑，不创建持久化会话
type DebugRequest struct {
	AgentID string `json:"agent_id" binding:"required"` // 目标Agent ID
	Message string `json:"message" binding:"required"`  // 调试输入内容
}

// CreateConvRequest 创建会话请求
// 初始化一个新的对话窗口
type CreateConvRequest struct {
	AgentID string `json:"agent_id" binding:"required"` // 关联的Agent ID
}

// ConvRequest 发送消息请求
// 在已有会话中发送新消息
type ConvRequest struct {
	AgentID string `json:"agent_id" binding:"required"` // Agent ID(通常用于校验)
	Message string `json:"message" binding:"required"`  // 用户发送的消息内容
	ConvID  string `json:"conv_id"`                     // 会话ID(如果是首条消息可能为空，但推荐先创建会话)
}
