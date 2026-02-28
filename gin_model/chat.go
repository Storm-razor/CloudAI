package ginmodel

import "github.com/cloudwego/eino/schema"

// ChatResponse 普通对话响应
// 包含AI的回答以及引用的知识库文档
type ChatResponse struct {
	Response   string             `json:"response"`   // AI回复的完整文本
	References []*schema.Document `json:"references"` // 引用来源(RAG检索到的文档片段)
}

// ChatRequest 对话请求
// 最基础的单轮对话请求
type ChatRequest struct {
	Query string   `json:"query"` // 用户问题
	KBs   []string `json:"kbs"`   // 指定搜索的知识库ID列表
}

// ChatStreamResponse 流式对话响应(SSE)
// 兼容 OpenAI API 格式，用于前端打字机效果展示
type ChatStreamResponse struct {
	ID      string             `json:"id"`      // 消息ID
	Object  string             `json:"object"`  // 对象类型(通常为 "chat.completion.chunk")
	Created int64              `json:"created"` // 创建时间戳
	Model   string             `json:"model"`   // 使用的模型名称
	Choices []ChatStreamChoice `json:"choices"` // 响应选项列表
}

// ChatStreamChoice 流式响应选项
type ChatStreamChoice struct {
	Delta        ChatStreamDelta `json:"delta"`         // 增量内容(这一帧生成的字)
	Index        int             `json:"index"`         // 索引
	FinishReason *string         `json:"finish_reason"` // 结束原因(null表示生成中, "stop"表示生成结束)
}

// ChatStreamDelta 流式响应增量内容
type ChatStreamDelta struct {
	Role    string `json:"role,omitempty"`    // 角色(仅第一帧返回)
	Content string `json:"content,omitempty"` // 文本内容(后续帧返回)
}
