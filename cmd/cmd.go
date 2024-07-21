package cmd

import (
	"fmt"

	"github.com/jalalx/kdb/llms"
	"github.com/jalalx/kdb/repos"
	"github.com/spf13/cobra"
)

func NewCLI(
	repo repos.EmbeddingRepo,
	llmProvider llms.LlmProvider,
	version string,
	githash string) *cobra.Command {
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
			if v, _ := cmd.Flags().GetBool("version"); v {
				versionHandler(version, githash)
				return
			}

			cmd.Print(cmd.UsageString())
		},
	}

	rootCmd.Flags().BoolP("version", "v", false, "Show version information")

	rootCmd.AddCommand(listCommand(repo))
	rootCmd.AddCommand(queryCommand(repo, llmProvider))
	rootCmd.AddCommand(embedCommand(repo, llmProvider))
	rootCmd.AddCommand(deleteCommand(repo))
	return rootCmd
}

func versionHandler(version string, githash string) {
	fmt.Printf("Version: %s-%s\n", version, githash)
}
