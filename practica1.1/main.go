package main

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/mail"
	"strconv"
	"strings"
	"sync"
	"time"
)

type User struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
}

type Student struct {
	ID       int    `json:"id"`
	FullName string `json:"fullName"`
	Class    string `json:"class"`
	Age      int    `json:"age"`
	Email    string `json:"email"`
	UserID   int    `json:"userId"`
}

type registerRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type studentRequest struct {
	FullName string `json:"fullName"`
	Class    string `json:"class"`
	Age      int    `json:"age"`
	Email    string `json:"email"`
}

type app struct {
	mu            sync.RWMutex
	users         map[int]User
	passwords     map[int]string
	tokens        map[string]int
	students      map[int]Student
	nextUserID    int
	nextStudentID int
}

func newApp() *app {
	return &app{
		users:         make(map[int]User),
		passwords:     make(map[int]string),
		tokens:        make(map[string]int),
		students:      make(map[int]Student),
		nextUserID:    1,
		nextStudentID: 1,
	}
}

func main() {
	application := newApp()
	server := &http.Server{
		Addr:              ":8080",
		Handler:           application,
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Println("SchoolHub запущен: http://localhost:8080")
	log.Fatal(server.ListenAndServe())
}

func (a *app) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == http.MethodPost && r.URL.Path == "/auth/register":
		a.register(w, r)
	case r.Method == http.MethodPost && r.URL.Path == "/auth/login":
		a.login(w, r)
	case r.URL.Path == "/me":
		a.withAuth(w, r, a.me)
	case strings.HasPrefix(r.URL.Path, "/students"):
		a.withAuth(w, r, a.studentsHandler)
	default:
		writeError(w, http.StatusNotFound, "маршрут не найден")
	}
}

func (a *app) register(w http.ResponseWriter, r *http.Request) {
	var request registerRequest
	if !readJSON(w, r, &request) {
		return
	}

	request.Email = normalizeEmail(request.Email)
	if err := validateEmail(request.Email); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if len(request.Password) < 6 {
		writeError(w, http.StatusBadRequest, "пароль должен содержать минимум 6 символов")
		return
	}

	a.mu.Lock()
	defer a.mu.Unlock()

	for _, user := range a.users {
		if user.Email == request.Email {
			writeError(w, http.StatusConflict, "этот email уже зарегистрирован")
			return
		}
	}

	user := User{ID: a.nextUserID, Email: request.Email}
	a.nextUserID++
	a.users[user.ID] = user
	a.passwords[user.ID] = hashPassword(request.Password)
	writeJSON(w, http.StatusCreated, user)
}

func (a *app) login(w http.ResponseWriter, r *http.Request) {
	var request loginRequest
	if !readJSON(w, r, &request) {
		return
	}

	email := normalizeEmail(request.Email)
	a.mu.Lock()
	defer a.mu.Unlock()

	var user User
	for _, savedUser := range a.users {
		if savedUser.Email == email {
			user = savedUser
			break
		}
	}

	if user.ID == 0 || a.passwords[user.ID] != hashPassword(request.Password) {
		writeError(w, http.StatusUnauthorized, "неверный email или пароль")
		return
	}

	token, err := newToken()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "не удалось создать токен")
		return
	}
	a.tokens[token] = user.ID
	writeJSON(w, http.StatusOK, map[string]interface{}{"token": token, "user": user})
}

func (a *app) withAuth(w http.ResponseWriter, r *http.Request, handler func(http.ResponseWriter, *http.Request, int)) {
	authorization := strings.TrimSpace(r.Header.Get("Authorization"))
	if !strings.HasPrefix(authorization, "Bearer ") {
		writeError(w, http.StatusUnauthorized, "нужен токен Bearer")
		return
	}
	token := strings.TrimSpace(strings.TrimPrefix(authorization, "Bearer "))

	a.mu.RLock()
	userID, ok := a.tokens[token]
	a.mu.RUnlock()
	if !ok {
		writeError(w, http.StatusUnauthorized, "неверный токен")
		return
	}
	handler(w, r, userID)
}

func (a *app) me(w http.ResponseWriter, _ *http.Request, userID int) {
	a.mu.RLock()
	user, ok := a.users[userID]
	a.mu.RUnlock()
	if !ok {
		writeError(w, http.StatusUnauthorized, "пользователь не найден")
		return
	}
	writeJSON(w, http.StatusOK, user)
}

