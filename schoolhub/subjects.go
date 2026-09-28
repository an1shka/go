package main

import (
	"net/http"
	"strings"
)

func SubjectsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		subjectsMutex.RLock()
		result := append([]Subject(nil), subjects...)
		subjectsMutex.RUnlock()
		writeJSON(w, http.StatusOK, result)
	case http.MethodPost:
		var subject Subject
		if !decodeJSON(w, r, &subject) {
			return
		}
		if strings.TrimSpace(subject.Title) == "" {
			writeError(w, http.StatusBadRequest, "title is required")
			return
		}
		if !teacherExists(subject.TeacherID) {
			writeError(w, http.StatusBadRequest, "teacher not found")
			return
		}
		subjectsMutex.Lock()
		subject.ID = nextSubjectID
		nextSubjectID++
		subjects = append(subjects, subject)
		subjectsMutex.Unlock()
		writeJSON(w, http.StatusCreated, subject)
	default:
		methodNotAllowed(w)
	}
}

func SubjectByIDHandler(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(r.URL.Path, "/subjects/")
	if !ok {
		writeError(w, http.StatusBadRequest, "invalid subject ID")
		return
	}
	subjectsMutex.Lock()
	defer subjectsMutex.Unlock()
	for i, subject := range subjects {
		if subject.ID != id {
			continue
		}
		switch r.Method {
		case http.MethodGet:
			writeJSON(w, http.StatusOK, subject)
		case http.MethodPut:
			var updated Subject
			if !decodeJSON(w, r, &updated) {
				return
			}
			if strings.TrimSpace(updated.Title) == "" {
				writeError(w, http.StatusBadRequest, "title is required")
				return
			}
			if !teacherExists(updated.TeacherID) {
				writeError(w, http.StatusBadRequest, "teacher not found")
				return
			}
			updated.ID = id
			subjects[i] = updated
			writeJSON(w, http.StatusOK, updated)
		case http.MethodDelete:
			subjects = append(subjects[:i], subjects[i+1:]...)
			writeJSON(w, http.StatusOK, map[string]string{"message": "Subject deleted"})
		default:
			methodNotAllowed(w)
		}
		return
	}
	writeError(w, http.StatusNotFound, "subject not found")
}

func teacherExists(id int) bool {
	teachersMutex.RLock()
	defer teachersMutex.RUnlock()
	for _, teacher := range teachers {
		if teacher.ID == id {
			return true
		}
	}
	return false
}
