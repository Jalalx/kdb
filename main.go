package main

import (
	"context"
	"log"
	"os"
	"path/filepath"

	_ "github.com/marcboeker/go-duckdb"
	"github.com/spf13/cobra"

	"github.com/jalalx/kdb/cmd"
	"github.com/jalalx/kdb/config"
	"github.com/jalalx/kdb/llms"
	"github.com/jalalx/kdb/repos"
)

const (
	KDB_DIR = ".kdb"
)

var (
	Version string
	GitHash string
)

func main() {
	cfg := getConfig()

	// Initialize the repo
	repo, err := repos.NewRepository(&cfg.Storage)
	if err != nil {
		panic(err)
	}

	defer repo.Close()

	llmProvider, err := llms.NewLlmProvider(&cfg.Embedding)
	if err != nil {
		panic(err)
	}

	err = repo.Connect()
	if err != nil {
		panic(err)
	}

	err = repo.Init(cfg.Embedding.Model.Dimensions)
	if err != nil {
		panic(err)
	}

	rootCmd := cmd.NewCLI(repo, llmProvider, Version, GitHash)
	cobra.CheckErr(rootCmd.ExecuteContext(context.Background()))
}

func getConfig() *config.KdbConfig {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	kdbDir := filepath.Join(homeDir, KDB_DIR)

	if _, err := os.Stat(kdbDir); os.IsNotExist(err) {
		if err := os.Mkdir(kdbDir, 0755); err != nil {
			log.Fatal(err)
		}
	}

	var cfg *config.KdbConfig
	configFile := filepath.Join(kdbDir, config.DEFAULT_FILENAME)
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		cfg = config.NewDefaultConfig()
		config.SaveConfig(cfg, configFile)
	} else {
		cfg, err = config.LoadConfig(configFile)
		if err != nil {
			log.Fatal(err)
		}
	}

	return cfg
}
