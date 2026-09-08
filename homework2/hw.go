package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Product struct {
	ID    int     `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

var products = []Product{
	{
		ID:    1,
		Name:  "Laptop",
		Price: 500000,
	},
	{
		ID:    2,
		Name:  "Mouse",
		Price: 12000,
	},
}

func getProducts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	err := json.NewEncoder(w).Encode(products)

	json.NewEncoder(w).Encode(products)

	if err != nil {
		http.Error(
			w,
			"Failed to encode products",
			http.StatusInternalServerError,
		)
	}
}

func createProduct(w http.ResponseWriter, r *http.Request) {
	var product Product

	err := json.NewDecoder(r.Body).Decode(&product)

	if err != nil {
		http.Error(
			w,
			"Invalid JSON",
			http.StatusBadRequest,
		)
		return
	}

	if product.Name == "" {
		http.Error(
			w,
			"Name is required",
			http.StatusBadRequest,
		)
		return
	}

	if product.Price <= 0 {
		http.Error(
			w,
			"Price must be greater than 0",
			http.StatusBadRequest,
		)
		return
	}

	product.ID = len(products) + 1

	products = append(products, product)

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(product)
}

func productsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getProducts(w, r)

	case http.MethodPost:
		createProduct(w, r)

	default:
		http.Error(
			w,
			"Method not allowed",
			http.StatusMethodNotAllowed,
		)
	}
}

func main() {
	http.HandleFunc("/products", productsHandler)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}

		http.Redirect(w, r, "/products", http.StatusFound)
	})

	fmt.Println("Server started on http://localhost:8080")

	err := http.ListenAndServe(":8080", nil)

	if err != nil {
		fmt.Println("Server error:", err)
	}
}
