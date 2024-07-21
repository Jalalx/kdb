package cmd

import (
	"fmt"
	"log"
	"math"
	"os"
	"time"

	"github.com/jalalx/kdb/repos"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

func listCommand(repo repos.EmbeddingRepo) *cobra.Command {
	var useUtcDate bool
	var limit int
	var command = &cobra.Command{
		Use:   "list",
		Short: "Lists the embedded entries.",
		Run: func(cmd *cobra.Command, args []string) {
			listHandler(useUtcDate, repo, limit)
		},
	}

	command.Flags().BoolVarP(&useUtcDate, "utc-date", "u", false, "If set, dates will be printed in the UTC format.")
	command.Flags().IntVarP(&limit, "limit", "n", 5, "Number of entries to be returned.")

	return command
}

func listHandler(
	useUtcDate bool,
	repo repos.EmbeddingRepo,
	limit int) error {

	items, err := repo.List(limit)
	if err != nil {
		log.Fatalln(err)
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Id", "Created At", "Content"})
	table.SetColMinWidth(2 /*Content Column Index*/, 50)
	for _, item := range items {
		dateStr := humanReadableTime(item.CreatedAt)
		if useUtcDate {
			dateStr = item.CreatedAt.String()
		}

		table.Append([]string{item.Id.String(), dateStr, item.Content})
	}
	table.Render()
	// No errors
	return nil
}

func humanReadableTime(t time.Time) string {
	now := time.Now()
	diff := now.Sub(t)

	if diff.Seconds() < 60 {
		return fmt.Sprintf("%d seconds ago", int(diff.Seconds()))
	}
	if diff.Minutes() < 60 {
		return fmt.Sprintf("%d minutes ago", int(diff.Minutes()))
	}
	if diff.Hours() < 24 {
		return fmt.Sprintf("%d hours ago", int(diff.Hours()))
	}
	if diff.Hours() < 48 {
		return "yesterday"
	}
	if diff.Hours() < 24*7 {
		return fmt.Sprintf("%d days ago", int(math.Floor(diff.Hours()/24)))
	}
	if diff.Hours() < 24*30 {
		return fmt.Sprintf("%d weeks ago", int(math.Floor(diff.Hours()/(24*7))))
	}
	if diff.Hours() < 24*365 {
		return fmt.Sprintf("%d months ago", int(math.Floor(diff.Hours()/(24*30))))
	}
	return fmt.Sprintf("%d years ago", int(math.Floor(diff.Hours()/(24*365))))
}
