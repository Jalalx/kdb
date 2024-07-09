package llms

import (
	ollamaclient "github.com/xyproto/ollamaclient"
)

type OllamaLlmProvider struct {
	oc *ollamaclient.Config
}

func NewOllamaLlmProvider() *OllamaLlmProvider {
	return &OllamaLlmProvider{}
}

func (p *OllamaLlmProvider) GetEmbedding(prompt string, model string) ([]float64, error) {
	p.oc = ollamaclient.NewWithModel(model)
	p.oc.Verbose = false

	if err := p.oc.PullIfNeeded(); err != nil {
		return nil, err
	}

	embeddings, err := p.oc.Embeddings(prompt)
	if err != nil {
		return nil, err
	}

	return embeddings, nil
}
