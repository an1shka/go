package main

import (
    "fmt"
    "net/http"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, "Главная страница")
}

func aboutHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, "О нас")
}

func contactsHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, "Контакты")
}

func productsHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, "Товары")
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, "Привет!")
}

func main() {
    http.HandleFunc("/", homeHandler)
    http.HandleFunc("/about", aboutHandler)
    http.HandleFunc("/contacts", contactsHandler)
    http.HandleFunc("/products", productsHandler)
    http.HandleFunc("/hello", helloHandler)

    fmt.Println("Сервер запущен на http://localhost:8080")

    http.ListenAndServe(":8080", nil)
}