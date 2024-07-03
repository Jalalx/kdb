package main

import (
	"database/sql"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	_ "github.com/marcboeker/go-duckdb"
	"github.com/spf13/cobra"
)

const (
	CONTEXT_LENGTH             = 1024
	EMBEDDING_MODEL_NAME       = "nomic-embed-text"
	EMBEDDING_MODEL_DIMENSIONS = 768
	DATABASE_VENDOR            = "duckdb"
	DATABASE_NAME              = "db/knowledgebase.ddb"
)

func readStdInput() string {
	stdInputBytes, err := io.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}

	return string(stdInputBytes)
}

func main() {
	var args InputArgs

	// Initialize the database
	db, err := Connect()
	if err != nil {
		panic(err)
	}

	defer db.Close()

	var rootCmd = &cobra.Command{
		Use:   "kdb",
		Short: "A knowledge database available as a command line tool",
		Run: func(cmd *cobra.Command, _ []string) {
			processInput(&args, db)
		},
	}

	rootCmd.Flags().BoolVarP(&args.Stdin, "stdin", "s", false, "Read from stdin")
	rootCmd.Flags().StringVarP(&args.Query, "query", "q", "", "Search query")
	rootCmd.Flags().StringVarP(&args.Embed, "embed", "e", "", "Text to embed")
	rootCmd.Flags().IntVarP(&args.List, "list", "l", 0, "List embedded texts")
	rootCmd.Flags().IntVarP(&args.Top, "top", "t", 0, "Sets the max number of results to return for query")
	rootCmd.Flags().BoolVarP(&args.Verbose, "verbose", "v", false, "Verbose output")

	// Execute the command
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func processInput(args *InputArgs, db *sql.DB) {

	embed := strings.TrimSpace(args.Embed)
	if args.Stdin && embed != "" {
		fmt.Println("Cannot use --stdin and --embed at the same time.")
		os.Exit(1)
	}

	query := strings.TrimSpace(args.Query)
	if args.Top > 0 && query == "" {
		fmt.Println("Cannot use --top without --query.")
		os.Exit(1)
	}

	if query != "" && len(query) > CONTEXT_LENGTH {
		fmt.Printf("Query must be less than %d characters.\n", CONTEXT_LENGTH)
		os.Exit(1)
	}

	if query != "" && args.Top == 0 {
		args.Top = 3
	}

	if embed != "" && len(embed) > CONTEXT_LENGTH {
		fmt.Printf("Embedding text must be less than %d characters.\n", CONTEXT_LENGTH)
		os.Exit(1)
	}

	if args.Top > 0 && args.List > 0 {
		fmt.Println("Cannot use --top and --list at the same time.")
		os.Exit(1)
	}

	if args.Stdin {
		stdin := readStdInput()
		eb, err := Embedd(strings.TrimSpace(stdin), EMBEDDING_MODEL_NAME, db)
		if err != nil {
			log.Fatalln(err)
		}
		InsertEmbedding(stdin, eb, db)
	} else if embed != "" {
		eb, err := Embedd(embed, EMBEDDING_MODEL_NAME, db)
		if err != nil {
			log.Fatalln(err)
		}
		InsertEmbedding(embed, eb, db)
	} else if query != "" {
		Query(query, args.Top, db)
	} else if args.List > 0 {
		ListEmbeddings(args.List, db)
	} else {
		fmt.Println("No action specified.")
	}
}
