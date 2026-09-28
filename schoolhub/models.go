package main

import "time"

type User struct {
	Username string `json:"username"`
	Password string `json:"password,omitempty"`
}

type Student struct {
	ID        int    `json:"id"`
	FullName  string `json:"full_name"`
	Email     string `json:"email"`
	ClassName string `json:"class_name"`
}

type Teacher struct {
	ID       int    `json:"id"`
	FullName string `json:"full_name"`
	Email    string `json:"email"`
}

type Subject struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	TeacherID int    `json:"teacher_id"`
}

type Grade struct {
	ID        int       `json:"id"`
	StudentID int       `json:"student_id"`
	SubjectID int       `json:"subject_id"`
	Value     int       `json:"value"`
	Date      time.Time `json:"date"`
	Comment   string    `json:"comment"`
}

type Homework struct {
	ID          int       `json:"id"`
	SubjectID   int       `json:"subject_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	IssuedAt    time.Time `json:"issued_at"`
	Deadline    time.Time `json:"deadline"`
}
