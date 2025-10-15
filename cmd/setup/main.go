package main

import (
	"demo/internal/database"
	"fmt"
	"log"
)

// This program connects to the database and creates the necessary tables.
func main() {
	// Connect to the database using the function from your database package
	database.Connect()

	// Check if the connection is successful
	if database.DB == nil {
		log.Fatal("Failed to connect to the database.")
	}
	defer database.DB.Close()

	fmt.Println("Successfully connected to the database.")

	// Run the setup function to create tables
	if err := database.SetupDatabase(); err != nil {
		log.Fatalf("Failed to set up database: %v", err)
	}
}
