package main

import (
	"fmt"
	"net/http"
	"encoding/json"
	"strconv"
)
//хендлеры
func healthHandler(w http.ResponseWriter, r *http.Request){
	fmt.Fprintln(w, "Users API is working")
}
func getUsers(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")

	err := json.NewEncoder(w).Encode(users)

	if err != nil{
		http.Error(w,"Failed to encode users",http.StatusInternalServerError)
	}
}

id,err := strconv.Atoi(idString)

func findUserByID(id int)(*User error){
	for i := renge users{
		if users[i].ID == id{
			return &users[i], nil
		}
	}
	return nil, errors.New("user not found")
}

func getUserById

// Структура пользователя
type User struct{
	ID int `json:"id"`
	Name string `json:"name"`
	Age int `json:"age"`
}
// Слайс юзеров
var users = []User{
		{
			ID:1,
			Name: "Alex",
			Age: 20,
		},
		{
			ID:2,
			Name: "John",
			Age: 21,
		},
	}

func main(){
	//маршруты
	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("GET /users",getUsers)
	http.HandleFunc("GET /users/{id}",)
	// Запуск сервера
	fmt.Println("Server started on http://localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil{
		fmt.Println("Server error", err)
	}
}