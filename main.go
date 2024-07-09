package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	_ "github.com/marcboeker/go-duckdb"
	"github.com/spf13/cobra"

	"github.com/jalalx/kdb/database"
)

const (
	CONTEXT_LENGTH             = 4000
	EMBEDDING_MODEL_NAME       = "nomic-embed-text"
	EMBEDDING_MODEL_DIMENSIONS = 768
)

var (
	Version string
	GitHash string
)

func readStdInput() string {
	stdInputBytes, err := io.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}

	return string(stdInputBytes)
}

func main() {
	MakeKdbDirIfNeeded()

	// Initialize the repo
	repo, err := database.NewRepository()
	if err != nil {
		panic(err)
	}

	defer repo.Close()

	err = repo.Connect()
	if err != nil {
		panic(err)
	}

	err = repo.Init(EMBEDDING_MODEL_DIMENSIONS)
	if err != nil {
		panic(err)
	}

	var args InputArgs
	var rootCmd = &cobra.Command{
		Use:   "kdb",
		Short: "A knowledge database available as a command line tool",
		Run: func(cmd *cobra.Command, _ []string) {
			processInput(&args, repo)
		},
	}

	rootCmd.Flags().BoolVarP(&args.Stdin, "stdin", "s", false, "Read from stdin")
	rootCmd.Flags().StringVarP(&args.Query, "query", "q", "", "Search query")
	rootCmd.Flags().StringVarP(&args.Embed, "embed", "e", "", "Text to embed")
	rootCmd.Flags().IntVarP(&args.List, "list", "l", 0, "List embedded texts")
	rootCmd.Flags().IntVarP(&args.Top, "top", "t", 0, "Sets the max number of results to return for query")
	rootCmd.Flags().BoolVarP(&args.Version, "version", "v", false, "Prints the version")

	// Execute the command
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func processInput(args *InputArgs, repo database.EmbeddingRepo) {

	if args.Version {
		fmt.Printf("Version: %s-%s\n", Version, GitHash)
		os.Exit(0)
	}

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
		content := strings.TrimSpace(readStdInput())
		performInsert(content, repo)
	} else if embed != "" {
		performInsert(embed, repo)
	} else if query != "" {
		performQuery(query, args.Top, repo)
	} else if args.List > 0 {
		performList(args.List, repo)
	} else {
		fmt.Println("No action specified.")
	}
}

func performList(limit int, repo database.EmbeddingRepo) {
	items, err := repo.List(limit)
	if err != nil {
		log.Fatalln(err)
	}

	for _, item := range items {
		fmt.Printf("%s\t%s\t%s", item.Id, item.CreatedAt, item.Content)
	}
}

func performQuery(query string, top int, repo database.EmbeddingRepo) {
	eb, err := Embedd(query, EMBEDDING_MODEL_NAME)
	if err != nil {
		log.Fatalln(err)
	}

	items, err := repo.Query(eb, top)
	if err != nil {
		log.Fatalln(err)
	}

	for _, item := range items {
		fmt.Println(item.Content)
	}
}

func performInsert(content string, repo database.EmbeddingRepo) {
	eb, err := Embedd(content, EMBEDDING_MODEL_NAME)
	if err != nil {
		log.Fatalln(err)
	}

	result, err := repo.Insert(content, eb)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("Inserted item: %s\n", result.Id)
}
