package main

import (
	"fmt"
	"log"
	"net/http"

	"demo/internal/api"
	"demo/internal/database" // Import the database package
)

func main() {
	// Connect to the database
	database.Connect()
	if database.DB == nil {
		log.Fatal("Failed to connect to the database.")
	}
	defer database.DB.Close()
	fmt.Println("Successfully connected to the database.")

	// Use handlers from the 'api' package
	http.HandleFunc("/", api.HomeHandler)
	http.HandleFunc("/login", api.LoginHandler)
	http.HandleFunc("/signup", api.SignupHandler)
	http.HandleFunc("/dashboard", api.DashboardHandler)

	fmt.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
