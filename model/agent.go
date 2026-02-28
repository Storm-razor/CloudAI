package model

import (
	"time"
)

// Agent
type Agent struct {
	ID          string    `gorm:"primaryKey;type:char(36)"`
	UserID      uint      `gorm:"index"`
	Name        string    `gorm:"not null"`
	Description string    `gorm:"type:text"`
	AgentSchema string    `gorm:"type:json"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
}

// AgentSchema 配置Agent
type AgentSchema struct {
	LLMConfig LLMConfig       `json:"llm_config"`
	MCP       MCPConfig       `json:"mcp"`
	Tools     ToolsConfig     `json:"tools"`
	Prompt    string          `json:"prompt"`
	Knowledge KnowledgeConfig `json:"knowledge"`
}

// LLMConfig 配置Agent关联的LLM模型
type LLMConfig struct {
	ModelID         string  `json:"model_id"`
	Temperature     float64 `json:"temperature"`
	TopP            float64 `json:"top_p"`
	MaxOutputLength int     `json:"max_output_length"`
	Thinking        bool    `json:"thinking"`
}

// MCPConfig 配置MCP SSE服务器
type MCPConfig struct {
	Servers []string `json:"servers"`
}

// ToolsConfig 配置Agent关联的工具IDs（考虑到MCP到存在，目前没有实现Tools模块）
type ToolsConfig struct {
	ToolIDs []string `json:"tool_ids"`
}

// KnowledgeConfig Agent关联的知识库IDs
type KnowledgeConfig struct {
	KnowledgeIDs []string `json:"knowledge_ids"`
	TopK         int      `json:"top_k"`
}
