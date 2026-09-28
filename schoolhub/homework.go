package main

import (
	"net/http"
	"time"
)

func HomeworkHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		subjectID, err := optionalQueryID(r, "subject_id")
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid subject_id")
			return
		}
		overdue := r.URL.Query().Get("overdue") == "true"
		now := time.Now()
		homeworksMutex.RLock()
		result := make([]Homework, 0)
		for _, homework := range homeworks {
			if subjectID != nil && homework.SubjectID != *subjectID {
				continue
			}
			if overdue && !homework.Deadline.Before(now) {
				continue
			}
			result = append(result, homework)
		}
		homeworksMutex.RUnlock()
		writeJSON(w, http.StatusOK, result)
	case http.MethodPost:
		var homework Homework
		if !decodeJSON(w, r, &homework) {
			return
		}
		if err := validateHomework(homework); err != nil {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		if homework.IssuedAt.IsZero() {
			homework.IssuedAt = time.Now().UTC()
		}
		homeworksMutex.Lock()
		homework.ID = nextHomeworkID
		nextHomeworkID++
		homeworks = append(homeworks, homework)
		homeworksMutex.Unlock()
		writeJSON(w, http.StatusCreated, homework)
	default:
		methodNotAllowed(w)
	}
}

func HomeworkByIDHandler(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(r.URL.Path, "/homework/")
	if !ok {
		writeError(w, http.StatusBadRequest, "invalid homework ID")
		return
	}
	homeworksMutex.Lock()
	defer homeworksMutex.Unlock()
	for i, homework := range homeworks {
		if homework.ID != id {
			continue
		}
		switch r.Method {
		case http.MethodGet:
			writeJSON(w, http.StatusOK, homework)
		case http.MethodPut:
			var updated Homework
			if !decodeJSON(w, r, &updated) {
				return
			}
			if err := validateHomework(updated); err != nil {
				writeError(w, http.StatusBadRequest, err.Error())
				return
			}
			if updated.IssuedAt.IsZero() {
				updated.IssuedAt = homework.IssuedAt
			}
			updated.ID = id
			homeworks[i] = updated
			writeJSON(w, http.StatusOK, updated)
		case http.MethodDelete:
			homeworks = append(homeworks[:i], homeworks[i+1:]...)
			writeJSON(w, http.StatusOK, map[string]string{"message": "Homework deleted"})
		default:
			methodNotAllowed(w)
		}
		return
	}
	writeError(w, http.StatusNotFound, "homework not found")
}

func validateHomework(homework Homework) error {
	if homework.Title == "" {
		return &validationError{"title is required"}
	}
	if !subjectExists(homework.SubjectID) {
		return &validationError{"subject not found"}
	}
	issuedAt := homework.IssuedAt
	if issuedAt.IsZero() {
		issuedAt = time.Now()
	}
	if homework.Deadline.IsZero() {
		return &validationError{"deadline is required"}
	}
	if homework.Deadline.Before(issuedAt) {
		return &validationError{"deadline cannot be earlier than issued_at"}
	}
	return nil
}
