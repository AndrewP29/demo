package main

import (
	"fmt"
	"log"
	"net/http"
	"html/template"
)

func main() {
	http.HandleFunc("/", homeHandler)

	fmt.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("web/templates/home.html"))
	tmpl.Execute(w, nil)
}