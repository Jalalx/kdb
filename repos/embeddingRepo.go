package repos

import (
	"time"

	"github.com/google/uuid"
)

type EmbeddingListItem struct {
	Id        uuid.UUID
	Content   string
	CreatedAt time.Time
}

type EmbeddingQueryItem struct {
	Id        string
	Content   string
	Distance  float64
	CreatedAt time.Time
}

func NewRepository() (EmbeddingRepo, error) {
	return NewDuckDbEmbeddingRepo(), nil
}

type EmbeddingRepo interface {
	Connect() error

	Init(dims int) error

	Insert(content string, embeddings []float64) (EmbeddingListItem, error)

	List(limit int) ([]EmbeddingListItem, error)

	Query(vector []float64, top int) ([]EmbeddingQueryItem, error)

	Close() error
}
