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

	if user == "" || password == "" || dbname == "" {
		return nil, fmt.Errorf("missing required database environment variables")
	}

	if host == "" {
		host = "localhost"
	}

	connStr := fmt.Sprintf("user=%s password=%s dbname=%s host=%s sslmode=disable",
		user, password, dbname, host)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}