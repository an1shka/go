package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

type Handler struct {
	storage *MemoryStorage
}

func NewHandler(storage *MemoryStorage) *Handler {
	return &Handler{storage: storage}
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}

func (h *Handler) MovieRoutes(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/movies")
	path = strings.Trim(path, "/")

	if path == "" {
		switch r.Method {
		case http.MethodGet:
			h.GetMovies(w, r)
		case http.MethodPost:
			h.CreateMovie(w, r)
		default:
			writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		}
		return
	}

	// Валидация ID: если передана строка не из цифр (например, /movies/abc),
	// возвращаем 400 Bad Request
	id, err := strconv.Atoi(path)
	if err != nil {
		writeError(w, http.StatusBadRequest, "ID must be a number")
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.GetMovieByID(w, r, id)
	case http.MethodPut:
		h.UpdateMovie(w, r, id)
	case http.MethodDelete:
		h.DeleteMovie(w, r, id)
	default:
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

func (h *Handler) GetMovies(w http.ResponseWriter, r *http.Request) {
	movies := h.storage.GetAll()
	writeJSON(w, http.StatusOK, movies)
}

func (h *Handler) GetMovieByID(w http.ResponseWriter, r *http.Request, id int) {
	movie, err := h.storage.GetByID(id)
	if err != nil {
		writeError(w, http.StatusNotFound, "Movie not found")
		return
	}
	writeJSON(w, http.StatusOK, movie)
}

func (h *Handler) CreateMovie(w http.ResponseWriter, r *http.Request) {
	var input MovieInput

	// Ошибка парсинга JSON -> 400 Bad Request
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid JSON payload")
		return
	}

	// Ошибка валидации -> 400 Bad Request
	if err := validateMovie(input); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	movie := h.storage.Create(input)
	writeJSON(w, http.StatusCreated, movie)
}

func (h *Handler) UpdateMovie(w http.ResponseWriter, r *http.Request, id int) {
	var input MovieInput

	// Ошибка парсинга JSON -> 400 Bad Request
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid JSON payload")
		return
	}

	// Ошибка валидации -> 400 Bad Request
	if err := validateMovie(input); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	movie, err := h.storage.Update(id, input)
	if err != nil {
		writeError(w, http.StatusNotFound, "Movie not found")
		return
	}

	writeJSON(w, http.StatusOK, movie)
}

func (h *Handler) DeleteMovie(w http.ResponseWriter, r *http.Request, id int) {
	err := h.storage.Delete(id)
	if err != nil {
		writeError(w, http.StatusNotFound, "Movie not found")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}