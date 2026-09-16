package main

import (
	"log"
	"net/http"
)

func main() {
	storage := NewMemoryStorage()
	
	storage.Create(MovieInput{
		Title:    "Inception",
		Year:     2010,
		Rating:   8.8,
		Director: "Christopher Nolan",
	})

	handler := NewHandler(storage)

	mux := http.NewServeMux()
	mux.HandleFunc("/movies", handler.MovieRoutes)
	mux.HandleFunc("/movies/", handler.MovieRoutes)

	log.Println("Server is running on http://localhost:8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}