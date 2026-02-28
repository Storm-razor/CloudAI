package ginmodel

// CreateModelRequest 注册模型请求
// 支持添加 LLM (大语言模型) 和 Embedding (向量模型)
type CreateModelRequest struct {
	// 基础信息
	Type      string `json:"type" binding:"required,oneof=embedding llm"` // 模型类型: embedding 或 llm
	ShowName  string `json:"name" binding:"required"`                     // 显示名称(用户自定义)
	Server    string `json:"server" binding:"required"`                   // 供应商(如 openai, zhipu, ollama)
	BaseURL   string `json:"base_url" binding:"required,url"`             // API地址
	ModelName string `json:"model" binding:"required"`                    // 实际模型标识(如 gpt-4, llama3)
	APIKey    string `json:"api_key"`                                     // API密钥

	// Embedding 特有参数
	Dimension int `json:"dimension"` // 向量维度(如 1536, 1024)

	// LLM 特有参数
	MaxOutputLength int  `json:"max_output_length"` // 最大输出长度
	Function        bool `json:"function"`          // 是否支持函数调用

	// 通用字段
	MaxTokens int `json:"max_tokens"` // 上下文最大Token数
}

// PageModelRequest 分页查询模型请求
type PageModelRequest struct {
	Type string `form:"type"`            // 按类型筛选(embedding/llm)
	Page int    `form:"page,default=1"`  // 页码
	Size int    `form:"size,default=10"` // 每页数量
}

// UpdateModelRequest 更新模型配置请求
type UpdateModelRequest struct {
	ID        string `json:"id" binding:"required"`       // 模型ID
	// 基础信息
	ShowName  string `json:"name"`                        // 新显示名称
	Server    string `json:"server"`                      // 新供应商
	BaseURL   string `json:"base_url"`                    // 新API地址
	ModelName string `json:"model"`                       // 新模型标识
	APIKey    string `json:"api_key"`                     // 新API密钥

	// Embedding
	Dimension int `json:"dimension"`                      // 新向量维度

	// LLM
	MaxOutputLength int  `json:"max_output_length"`       // 新最大输出长度
	Function        bool `json:"function"`                // 是否支持函数调用

	// 通用字段
	MaxTokens int `json:"max_tokens"`                     // 新上下文最大Token数
}
