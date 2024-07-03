package main

import (
	"log"
	"os"
	"path/filepath"
)

func MakeKdbDirIfNeeded() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	kdbDir := filepath.Join(homeDir, ".kdb")

	if _, err := os.Stat(kdbDir); os.IsNotExist(err) {
		if err := os.Mkdir(kdbDir, 0755); err != nil {
			log.Fatal(err)
		}
	}
}
