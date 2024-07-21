package llms

import (
	"fmt"
	"strings"

	"github.com/jalalx/kdb/config"
	ollamaclient "github.com/xyproto/ollamaclient"
)

type OllamaLlmProvider struct {
	oc *ollamaclient.Config
}

func NewOllamaLlmProvider(cfg *config.KdbEmbeddingConfig) (*OllamaLlmProvider, error) {

	if strings.ToLower(cfg.Provider) == "ollama" {
		oc := ollamaclient.NewWithModelAndAddr(cfg.Model.Name, cfg.URL)
		oc.Verbose = false
		return &OllamaLlmProvider{
			oc: oc,
		}, nil
	}

	return nil, fmt.Errorf("llm provider is not supported: %s", cfg.Provider)
}

func (p *OllamaLlmProvider) GetEmbedding(prompt string) ([]float64, error) {
	if err := p.oc.PullIfNeeded(); err != nil {
		return nil, err
	}

	embeddings, err := p.oc.Embeddings(prompt)
	if err != nil {
		return nil, err
	}

	return embeddings, nil
}
