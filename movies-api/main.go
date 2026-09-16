package main

import (
	"log"
	"net/http"
)

func main() {
	storage := NewMemoryStorage()

	// Начальный фильм для демонстрации (ID = 1)
	storage.Create(MovieInput{
		Title:    "The Shawshank Redemption",
		Year:     1994,
		Rating:   9.3,
		Director: "Frank Darabont",
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