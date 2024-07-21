package cmd

import (
	"fmt"

	"github.com/jalalx/kdb/llms"
	"github.com/jalalx/kdb/repos"
	"github.com/spf13/cobra"
)

var (
	Version string
	GitHash string
)

func NewCLI(repo repos.EmbeddingRepo, llmProvider llms.LlmProvider, embeddingModelName string) *cobra.Command {
	cobra.EnableCommandSorting = false

	var rootCmd = &cobra.Command{
		Use:           "kdb",
		Short:         "Knowledgebase in your command line",
		SilenceUsage:  true,
		SilenceErrors: true,
		CompletionOptions: cobra.CompletionOptions{
			DisableDefaultCmd: true,
		},
		Run: func(cmd *cobra.Command, args []string) {
			if version, _ := cmd.Flags().GetBool("version"); version {
				versionHandler(cmd, args)
				return
			}

			cmd.Print(cmd.UsageString())
		},
	}

	rootCmd.Flags().BoolP("version", "v", false, "Show version information")

	rootCmd.AddCommand(listCommand(repo))
	rootCmd.AddCommand(queryCommand(repo, llmProvider, embeddingModelName))
	rootCmd.AddCommand(embedCommand(repo, llmProvider, embeddingModelName))
	rootCmd.AddCommand(deleteCommand(repo))
	return rootCmd
}

func versionHandler(_ *cobra.Command, _ []string) {
	fmt.Printf("Version: %s-%s\n", Version, GitHash)
}
