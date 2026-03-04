package llm

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/wwwzy/CloudAI/model"

	"github.com/cloudwego/eino-ext/components/model/ollama"
	"github.com/cloudwego/eino-ext/components/model/openai"

	eino_model "github.com/cloudwego/eino/components/model"
)

const (
	defaultLLMTimeout = 60 * time.Second
	serverOllama      = "ollama"
	serverOpenAI      = "openai"
	modelTypeLLM      = "llm"
)

// ---------------------------
// @brief 传入配置获取LLM模型
// ---------------------------
func GetLLMClient(ctx context.Context, cfg *model.Model) (eino_model.ToolCallingChatModel, error) {
	if cfg == nil {
		return nil, errors.New("input model configuration is nil")
	}
	if cfg.Type != modelTypeLLM {
		return nil, fmt.Errorf("model type is '%s', but expected '%s'", cfg.Type, modelTypeLLM)
	}

	switch strings.ToLower(cfg.Server) {
	case serverOllama:
		ollamaCfg := &ollama.ChatModelConfig{
			Model:   cfg.ModelName,
			BaseURL: cfg.BaseURL,
			Timeout: defaultLLMTimeout,
			Options: &ollama.Options{},
		}
		return ollama.NewChatModel(ctx, ollamaCfg)
	case serverOpenAI:
		openAICfg := &openai.ChatModelConfig{
			APIKey:    cfg.APIKey,
			Timeout:   defaultLLMTimeout,
			BaseURL:   cfg.BaseURL,
			Model:     cfg.ModelName,
			MaxTokens: &cfg.MaxTokens,
		}
		return openai.NewChatModel(ctx, openAICfg)

	default:
		return nil, fmt.Errorf("unsupported LLM server type: '%s'", cfg.Server)
	}
}
