package cmd

import (
	"errors"
	"fmt"
	"log"
	"os"
	"text/tabwriter"

	"github.com/jalalx/kdb/repos"
	"github.com/spf13/cobra"
)

func listCommand(repo repos.EmbeddingRepo) *cobra.Command {
	var command = &cobra.Command{
		Use:   "list",
		Short: "Lists the embedded entries.",
		Run: func(cmd *cobra.Command, args []string) {
			listHandler(cmd, args, repo)
		},
	}
	command.Flags().Int("limit", 5, "Number of entries to be returned.")
	command.Flags().String("separator", "\n", "Separator to be used between entries.")

	return command
}

func listHandler(cmd *cobra.Command, _ []string, repo repos.EmbeddingRepo) error {
	limit, errLimit := cmd.Flags().GetInt("limit")
	if errLimit != nil {
		return errors.New("failure in getting the limit")
	}

	separator, _ := cmd.Flags().GetString("separator")

	items, err := repo.List(limit)
	if err != nil {
		log.Fatalln(err)
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "Id\tCreated At\tContent\t")
	fmt.Fprintln(w, "--\t----------\t-------\t")
	for _, item := range items {
		row := fmt.Sprintf("%s\t%s\t%s%s", item.Id, item.CreatedAt, item.Content, separator)
		fmt.Fprint(w, row)
	}

	w.Flush()

	// No errors
	return nil
}
