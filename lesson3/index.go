package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

type Item struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

var (
	items  = make(map[int]Item)
	nextID = 1
	mu     sync.Mutex // защита от гонки данных при параллельных запросах
)

func main() {
	// Настраиваем маршруты
	http.HandleFunc("/api/items", handleItems)      // GET (все) и POST
	http.HandleFunc("/api/items/", handleItemByID)  // PUT, PATCH, DELETE (по ID)
	http.HandleFunc("/", handleIndex)               // Раздача фронтенда

	fmt.Println("Сервер запущен на http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

// Обработчик для /api/items/
func handleItems(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodGet:
		mu.Lock()
		list := make([]Item, 0, len(items))
		for _, item := range items {
			list = append(list, item)
		}
		mu.Unlock()
		json.NewEncoder(w).Encode(list)

	case http.MethodPost:
		var item Item
		if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
			http.Error(w, "Некорректный JSON", http.StatusBadRequest)
			return
		}
		mu.Lock()
		item.ID = nextID
		nextID++
		items[item.ID] = item
		mu.Unlock()

		w.WriteHeader(http.StatusCreated) // 201 Created
		json.NewEncoder(w).Encode(item)

	default:
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
	}
}

// Задания 1 и 2: Обработчики PUT, PATCH, DELETE по ID (/api/items/{id})
func handleItemByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Получаем ID из URL (/api/items/1 -> "1")
	idStr := strings.TrimPrefix(r.URL.Path, "/api/items/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Некорректный ID", http.StatusBadRequest)
		return
	}

	mu.Lock()
	existingItem, exists := items[id]
	mu.Unlock()

	if !exists {
		http.Error(w, "Элемент не найден", http.StatusNotFound) // 404
		return
	}

	switch r.Method {

	// 1. Метод PUT — Полная замена объекта
	case http.MethodPut:
		var newItem Item
		if err := json.NewDecoder(r.Body).Decode(&newItem); err != nil {
			http.Error(w, "Некорректный JSON", http.StatusBadRequest)
			return
		}
		newItem.ID = id // Задаем старый ID

		mu.Lock()
		items[id] = newItem
		mu.Unlock()

		json.NewEncoder(w).Encode(newItem)

	// 1. Метод PATCH — Частичное обновление объекта
	case http.MethodPatch:
		// В реальных приложениях декодируют в map или struct с указателями,
		// чтобы менять только переданные поля.
		var updates map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
			http.Error(w, "Некорректный JSON", http.StatusBadRequest)
			return
		}

		mu.Lock()
		if newName, ok := updates["name"].(string); ok {
			existingItem.Name = newName
		}
		items[id] = existingItem
		mu.Unlock()

		json.NewEncoder(w).Encode(existingItem)

	// 2. Метод DELETE — Удаление ресурса
	case http.MethodDelete:
		mu.Lock()
		delete(items, id)
		mu.Unlock()

		w.WriteHeader(http.StatusNoContent) // 204 No Content

	default:
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
	}
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}