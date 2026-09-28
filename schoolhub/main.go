package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
)

const (
	serverPort = ":8080"
	demoToken  = "schoolhub-demo-token"
)

func main() {
	seedData()
	mux := http.NewServeMux()
	mux.HandleFunc("/auth/login", LoginHandler)
	mux.HandleFunc("/students", protected(StudentsHandler))
	mux.HandleFunc("/students/", protected(studentRoute))
	mux.HandleFunc("/teachers", protected(TeachersHandler))
	mux.HandleFunc("/teachers/", protected(TeacherByIDHandler))
	mux.HandleFunc("/subjects", protected(SubjectsHandler))
	mux.HandleFunc("/subjects/", protected(subjectRoute))
	mux.HandleFunc("/grades", protected(GradesHandler))
	mux.HandleFunc("/grades/", protected(GradeByIDHandler))
	mux.HandleFunc("/homework", protected(HomeworkHandler))
	mux.HandleFunc("/homework/", protected(HomeworkByIDHandler))

	log.Printf("SchoolHub API is running on http://localhost%s", serverPort)
	log.Fatal(http.ListenAndServe(serverPort, mux))
}

func studentRoute(w http.ResponseWriter, r *http.Request) {
	if strings.HasSuffix(r.URL.Path, "/grades") {
		StudentGradesHandler(w, r)
		return
	}
	StudentByIDHandler(w, r)
}

func subjectRoute(w http.ResponseWriter, r *http.Request) {
	if strings.HasSuffix(r.URL.Path, "/homework") {
		id, ok := pathID(strings.TrimSuffix(r.URL.Path, "/homework"), "/subjects/")
		if !ok || r.Method != http.MethodGet {
			if !ok {
				writeError(w, http.StatusBadRequest, "invalid subject ID")
			} else {
				methodNotAllowed(w)
			}
			return
		}
		if !subjectExists(id) {
			writeError(w, http.StatusNotFound, "subject not found")
			return
		}
		r2 := r.Clone(r.Context())
		query := r2.URL.Query()
		query.Set("subject_id", stringID(id))
		r2.URL.RawQuery = query.Encode()
		HomeworkHandler(w, r2)
		return
	}
	SubjectByIDHandler(w, r)
}

func stringID(id int) string {
	return strconv.Itoa(id)
}

func protected(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			next(w, r)
			return
		}
		auth := strings.TrimSpace(r.Header.Get("Authorization"))
		if auth != "Bearer "+demoToken {
			writeError(w, http.StatusUnauthorized, "missing or invalid bearer token")
			return
		}
		next(w, r)
	}
}

func decodeJSON(w http.ResponseWriter, r *http.Request, destination any) bool {
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(destination); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return false
	}
	return true
}
