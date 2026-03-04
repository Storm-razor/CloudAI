package embedding

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"

	einoEmbedding "github.com/cloudwego/eino/components/embedding"
	"github.com/ollama/ollama/api"
	"github.com/wwwzy/CloudAI/model"
)

func init() {
	register(ProviderOllama, &ollamaEmbedder{})
}

type OllamaEmbeddingConfig struct {
	BaseURL    string
	Model      string
	Dimension  *int
	Timeout    *time.Duration
	HTTPClient *http.Client
}

type ollamaEmbedder struct {
	cli  *api.Client
	conf *OllamaEmbeddingConfig
}

func (o *ollamaEmbedder) New(ctx context.Context, cfg *model.Model, opts ...EmbeddingOption) (EmbeddingService, error) {
	if err := checkCfg(cfg); err != nil {
		return nil, err
	}

	options := &EmbeddingOptions{}
	for _, opt := range opts {
		opt(options)
	}

	config := &OllamaEmbeddingConfig{
		BaseURL:   cfg.BaseURL,
		Model:     cfg.ModelName,
		Dimension: &cfg.Dimension,
		Timeout:   options.Timeout,
	}

	//构造client
	var httpClient *http.Client
	if config.HTTPClient != nil {
		httpClient = config.HTTPClient
	} else {
		httpClient = &http.Client{Timeout: *config.Timeout}
	}

	//构造url
	baseURL, err := url.Parse(config.BaseURL)
	if err != nil {
		return nil, fmt.Errorf("invalid base URL: %w", err)
	}

	cli := api.NewClient(baseURL, httpClient)

	return &ollamaEmbedder{
		cli:  cli,
		conf: config,
	}, nil

}

func (o *ollamaEmbedder) EmbedStrings(ctx context.Context, texts []string, opts ...einoEmbedding.Option) (
	embeddings [][]float64, err error) {
	req := &api.EmbedRequest{
		Model: o.conf.Model,
		Input: texts,
	}
	resp, err := o.cli.Embed(ctx, req)
	if err != nil {
		return nil, err
	}

	embeddings = make([][]float64, len(resp.Embeddings))
	for i, d := range resp.Embeddings {
		res := make([]float64, len(d))
		for j, emb := range d {
			res[j] = float64(emb)
		}
		embeddings[i] = res
	}

	return embeddings, nil
}

func (o *ollamaEmbedder) GetDimension() int {
	return *o.conf.Dimension
}

// ---------------------------
// @brief 检查model信息
// ---------------------------
func checkCfg(cfg *model.Model) error {
	if cfg.BaseURL == "" {
		return fmt.Errorf("ollama base URL cannot be empty")
	}

	if _, err := url.Parse(cfg.BaseURL); err != nil {
		return fmt.Errorf("invalid ollama base URL: %w", err)
	}

	if cfg.ModelName == "" {
		return fmt.Errorf("ollama model name cannot be empty")
	}

	if cfg.Dimension <= 0 {
		return fmt.Errorf("ollama embedding dimension must be positive")
	}

	return nil
}
