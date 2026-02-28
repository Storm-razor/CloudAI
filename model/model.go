package model

import "time"

type Model struct {
	// 基础信息
	ID        string `gorm:"primaryKey;type:char(36)"`
	UserID    uint   `gorm:"index"`    // 用户ID
	Type      string `gorm:"not null"` // 模型的类型：embedding/llm
	ShowName  string `gorm:"not null"` // 显示名称
	Server    string `gorm:"not null"` // 模型的供应商：openai/ollama
	BaseURL   string `gorm:"not null"` // API基础地址
	ModelName string `gorm:"not null"` // 模型标识符，例如 deepseek-chat，text-embedding-v3
	APIKey    string // 访问密钥，ollama一般不需要

	// Embedding模型字段
	Dimension int // 向量维度(embedding必填)

	// LLM模型字段
	MaxOutputLength int  `gorm:"default:4096"`
	Function        bool `gorm:"default:false"`

	// 通用字段
	MaxTokens int       `gorm:"default:1024"` // 限制最大的输入长度
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}