func (a *app) studentsHandler(w http.ResponseWriter, r *http.Request, userID int) {
	path := strings.TrimPrefix(r.URL.Path, "/students")
	if path == "" || path == "/" {
		a.studentsCollection(w, r, userID)
		return
	}

	idText := strings.TrimPrefix(path, "/")
	id, err := strconv.Atoi(idText)
	if err != nil || id < 1 || strings.Contains(idText, "/") {
		writeError(w, http.StatusNotFound, "ученик не найден")
		return
	}
	a.studentByID(w, r, userID, id)
}

func (a *app) studentsCollection(w http.ResponseWriter, r *http.Request, userID int) {
	switch r.Method {
	case http.MethodGet:
		a.mu.RLock()
		students := make([]Student, 0, len(a.students))
		for _, student := range a.students {
			students = append(students, student)
		}
		a.mu.RUnlock()
		writeJSON(w, http.StatusOK, students)
	case http.MethodPost:
		var request studentRequest
		if !readJSON(w, r, &request) {
			return
		}
		if err := validateStudent(request); err != nil {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}

		a.mu.Lock()
		student := Student{
			ID:       a.nextStudentID,
			FullName: strings.TrimSpace(request.FullName),
			Class:    strings.TrimSpace(request.Class),
			Age:      request.Age,
			Email:    normalizeEmail(request.Email),
			UserID:   userID,
		}
		a.nextStudentID++
		a.students[student.ID] = student
		a.mu.Unlock()
		writeJSON(w, http.StatusCreated, student)
	default:
		writeError(w, http.StatusMethodNotAllowed, "метод не поддерживается")
	}
}

func (a *app) studentByID(w http.ResponseWriter, r *http.Request, userID, id int) {
	switch r.Method {
	case http.MethodGet:
		a.mu.RLock()
		student, ok := a.students[id]
		a.mu.RUnlock()
		if !ok {
			writeError(w, http.StatusNotFound, "ученик не найден")
			return
		}
		writeJSON(w, http.StatusOK, student)
	case http.MethodPut, http.MethodPatch:
		var request studentRequest
		if !readJSON(w, r, &request) {
			return
		}
		if err := validateStudent(request); err != nil {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}

		a.mu.Lock()
		if _, ok := a.students[id]; !ok {
			a.mu.Unlock()
			writeError(w, http.StatusNotFound, "ученик не найден")
			return
		}
		student := Student{
			ID:       id,
			FullName: strings.TrimSpace(request.FullName),
			Class:    strings.TrimSpace(request.Class),
			Age:      request.Age,
			Email:    normalizeEmail(request.Email),
			UserID:   userID,
		}
		a.students[id] = student
		a.mu.Unlock()
		writeJSON(w, http.StatusOK, student)
	case http.MethodDelete:
		a.mu.Lock()
		if _, ok := a.students[id]; !ok {
			a.mu.Unlock()
			writeError(w, http.StatusNotFound, "ученик не найден")
			return
		}
		delete(a.students, id)
		a.mu.Unlock()
		writeJSON(w, http.StatusOK, map[string]string{"message": "ученик удалён"})
	default:
		writeError(w, http.StatusMethodNotAllowed, "метод не поддерживается")
	}
}

func readJSON(w http.ResponseWriter, r *http.Request, destination interface{}) bool {
	decoder := json.NewDecoder(http.MaxBytesReader(w, r.Body, 1<<20))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(destination); err != nil {
		writeError(w, http.StatusBadRequest, "некорректный JSON")
		return false
	}
	if decoder.Decode(&struct{}{}) == nil {
		writeError(w, http.StatusBadRequest, "в запросе должен быть один JSON-объект")
		return false
	}
	return true
}

func validateStudent(request studentRequest) error {
	if strings.TrimSpace(request.FullName) == "" {
		return errors.New("поле fullName обязательно")
	}
	if strings.TrimSpace(request.Class) == "" {
		return errors.New("поле class обязательно")
	}
	if request.Age < 6 || request.Age > 100 {
		return errors.New("возраст должен быть от 6 до 100 лет")
	}
	return validateEmail(normalizeEmail(request.Email))
}

func validateEmail(email string) error {
	if email == "" {
		return errors.New("поле email обязательно")
	}
	address, err := mail.ParseAddress(email)
	if err != nil || address.Address != email {
		return errors.New("укажите корректный email")
	}
	return nil
}

func normalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

func hashPassword(password string) string {
	hash := sha256.Sum256([]byte(password))
	return hex.EncodeToString(hash[:])
}

func newToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("generate token: %w", err)
	}
	return hex.EncodeToString(bytes), nil
}

func writeJSON(w http.ResponseWriter, status int, value interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(value); err != nil {
		log.Printf("ошибка отправки JSON: %v", err)
	}
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}
