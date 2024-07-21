package cmd

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/jalalx/kdb/llms"
	"github.com/jalalx/kdb/repos"
	"github.com/spf13/cobra"
)

func embedCommand(
	repo repos.EmbeddingRepo,
	llmProvider llms.LlmProvider) *cobra.Command {

	var content string
	var silent bool
	var fromStdInput bool

	var command = &cobra.Command{
		Use:   "embed [text]",
		Short: "Stores the given text in the vector database.",
		Run: func(cmd *cobra.Command, args []string) {
			if fromStdInput {
				if len(args) != 0 {
					cmd.Println("When using --stdin, no text must be provided.")
					return
				}

				content = readStdInput()
			} else {
				if len(args) < 1 {
					cmd.Println("No text is provided.")
					return
				}

				if len(args) > 1 {
					cmd.Println("Too many parameters are provided.")
					return
				}

				content = args[0]
			}

			embedHandler(repo, llmProvider, content, silent)
		},
		Aliases: []string{"save"},
	}
	command.Flags().BoolVar(&silent, "silent", false, "If set, no response will be printed unless there is an error.")
	command.Flags().BoolVar(&fromStdInput, "stdin", false, "If set, Content will be read from the standard input.")

	return command
}

func embedHandler(
	repo repos.EmbeddingRepo,
	llmProvider llms.LlmProvider,
	content string,
	silent bool) error {

	vector, err := llmProvider.GetEmbedding(content)
	if err != nil {
		return err
	}

	result, err := repo.Insert(content, vector)
	if err != nil {
		log.Fatalln(err)
	}

	if !silent {
		fmt.Printf("%s\n", result.Id)
	}

	return nil
}

func readStdInput() string {
	stdInputBytes, err := io.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}

	return string(stdInputBytes)
}
