package embedding

import (
	"context"
	"fmt"
	"time"

	einoEmbedding "github.com/cloudwego/eino/components/embedding"
	"github.com/wwwzy/CloudAI/model"
)

const (
	ProviderOpenAI = "openai" //openAI模型
	ProviderOllama = "ollama" //本地模型
)

type EmbeddingOption func(*EmbeddingOptions)

type EmbeddingOptions struct {
	Timeout *time.Duration
}

func WithTimeout(timeout time.Duration) EmbeddingOption {
	return func(o *EmbeddingOptions) {
		o.Timeout = &timeout
	}
}

// embedding服务通用接口
type EmbeddingService interface {
	New(ctx context.Context, cfg *model.Model, opts ...EmbeddingOption) (EmbeddingService, error)
	// EmbedStrings 将文本转换为向量表示
	EmbedStrings(ctx context.Context, texts []string, opts ...einoEmbedding.Option) ([][]float64, error)
	// GetDimension 返回嵌入向量的维度
	GetDimension() int
}

var embeddingMap = make(map[string]EmbeddingService)

// ---------------------------
// @brief 注册embedding服务
// ---------------------------
func register(name string, embeddingService EmbeddingService) {
	embeddingMap[name] = embeddingService
}

// ---------------------------
// @brief 外界调用embedding服务
// ---------------------------
func NewEmbeddingService(ctx context.Context, cfg *model.Model, opts ...EmbeddingOption) (EmbeddingService, error) {
	if cfg == nil {
		return nil, fmt.Errorf("embedding config is nil")
	}

	if cfg.Server == "" {
		return nil, fmt.Errorf("embedding config server is empty")
	}

	// 获取实例
	if embedding, ok := embeddingMap[cfg.Server]; ok {
		return embedding.New(ctx, cfg, opts...)
	}
	return nil, fmt.Errorf("不支持的嵌入服务提供者: %s", cfg.Server)
}
