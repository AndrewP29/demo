package database

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

func NewDB() (*sql.DB, error) {
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")

	if user == "" || password == "" || dbname == "" {
		return nil, fmt.Errorf("missing required database environment variables")
	}

	if host == "" {
		host = "localhost"
	}

	if port == "" {
		port = "5432"
	}

	connStr := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
		user, password, dbname, host, port)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	// Read and execute the setup.sql file
	setupSQL, err := os.ReadFile("sql/setup.sql")
	if err != nil {
		return nil, fmt.Errorf("failed to read setup.sql: %w", err)
	}

	if _, err = db.Exec(string(setupSQL)); err != nil {
		return nil, fmt.Errorf("failed to execute setup.sql: %w", err)
	}

	return db, nil
}