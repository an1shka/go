package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/teachers", TeachersHandler)
	http.HandleFunc("/teachers/", TeacherByIDHandler)

	http.HandleFunc("/subjects", SubjectsHandler)

	http.HandleFunc("/grades", GradesHandler)
	http.HandleFunc("/students/", StudentGradesHandler)

	http.HandleFunc("/homework", HomeworkHandler)

	fmt.Println("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}