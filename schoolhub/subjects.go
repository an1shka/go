package main

import (
	"encoding/json"
	"net/http"
)

func SubjectsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodGet:
		subjectsMutex.Lock()
		defer subjectsMutex.Unlock()
		json.NewEncoder(w).Encode(subjects)

	case http.MethodPost:
		var s Subject
		if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
			http.Error(w, `{"error": "Invalid request body"}`, http.StatusBadRequest)
			return
		}

		teachersMutex.Lock()
		teacherExists := false
		for _, t := range teachers {
			if t.ID == s.TeacherID {
				teacherExists = true
				break
			}
		}
		teachersMutex.Unlock()

		if !teacherExists {
			http.Error(w, `{"error": "Teacher not found"}`, http.StatusBadRequest)
			return
		}

		subjectsMutex.Lock()
		s.ID = nextSubjectID
		nextSubjectID++
		subjects = append(subjects, s)
		subjectsMutex.Unlock()

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(s)

	default:
		http.Error(w, `{"error": "Method not allowed"}`, http.StatusMethodNotAllowed)
	}
}