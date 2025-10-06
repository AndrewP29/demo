package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	//Whenever someone visits the home page, use this function
	// w is the response and r is the request
	// * indicates a pointer
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "My demo API")
	})

	fmt.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}