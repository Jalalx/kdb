package repos

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jalalx/kdb/config"
)

func NewDuckDbEmbeddingRepo(cfg *config.KdbStorageConfig) *EmbeddingDuckDbRepo {
	return &EmbeddingDuckDbRepo{
		cfg: cfg,
	}
}

type EmbeddingDuckDbRepo struct {
	cfg *config.KdbStorageConfig
}

func (repo *EmbeddingDuckDbRepo) connect() (*sql.DB, error) {
	providerName := strings.ToLower(repo.cfg.Provider)
	db, err := sql.Open(providerName, repo.cfg.URL)
	if err != nil {
		return nil, err
	}

	db.Exec("SET autoinstall_known_extensions=true;")
	db.Exec("SET autoload_known_extensions=true;")
	db.Exec("INSTALL vss")
	db.Exec("LOAD vss")
	db.Exec("SET hnsw_enable_experimental_persistence=true;")

	return db, nil
}

func (repo *EmbeddingDuckDbRepo) Init(dims int) error {
	batch := []string{
		fmt.Sprintf(`CREATE TABLE IF NOT EXISTS embeddings (
			id uuid DEFAULT gen_random_uuid(),
			content TEXT NOT NULL,
			vector FLOAT[%d],
			created_at TIMESTAMP DEFAULT current_timestamp);`, dims),
		"CREATE INDEX IF NOT EXISTS idx_hnsw_vector ON embeddings USING HNSW (vector);",
	}

	db, err := repo.connect()
	if err != nil {
		panic(err)
	}

	defer db.Close()
	for _, sql := range batch {
		_, err := db.Exec(sql)
		if err != nil {
			return err
		}
	}

	return nil
}

func (repo *EmbeddingDuckDbRepo) Insert(content string, embeddings []float64) (EmbeddingListItem, error) {

	db, err := repo.connect()
	if err != nil {
		panic(err)
	}

	defer db.Close()

	// Prepare the SQL statement
	stmt, err := db.Prepare("INSERT INTO embeddings(id, content, vector, created_at) VALUES (?, ?, ?, ?)")
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

	db, err := repo.connect()
	if err != nil {
		panic(err)
	}

	defer db.Close()

	query := "SELECT id, content, created_at FROM embeddings ORDER BY created_at DESC limit ?"
	rows, err := db.Query(query, limit)

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

	db, err := repo.connect()
	if err != nil {
		panic(err)
	}

	defer db.Close()

	// TODO: Use prepared statement when the support for arrays were added to duckdb
	query := fmt.Sprintf(
		"SELECT id, content, array_distance(vector, %s) as distance, created_at FROM embeddings ORDER BY distance LIMIT ?",
		stringifyWithType(vector))
	rows, err := db.Query(query, top)

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

	db, err := repo.connect()
	if err != nil {
		panic(err)
	}

	defer db.Close()

	query := "DELETE FROM embeddings WHERE id = ?"
	result, err := db.Exec(query, id)
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
