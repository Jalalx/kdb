package cmd

import (
	"fmt"

	"github.com/jalalx/kdb/llms"
	"github.com/jalalx/kdb/repos"
	"github.com/spf13/cobra"
)

func queryCommand(
	repo repos.EmbeddingRepo,
	llmProvider llms.LlmProvider,
	modelName string) *cobra.Command {

	var limit int
	var showIdentifiers bool

	var command = &cobra.Command{
		Use:   "query [text]",
		Short: "Search in the embedded entries.",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			queryText := args[0]
			queryHandler(repo, llmProvider, queryText, limit, showIdentifiers, modelName)
		},
	}
	command.Flags().IntVar(&limit, "limit", 5, "Number of entries to be returned.")
	command.Flags().BoolVar(&showIdentifiers, "id", false, "Number of entries to be returned.")

	return command
}

func queryHandler(
	repo repos.EmbeddingRepo,
	llmProvider llms.LlmProvider,
	query string,
	limit int,
	showIdentifiers bool,
	modelName string) error {

	vector, err := llmProvider.GetEmbedding(query, modelName)
	if err != nil {
		return err
	}

	items, err := repo.Query(vector, limit)
	if err != nil {
		return err
	}

	for _, item := range items {
		if showIdentifiers {
			fmt.Printf("%s\t%s\n", item.Id, item.Content)
		} else {
			fmt.Println(item.Content)
		}
	}

	return nil
}
