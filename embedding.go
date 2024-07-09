package main

import (
	ollamaclient "github.com/xyproto/ollamaclient"
)

func getClientConfig() (*ollamaclient.Config, error) {
	oc := ollamaclient.NewWithModel(EMBEDDING_MODEL_NAME)
	oc.Verbose = false

	if err := oc.PullIfNeeded(); err != nil {
		return nil, err
	}

	return oc, nil
}

func Embedd(prompt string, model string) ([]float64, error) {
	oc, err := getClientConfig()
	if err != nil {
		return nil, err
	}

	embeddings, err := oc.Embeddings(prompt)
	if err != nil {
		return nil, err
	}

	return embeddings, nil
}
