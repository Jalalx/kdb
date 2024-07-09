package llms

type LlmProvider interface {
	GetEmbedding(prompt string, model string) ([]float64, error)
}

func NewLlmProvider() (LlmProvider, error) {
	return NewOllamaLlmProvider(), nil
}
