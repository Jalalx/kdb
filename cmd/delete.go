package cmd

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/jalalx/kdb/repos"
	"github.com/spf13/cobra"
)

func deleteCommand(repo repos.EmbeddingRepo) *cobra.Command {

	var silent bool

	var command = &cobra.Command{
		Use:   "delete [id]",
		Short: "Deletes an entry by given id.",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			id, err := uuid.Parse(args[0])
			if err != nil {
				cmd.Printf("Given input is not a valid UUID: '%s'.", args[0])
				return
			}

			deleteHandler(repo, id, silent)
		},
		Aliases: []string{"remove"},
	}
	command.Flags().BoolVar(&silent, "silent", false, "If set, no response will be printed unless there is an error.")

	return command
}

func deleteHandler(
	repo repos.EmbeddingRepo,
	id uuid.UUID,
	silent bool) error {

	result, err := repo.Delete(id)
	if err != nil {
		return err
	}

	if !silent {
		if result {
			fmt.Printf("Entry deleted: %s\n", id)
		} else {
			fmt.Printf("Item not found for id '%s'\n", id)
		}
	}

	return nil
}
