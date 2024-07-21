package llms

import "github.com/jalalx/kdb/config"

type LlmProvider interface {
	GetEmbedding(prompt string) ([]float64, error)
}

func NewLlmProvider(cfg *config.KdbEmbeddingConfig) (LlmProvider, error) {
	return NewOllamaLlmProvider(cfg)
}
