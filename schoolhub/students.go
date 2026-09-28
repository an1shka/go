package main

import (
	"net/http"
	"strings"
)

func StudentsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		className := r.URL.Query().Get("class_name")
		if className == "" {
			className = r.URL.Query().Get("class")
		}
		studentsMutex.RLock()
		result := make([]Student, 0)
		for _, student := range students {
			if className == "" || strings.EqualFold(student.ClassName, className) {
				result = append(result, student)
			}
		}
		studentsMutex.RUnlock()
		writeJSON(w, http.StatusOK, result)
	case http.MethodPost:
		var student Student
		if !decodeJSON(w, r, &student) {
			return
		}
		if strings.TrimSpace(student.FullName) == "" || strings.TrimSpace(student.Email) == "" || strings.TrimSpace(student.ClassName) == "" {
			writeError(w, http.StatusBadRequest, "full_name, email and class_name are required")
			return
		}
		studentsMutex.Lock()
		student.ID = nextStudentID
		nextStudentID++
		students = append(students, student)
		studentsMutex.Unlock()
		writeJSON(w, http.StatusCreated, student)
	default:
		methodNotAllowed(w)
	}
}

func StudentByIDHandler(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(r.URL.Path, "/students/")
	if !ok {
		writeError(w, http.StatusBadRequest, "invalid student ID")
		return
	}
	studentsMutex.Lock()
	defer studentsMutex.Unlock()
	for i, student := range students {
		if student.ID != id {
			continue
		}
		switch r.Method {
		case http.MethodGet:
			writeJSON(w, http.StatusOK, student)
		case http.MethodPut:
			var updated Student
			if !decodeJSON(w, r, &updated) {
				return
			}
			if strings.TrimSpace(updated.FullName) == "" || strings.TrimSpace(updated.Email) == "" || strings.TrimSpace(updated.ClassName) == "" {
				writeError(w, http.StatusBadRequest, "full_name, email and class_name are required")
				return
			}
			updated.ID = id
			students[i] = updated
			writeJSON(w, http.StatusOK, updated)
		case http.MethodDelete:
			students = append(students[:i], students[i+1:]...)
			writeJSON(w, http.StatusOK, map[string]string{"message": "Student deleted"})
		default:
			methodNotAllowed(w)
		}
		return
	}
	writeError(w, http.StatusNotFound, "student not found")
}

func studentExists(id int) bool {
	studentsMutex.RLock()
	defer studentsMutex.RUnlock()
	for _, student := range students {
		if student.ID == id {
			return true
		}
	}
	return false
}

func subjectExists(id int) bool {
	subjectsMutex.RLock()
	defer subjectsMutex.RUnlock()
	for _, subject := range subjects {
		if subject.ID == id {
			return true
		}
	}
	return false
}
