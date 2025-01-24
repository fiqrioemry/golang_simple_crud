package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	// Create a new router
	router := mux.NewRouter()

	// Placeholder: Add your routes here
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to the Job Portal API!"))
	}).Methods("GET")

	// Start the server
	log.Println("Starting server on :8080...")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
