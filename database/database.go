package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/marcboeker/go-duckdb"
)

func Connect(vendor string, dbname string, dims int) (*sql.DB, error) {
	db, err := sql.Open(vendor, dbname)
	if err != nil {
		return nil, err
	}

	initDb(db, dims)
	return db, nil
}

func initDb(db *sql.DB, dims int) {

	batch := []string{
		"SET autoinstall_known_extensions=true;",
		"SET autoload_known_extensions=true;",
		"INSTALL vss",
		"LOAD vss",
		"SET hnsw_enable_experimental_persistence=true;",
		fmt.Sprintf(`CREATE TABLE IF NOT EXISTS embeddings (
			content TEXT NOT NULL,
			vector FLOAT[%d],
			created_at TIMESTAMP DEFAULT current_timestamp);`, dims),
		"CREATE INDEX IF NOT EXISTS idx_hnsw_vector ON embeddings USING HNSW (vector);",
	}

	for _, sql := range batch {
		_, err := db.Exec(sql)
		if err != nil {
			panic(err)
		}
	}
}

func InsertEmbedding(content string, embeddings []float64, db *sql.DB) (bool, error) {
	// Prepare the SQL statement
	stmt, err := db.Prepare("INSERT INTO embeddings(content, vector) VALUES (?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	// Execute the statement
	_, err = stmt.Exec(content, stringify(embeddings))
	if err != nil {
		log.Fatal(err)
		return false, err
	}

	return true, nil
}

func ListEmbeddings(limit int, db *sql.DB) {
	var (
		content string
	)

	query := "SELECT content FROM embeddings ORDER BY created_at DESC limit ?"
	rows, err := db.Query(query, limit)

	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	for rows.Next() {
		_ = rows.Scan(&content)
		fmt.Println(content)
	}

	// Check for errors from iterating over rows
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
}

func FindNearestEmbeddings(vector []float64, top int, db *sql.DB) {
	var (
		content  string
		distance float64
	)

	// TODO: Use prepared statement when the support for arrays were added to duckdb
	query := fmt.Sprintf(
		"SELECT content, array_distance(vector, %s) as distance FROM embeddings ORDER BY distance LIMIT ?",
		stringifyWithType(vector))
	rows, err := db.Query(query, top)

	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	for rows.Next() {
		_ = rows.Scan(&content, &distance)
		// Distance is not being used for now
		fmt.Println(content)
	}

	// Check for errors from iterating over rows
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
}
