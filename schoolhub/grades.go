package main

import (
	"net/http"
	"strconv"
	"strings"
	"time"
)

func GradesHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		studentID, err := optionalQueryID(r, "student_id")
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid student_id")
			return
		}
		subjectID, err := optionalQueryID(r, "subject_id")
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid subject_id")
			return
		}
		gradesMutex.RLock()
		result := make([]Grade, 0)
		for _, grade := range grades {
			if studentID != nil && grade.StudentID != *studentID {
				continue
			}
			if subjectID != nil && grade.SubjectID != *subjectID {
				continue
			}
			result = append(result, grade)
		}
		gradesMutex.RUnlock()
		writeJSON(w, http.StatusOK, result)
	case http.MethodPost:
		var grade Grade
		if !decodeJSON(w, r, &grade) {
			return
		}
		if err := validateGrade(grade); err != nil {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		if grade.Date.IsZero() {
			grade.Date = time.Now().UTC()
		}
		gradesMutex.Lock()
		grade.ID = nextGradeID
		nextGradeID++
		grades = append(grades, grade)
		gradesMutex.Unlock()
		writeJSON(w, http.StatusCreated, grade)
	default:
		methodNotAllowed(w)
	}
}

func GradeByIDHandler(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(r.URL.Path, "/grades/")
	if !ok {
		writeError(w, http.StatusBadRequest, "invalid grade ID")
		return
	}
	gradesMutex.Lock()
	defer gradesMutex.Unlock()
	for i, grade := range grades {
		if grade.ID != id {
			continue
		}
		switch r.Method {
		case http.MethodGet:
			writeJSON(w, http.StatusOK, grade)
		case http.MethodPut:
			var updated Grade
			if !decodeJSON(w, r, &updated) {
				return
			}
			if err := validateGrade(updated); err != nil {
				writeError(w, http.StatusBadRequest, err.Error())
				return
			}
			if updated.Date.IsZero() {
				updated.Date = grade.Date
			}
			updated.ID = id
			grades[i] = updated
			writeJSON(w, http.StatusOK, updated)
		case http.MethodDelete:
			grades = append(grades[:i], grades[i+1:]...)
			writeJSON(w, http.StatusOK, map[string]string{"message": "Grade deleted"})
		default:
			methodNotAllowed(w)
		}
		return
	}
	writeError(w, http.StatusNotFound, "grade not found")
}

func StudentGradesHandler(w http.ResponseWriter, r *http.Request) {
	studentID, ok := pathID(strings.TrimSuffix(r.URL.Path, "/grades"), "/students/")
	if !ok || r.Method != http.MethodGet {
		if !ok {
			writeError(w, http.StatusBadRequest, "invalid student ID")
		} else {
			methodNotAllowed(w)
		}
		return
	}
	if !studentExists(studentID) {
		writeError(w, http.StatusNotFound, "student not found")
		return
	}
	subjectID, err := optionalQueryID(r, "subject_id")
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid subject_id")
		return
	}
	gradesMutex.RLock()
	result := make([]Grade, 0)
	for _, grade := range grades {
		if grade.StudentID == studentID && (subjectID == nil || grade.SubjectID == *subjectID) {
			result = append(result, grade)
		}
	}
	gradesMutex.RUnlock()
	writeJSON(w, http.StatusOK, result)
}

func validateGrade(grade Grade) error {
	if grade.Value < 1 || grade.Value > 5 {
		return &validationError{"grade value must be between 1 and 5"}
	}
	if !studentExists(grade.StudentID) {
		return &validationError{"student not found"}
	}
	if !subjectExists(grade.SubjectID) {
		return &validationError{"subject not found"}
	}
	return nil
}

func optionalQueryID(r *http.Request, name string) (*int, error) {
	value := r.URL.Query().Get(name)
	if value == "" {
		return nil, nil
	}
	id, err := strconv.Atoi(value)
	if err != nil || id < 1 {
		return nil, strconv.ErrSyntax
	}
	return &id, nil
}

type validationError struct {
	message string
}

func (e *validationError) Error() string {
	return e.message
}
