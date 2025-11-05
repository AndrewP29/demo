package main

import (
	"fmt"
	"log"
	"net/http"

	"demo/internal/api"
	"demo/internal/database"
	"demo/internal/session"
)

func main() {
	// Connect to the database
	db, err := database.NewDB()
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	// This runs just before main exits
	defer db.Close()

	// A pointer to a struct where the DB field has value db
	store := &database.DBStore{DB: db}

	// Initialize the in-memory session store
	sessStore := session.NewMemoryStore()

	server := &api.Server{Store: store, SessionStore: sessStore}

	fmt.Println("Successfully connected to the database.")

	// API endpoints
	http.HandleFunc("/api/signup", server.SignupHandler)
	http.HandleFunc("/api/login", server.LoginHandler)

	// Serve static files for the frontend
	fs := http.FileServer(http.Dir("./web/static"))
	http.Handle("/", fs)

	fmt.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
