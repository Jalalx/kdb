package repos

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
)

const (
	DUCKDB_PROVIDER = "duckdb"
	DUCKDB_PATH     = "~/.kdb/knowledgebase.ddb"
)

func NewDuckDbEmbeddingRepo() *EmbeddingDuckDbRepo {
	return &EmbeddingDuckDbRepo{}
}

type EmbeddingDuckDbRepo struct {
	db *sql.DB
}

func (repo *EmbeddingDuckDbRepo) Connect() error {
	db, err := sql.Open(DUCKDB_PROVIDER, DUCKDB_PATH)
	if err != nil {
		return err
	}

	repo.db = db
	return nil
}

func (repo *EmbeddingDuckDbRepo) Init(dims int) error {
	batch := []string{
		"SET autoinstall_known_extensions=true;",
		"SET autoload_known_extensions=true;",
		"INSTALL vss",
		"LOAD vss",
		"SET hnsw_enable_experimental_persistence=true;",
		fmt.Sprintf(`CREATE TABLE IF NOT EXISTS embeddings (
			id uuid DEFAULT gen_random_uuid(),
			content TEXT NOT NULL,
			vector FLOAT[%d],
			created_at TIMESTAMP DEFAULT current_timestamp);`, dims),
		"CREATE INDEX IF NOT EXISTS idx_hnsw_vector ON embeddings USING HNSW (vector);",
	}

	for _, sql := range batch {
		_, err := repo.db.Exec(sql)
		if err != nil {
			return err
		}
	}

	return nil
}

func (repo *EmbeddingDuckDbRepo) Insert(content string, embeddings []float64) (EmbeddingListItem, error) {
	// Prepare the SQL statement
	stmt, err := repo.db.Prepare("INSERT INTO embeddings(id, content, vector, created_at) VALUES (?, ?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	item := EmbeddingListItem{}
	uuid := uuid.New()
	utcNow := time.Now().UTC()
	// Execute the statement
	_, err = stmt.Exec(uuid, content, stringify(embeddings), utcNow)
	if err != nil {
		log.Fatal(err)
		return item, err
	}

	item.Id = uuid
	item.Content = content
	item.CreatedAt = utcNow

	return item, nil
}

func (repo *EmbeddingDuckDbRepo) List(limit int) ([]EmbeddingListItem, error) {
	items := []EmbeddingListItem{}

	query := "SELECT id, content, created_at FROM embeddings ORDER BY created_at DESC limit ?"
	rows, err := repo.db.Query(query, limit)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		item := EmbeddingListItem{}
		_ = rows.Scan(&item.Id, &item.Content, &item.CreatedAt)
		items = append(items, item)
	}

	// Check for errors from iterating over rows
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

func (repo *EmbeddingDuckDbRepo) Query(vector []float64, top int) ([]EmbeddingQueryItem, error) {
	items := []EmbeddingQueryItem{}

	// TODO: Use prepared statement when the support for arrays were added to duckdb
	query := fmt.Sprintf(
		"SELECT id, content, array_distance(vector, %s) as distance, created_at FROM embeddings ORDER BY distance LIMIT ?",
		stringifyWithType(vector))
	rows, err := repo.db.Query(query, top)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		item := EmbeddingQueryItem{}
		_ = rows.Scan(&item.Id, &item.Content, &item.Distance, &item.CreatedAt)
		items = append(items, item)
		// fmt.Println(content)
	}

	// Check for errors from iterating over rows
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

func (repo *EmbeddingDuckDbRepo) Delete(id uuid.UUID) (bool, error) {
	query := "DELETE FROM embeddings WHERE id = ?"
	result, err := repo.db.Exec(query, id)
	if err != nil {
		return false, err
	}

	affectedRows, err := result.RowsAffected()
	if err != nil {
		return false, err
	}

	if affectedRows > 0 {
		return true, nil
	}

	return false, nil
}

func (repo *EmbeddingDuckDbRepo) Close() error {
	return repo.db.Close()
}
