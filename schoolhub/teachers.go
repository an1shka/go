package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

func TeachersHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		teachersMutex.RLock()
		result := append([]Teacher(nil), teachers...)
		teachersMutex.RUnlock()
		writeJSON(w, http.StatusOK, result)
	case http.MethodPost:
		var teacher Teacher
		if !decodeJSON(w, r, &teacher) {
			return
		}
		if strings.TrimSpace(teacher.FullName) == "" || strings.TrimSpace(teacher.Email) == "" {
			writeError(w, http.StatusBadRequest, "full_name and email are required")
			return
		}
		teachersMutex.Lock()
		teacher.ID = nextTeacherID
		nextTeacherID++
		teachers = append(teachers, teacher)
		teachersMutex.Unlock()
		writeJSON(w, http.StatusCreated, teacher)
	default:
		methodNotAllowed(w)
	}
}

func TeacherByIDHandler(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(r.URL.Path, "/teachers/")
	if !ok {
		writeError(w, http.StatusBadRequest, "invalid teacher ID")
		return
	}
	teachersMutex.Lock()
	defer teachersMutex.Unlock()
	for i, teacher := range teachers {
		if teacher.ID != id {
			continue
		}
		switch r.Method {
		case http.MethodGet:
			writeJSON(w, http.StatusOK, teacher)
		case http.MethodPut:
			var updated Teacher
			if !decodeJSON(w, r, &updated) {
				return
			}
			if strings.TrimSpace(updated.FullName) == "" || strings.TrimSpace(updated.Email) == "" {
				writeError(w, http.StatusBadRequest, "full_name and email are required")
				return
			}
			updated.ID = id
			teachers[i] = updated
			writeJSON(w, http.StatusOK, updated)
		case http.MethodDelete:
			teachers = append(teachers[:i], teachers[i+1:]...)
			writeJSON(w, http.StatusOK, map[string]string{"message": "Teacher deleted"})
		default:
			methodNotAllowed(w)
		}
		return
	}
	writeError(w, http.StatusNotFound, "teacher not found")
}

func pathID(path, prefix string) (int, bool) {
	value := strings.Trim(strings.TrimPrefix(path, prefix), "/")
	if value == "" || strings.Contains(value, "/") {
		return 0, false
	}
	id, err := strconv.Atoi(value)
	return id, err == nil && id > 0
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}

func methodNotAllowed(w http.ResponseWriter) {
	writeError(w, http.StatusMethodNotAllowed, "method not allowed")
}
