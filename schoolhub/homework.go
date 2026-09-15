package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"
)

func HomeworkHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodGet:
		subjectIDStr := r.URL.Query().Get("subject_id")
		overdueStr := r.URL.Query().Get("overdue")

		homeworksMutex.Lock()
		defer homeworksMutex.Unlock()

		var result []Homework
		now := time.Now()

		for _, hw := range homeworks {
			match := true

			if subjectIDStr != "" {
				subjID, _ := strconv.Atoi(subjectIDStr)
				if hw.SubjectID != subjID {
					match = false
				}
			}

			if overdueStr == "true" {
				if !hw.Deadline.Before(now) {
					match = false
				}
			}

			if match {
				result = append(result, hw)
			}
		}

		json.NewEncoder(w).Encode(result)

	case http.MethodPost:
		var hw Homework
		if err := json.NewDecoder(r.Body).Decode(&hw); err != nil {
			http.Error(w, `{"error": "Invalid request body"}`, http.StatusBadRequest)
			return
		}

		if hw.IssuedAt.IsZero() {
			hw.IssuedAt = time.Now()
		}

		if hw.Deadline.Before(hw.IssuedAt) {
			http.Error(w, `{"error": "Deadline cannot be earlier than issued date"}`, http.StatusBadRequest)
			return
		}

		homeworksMutex.Lock()
		hw.ID = nextHomeworkID
		nextHomeworkID++
		homeworks = append(homeworks, hw)
		homeworksMutex.Unlock()

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(hw)

	default:
		http.Error(w, `{"error": "Method not allowed"}`, http.StatusMethodNotAllowed)
	}
}