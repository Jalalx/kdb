package main

import (
	"database/sql"
	"log"

	"github.com/jalalx/kdb/database"
	ollamaclient "github.com/xyproto/ollamaclient"
)

func getClientConfig() (*ollamaclient.Config, error) {
	oc := ollamaclient.NewWithModel(EMBEDDING_MODEL_NAME)
	oc.Verbose = false

	if err := oc.PullIfNeeded(); err != nil {
		return nil, err
	}

	return oc, nil
}

func Embedd(prompt string, model string, db *sql.DB) ([]float64, error) {
	oc, err := getClientConfig()
	if err != nil {
		return nil, err
	}

	embeddings, err := oc.Embeddings(prompt)
	if err != nil {
		return nil, err
	}

	return embeddings, nil
}

func Query(query string, top int, db *sql.DB) {
	embeddings, err := Embedd(query, EMBEDDING_MODEL_NAME, db)
	if err != nil {
		log.Fatalln(err)
	}

	database.FindNearestEmbeddings(embeddings, top, db)
}
