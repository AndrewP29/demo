package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func Connect() {
	var err error

	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	host := os.Getenv("DB_HOST")

	if user == "" || password == "" || dbname == "" {
		log.Fatal("Missing required database environment variables")
	}

	if host == "" {
		host = "localhost"
	}

	connStr := fmt.Sprintf("user=%s password=%s dbname=%s host=%s sslmode=disable",
		user, password, dbname, host)

	DB, err = sql.Open("postgres", connStr)

	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}
}