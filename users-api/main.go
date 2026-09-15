package main
import (
	"fmt"
	"net/http"
	"encoding/json"
)

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

var users = []User{
 	{
 		ID: 1,
 		Name: "Alex",
 		Age: 20,
 	},
 	{
 		ID: 2,
 		Name: "Anna",
 		Age: 22,
 	},
}


func main() {
	fmt.Println("Server started on http://127.0.0.1:8080")
	// http.HandleFunc("GET /users", getUsers)
	http.HandleFunc("GET /users/{id}", getUserByID)
 	err := http.ListenAndServe("127.0.0.1:8080", nil)
 	if err != nil {
 		fmt.Println("Server error:", err)
 	}
}

// func homeHandler(w http.ResponseWriter, r *http.Request) {
//  	fmt.Fprintln(w, "Users API is working")
// }

// func homeHandler(w http.ResponseWriter, r *http.Request) {
//  	fmt.Println("Method:", r.Method)
//  	fmt.Println("Path:", r.URL.Path)
//  	fmt.Fprintln(w, "Users API is working")
// }

func getUsers(w http.ResponseWriter, r *http.Request) {
 	w.Header().Set("Content-Type", "application/json")
 	err := json.NewEncoder(w).Encode(users)
 	if err != nil {
 		http.Error(
 			w,
 			"Failed to encode users",
 			http.StatusInternalServerError,
 		)
 	}
}