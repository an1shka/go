package main

import (
	"log"
	"net/http"

	"main/internal/handlers"
	"main/internal/store"
)

func main() {
	// Инициализируем хранилище данных и обработчики
	st := store.NewStore()
	h := handlers.NewHandler(st)

	// Публичные маршруты (доступны всем)
	http.HandleFunc("/auth/register", h.Register)
	http.HandleFunc("/auth/login", h.Login)

	// Защищенные маршруты (требуется токен в Authorization header)
	http.HandleFunc("/me", h.AuthMiddleware(h.GetMe))
	http.HandleFunc("/students", h.AuthMiddleware(h.StudentsRouter))
	http.HandleFunc("/students/", h.AuthMiddleware(h.StudentsRouter))

	log.Println("Сервер SchoolHub запущен на http://localhost:8081")
	if err := http.ListenAndServe(":8081", nil); err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}