package main

import (
	"fmt"
	"log"
	"net/http"

	"demo/internal/api"
	"demo/internal/database"
)

func main() {
	// Connect to the database
	database.Connect()
	if database.DB == nil {
		log.Fatal("Failed to connect to the database.")
	}
	defer database.DB.Close()
	fmt.Println("Successfully connected to the database.")

	// API endpoints
	http.HandleFunc("/api/signup", api.SignupHandler)
	http.HandleFunc("/api/login", api.LoginHandler)
	// Future API endpoints like /api/login would go here

	// Serve static files for the frontend
	fs := http.FileServer(http.Dir("./web/static"))
	http.Handle("/", fs)

	fmt.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
