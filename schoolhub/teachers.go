package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

func TeachersHandler(w wWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodGet:
		teachersMutex.Lock()
		defer teachersMutex.Unlock()
		json.NewEncoder(w).Encode(teachers)

	case http.MethodPost:
		var t Teacher
		if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
			http.Error(w, `{"error": "Invalid request body"}`, http.StatusBadRequest)
			return
		}

		teachersMutex.Lock()
		t.ID = nextTeacherID
		nextTeacherID++
		teachers = append(teachers, t)
		teachersMutex.Unlock()

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(t)

	default:
		http.Error(w, `{"error": "Method not allowed"}`, http.StatusMethodNotAllowed)
	}
}

type wWriter = http.ResponseWriter

func TeacherByIDHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	idStr := strings.TrimPrefix(r.URL.Path, "/teachers/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, `{"error": "Invalid ID"}`, http.StatusBadRequest)
		return
	}

	teachersMutex.Lock()
	defer teachersMutex.Unlock()

	for i, t := range teachers {
		if t.ID == id {
			if r.Method == http.MethodGet {
				json.NewEncoder(w).Encode(t)
				return
			} else if r.Method == http.MethodDelete {
				teachers = append(teachers[:i], teachers[i+1:]...)
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(map[string]string{"message": "Teacher deleted"})
				return
			}
		}
	}

	http.Error(w, `{"error": "Teacher not found"}`, http.StatusNotFound)
}