package main

import (
	"context"

	_ "github.com/marcboeker/go-duckdb"
	"github.com/spf13/cobra"

	"github.com/jalalx/kdb/cmd"
	"github.com/jalalx/kdb/llms"
	"github.com/jalalx/kdb/repos"
)

const (
	CONTEXT_LENGTH             = 4000
	EMBEDDING_MODEL_NAME       = "nomic-embed-text"
	EMBEDDING_MODEL_DIMENSIONS = 768
)

func main() {
	MakeKdbDirIfNeeded()

	// Initialize the repo
	repo, err := repos.NewRepository()
	if err != nil {
		panic(err)
	}

	defer repo.Close()

	llmProvider, err := llms.NewLlmProvider()
	if err != nil {
		panic(err)
	}

	err = repo.Connect()
	if err != nil {
		panic(err)
	}

	err = repo.Init(EMBEDDING_MODEL_DIMENSIONS)
	if err != nil {
		panic(err)
	}

	rootCmd := cmd.NewCLI(repo, llmProvider, EMBEDDING_MODEL_NAME)
	cobra.CheckErr(rootCmd.ExecuteContext(context.Background()))
}
