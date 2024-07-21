package repos

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jalalx/kdb/config"
)

type EmbeddingListItem struct {
	Id        uuid.UUID
	Content   string
	CreatedAt time.Time
}

type EmbeddingQueryItem struct {
	Id        uuid.UUID
	Content   string
	Distance  float64
	CreatedAt time.Time
}

func NewRepository(cfg *config.KdbStorageConfig) (EmbeddingRepo, error) {
	if strings.ToLower(cfg.Provider) == "duckdb" {
		return NewDuckDbEmbeddingRepo(), nil
	}

	return nil, fmt.Errorf("database provider not supported: %s", cfg.Provider)
}

type EmbeddingRepo interface {
	Connect() error

	Init(dims int) error

	Insert(content string, embeddings []float64) (EmbeddingListItem, error)

	Delete(id uuid.UUID) (bool, error)

	List(limit int) ([]EmbeddingListItem, error)

	Query(vector []float64, top int) ([]EmbeddingQueryItem, error)

	Close() error
}
