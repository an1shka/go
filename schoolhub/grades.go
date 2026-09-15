package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func GradesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodGet:
		gradesMutex.Lock()
		defer gradesMutex.Unlock()
		json.NewEncoder(w).Encode(grades)

	case http.MethodPost:
		var g Grade
		if err := json.NewDecoder(r.Body).Decode(&g); err != nil {
			http.Error(w, `{"error": "Invalid request body"}`, http.StatusBadRequest)
			return
		}

		if g.Value < 1 || g.Value > 5 {
			http.Error(w, `{"error": "Grade value must be between 1 and 5"}`, http.StatusBadRequest)
			return
		}

		if g.Date.IsZero() {
			g.Date = time.Now()
		}

		gradesMutex.Lock()
		g.ID = nextGradeID
		nextGradeID++
		grades = append(grades, g)
		gradesMutex.Unlock()

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(g)

	default:
		http.Error(w, `{"error": "Method not allowed"}`, http.StatusMethodNotAllowed)
	}
}

func StudentGradesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	path := r.URL.Path
	parts := strings.Split(path, "/")
	if len(parts) < 4 {
		http.Error(w, `{"error": "Invalid URL"}`, http.StatusBadRequest)
		return
	}

	studentID, err := strconv.Atoi(parts[2])
	if err != nil {
		http.Error(w, `{"error": "Invalid student ID"}`, http.StatusBadRequest)
		return
	}

	subjectIDQuery := r.URL.Query().Get("subject_id")

	gradesMutex.Lock()
	defer gradesMutex.Unlock()

	var result []Grade
	for _, g := range grades {
		if g.StudentID == studentID {
			if subjectIDQuery != "" {
				subjID, _ := strconv.Atoi(subjectIDQuery)
				if g.SubjectID == subjID {
					result = append(result, g)
				}
			} else {
				result = append(result, g)
			}
		}
	}

	json.NewEncoder(w).Encode(result)
}