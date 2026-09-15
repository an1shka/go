package main

import (
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	serverAddr       = ":8080"
	tokenLifetime    = 24 * time.Hour
	passwordSaltSize = 16
	passwordRounds   = 120000
)

// User хранит учетную запись. Пароль сохраняется только в виде соли и хеша.
type User struct {
	ID           int    `json:"id"`
	Email        string `json:"email"`
	PasswordHash string `json:"-"`
	PasswordSalt string `json:"-"`
}

// Student описывает профиль ученика, который доступен через CRUD-маршруты.
type Student struct {
	ID          int    `json:"id"`
	FullName    string `json:"full_name"`
	Class       string `json:"class"`
	Age         int    `json:"age"`
	DateOfBirth string `json:"date_of_birth,omitempty"`
	Email       string `json:"email"`
	Gender      string `json:"gender"`
	UserID      int    `json:"user_id"`
}

// studentInput отделяет данные запроса от системных полей ID и UserID.
type studentInput struct {
	FullName    string `json:"full_name"`
	Class       string `json:"class"`
	Age         int    `json:"age"`
	DateOfBirth string `json:"date_of_birth"`
	Email       string `json:"email"`
	Gender      string `json:"gender"`
}

// registerInput и loginInput содержат тела запросов авторизации.
type registerInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// apiError задает единый формат ошибок для всех маршрутов.
type apiError struct {
	Error string `json:"error"`
}

// tokenClaims находятся внутри подписанного токена.
type tokenClaims struct {
	UserID    int   `json:"user_id"`
	ExpiresAt int64 `json:"expires_at"`
}

// store - простое in-memory хранилище для учебного проекта.
type store struct {
	mu            sync.RWMutex
	users         map[int]User
	usersByEmail  map[string]int
	students      map[int]Student
	nextUserID    int
	nextStudentID int
}

func newStore() *store {
	return &store{
		users:         make(map[int]User),
		usersByEmail:  make(map[string]int),
		students:      make(map[int]Student),
		nextUserID:    1,
		nextStudentID: 1,
	}
}

func main() {
	app := newStore()
	server := &http.Server{Addr: serverAddr, Handler: app.routes()}
	log.Printf("SchoolHub API запущен на http://localhost%s", serverAddr)
	log.Fatal(server.ListenAndServe())
}

// routes регистрирует публичные и защищенные маршруты API.
func (s *store) routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/auth/register", s.handleRegister)
	mux.HandleFunc("/auth/login", s.handleLogin)
	mux.HandleFunc("/me", s.requireAuth(s.handleMe))
	mux.HandleFunc("/students", s.requireAuth(s.handleStudents))
	mux.HandleFunc("/students/", s.requireAuth(s.handleStudentByID))
	mux.HandleFunc("/", handleHome)
	return jsonMiddleware(mux)
}

func (s *store) handleRegister(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "метод не поддерживается")
		return
	}
	var input registerInput
	if !decodeJSON(w, r, &input) || strings.TrimSpace(input.Email) == "" || input.Password == "" {
		writeError(w, http.StatusBadRequest, "email и password обязательны")
		return
	}

	email := normalizeEmail(input.Email)
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.usersByEmail[email]; exists {
		writeError(w, http.StatusConflict, "email уже зарегистрирован")
		return
	}
	salt, err := randomBytes(passwordSaltSize)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "не удалось создать пользователя")
		return
	}
	user := User{
		ID:           s.nextUserID,
		Email:        email,
		PasswordSalt: hex.EncodeToString(salt),
		PasswordHash: hashPassword(input.Password, salt),
	}
	s.nextUserID++
	s.users[user.ID] = user
	s.usersByEmail[email] = user.ID
	writeJSON(w, http.StatusCreated, map[string]any{"id": user.ID, "email": user.Email})
}

func (s *store) handleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "метод не поддерживается")
		return
	}
	var input loginInput
	if !decodeJSON(w, r, &input) || strings.TrimSpace(input.Email) == "" || input.Password == "" {
		writeError(w, http.StatusBadRequest, "email и password обязательны")
		return
	}

	s.mu.RLock()
	userID, exists := s.usersByEmail[normalizeEmail(input.Email)]
	user := s.users[userID]
	s.mu.RUnlock()
	if !exists || !passwordMatches(input.Password, user) {
		writeError(w, http.StatusUnauthorized, "неверный email или password")
		return
	}
	token, err := createToken(user.ID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "не удалось создать токен")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"token": token, "expires_in": int(tokenLifetime.Seconds())})
}

func (s *store) handleMe(w http.ResponseWriter, r *http.Request) {
	userID := userIDFromContext(r)
	s.mu.RLock()
	user := s.users[userID]
	s.mu.RUnlock()
	writeJSON(w, http.StatusOK, map[string]any{"id": user.ID, "email": user.Email})
}

func (s *store) handleStudents(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		s.mu.RLock()
		students := make([]Student, 0, len(s.students))
		for _, student := range s.students {
			students = append(students, student)
		}
		s.mu.RUnlock()
		writeJSON(w, http.StatusOK, students)
	case http.MethodPost:
		var input studentInput
		if !decodeJSON(w, r, &input) || !validStudentInput(input) {
			writeError(w, http.StatusBadRequest, "full_name, class, age и email обязательны")
			return
		}
		s.mu.Lock()
		student := Student{ID: s.nextStudentID, UserID: userIDFromContext(r), FullName: strings.TrimSpace(input.FullName), Class: strings.TrimSpace(input.Class), Age: input.Age, DateOfBirth: strings.TrimSpace(input.DateOfBirth), Email: normalizeEmail(input.Email), Gender: strings.TrimSpace(input.Gender)}
		s.nextStudentID++
		s.students[student.ID] = student
		s.mu.Unlock()
		writeJSON(w, http.StatusCreated, student)
	default:
		writeError(w, http.StatusMethodNotAllowed, "метод не поддерживается")
	}
}

func (s *store) handleStudentByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/students/"))
	if err != nil || id < 1 {
		writeError(w, http.StatusBadRequest, "id должен быть целым числом")
		return
	}

	switch r.Method {
	case http.MethodGet:
		s.mu.RLock()
		student, exists := s.students[id]
		s.mu.RUnlock()
		if !exists {
			writeError(w, http.StatusNotFound, "ученик не найден")
			return
		}
		writeJSON(w, http.StatusOK, student)
	case http.MethodPut, http.MethodPatch:
		var input studentInput
		if !decodeJSON(w, r, &input) || !validStudentInput(input) {
			writeError(w, http.StatusBadRequest, "full_name, class, age и email обязательны")
			return
		}
		s.mu.Lock()
		student, exists := s.students[id]
		if exists {
			student.FullName = strings.TrimSpace(input.FullName)
			student.Class = strings.TrimSpace(input.Class)
			student.Age = input.Age
			student.DateOfBirth = strings.TrimSpace(input.DateOfBirth)
			student.Email = normalizeEmail(input.Email)
			student.Gender = strings.TrimSpace(input.Gender)
			s.students[id] = student
		}
		s.mu.Unlock()
		if !exists {
			writeError(w, http.StatusNotFound, "ученик не найден")
			return
		}
		writeJSON(w, http.StatusOK, student)
	case http.MethodDelete:
		s.mu.Lock()
		_, exists := s.students[id]
		delete(s.students, id)
		s.mu.Unlock()
		if !exists {
			writeError(w, http.StatusNotFound, "ученик не найден")
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"message": "ученик удален"})
	default:
		writeError(w, http.StatusMethodNotAllowed, "метод не поддерживается")
	}
}

// requireAuth проверяет Bearer-токен до передачи запроса защищенному handler.
func (s *store) requireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := parseToken(r.Header.Get("Authorization"))
		if err != nil {
			writeError(w, http.StatusUnauthorized, "нужен корректный Bearer-токен")
			return
		}
		r = r.WithContext(withUserID(r, userID))
		next(w, r)
	}
}

type contextKey string

const userIDKey contextKey = "user_id"

func withUserID(r *http.Request, userID int) context.Context {
	return context.WithValue(r.Context(), userIDKey, userID)
}

func userIDFromContext(r *http.Request) int {
	userID, _ := r.Context().Value(userIDKey).(int)
	return userID
}

func createToken(userID int) (string, error) {
	claims := tokenClaims{UserID: userID, ExpiresAt: time.Now().Add(tokenLifetime).Unix()}
	payload, err := json.Marshal(claims)
	if err != nil {
		return "", err
	}
	encodedPayload := base64.RawURLEncoding.EncodeToString(payload)
	mac := hmac.New(sha256.New, tokenSecret())
	_, _ = mac.Write([]byte(encodedPayload))
	signature := base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
	return encodedPayload + "." + signature, nil
}

func parseToken(header string) (int, error) {
	parts := strings.Fields(header)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return 0, errors.New("invalid authorization header")
	}
	parts = strings.Split(parts[1], ".")
	if len(parts) != 2 {
		return 0, errors.New("invalid token")
	}
	mac := hmac.New(sha256.New, tokenSecret())
	_, _ = mac.Write([]byte(parts[0]))
	signature, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil || !hmac.Equal(signature, mac.Sum(nil)) {
		return 0, errors.New("invalid signature")
	}
	payload, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return 0, errors.New("invalid payload")
	}
	var claims tokenClaims
	if json.Unmarshal(payload, &claims) != nil || claims.UserID < 1 || claims.ExpiresAt <= time.Now().Unix() {
		return 0, errors.New("expired token")
	}
	return claims.UserID, nil
}

func tokenSecret() []byte {
	// Для учебного проекта секрет задается в коде; в реальном проекте используйте env/secret manager.
	return []byte("schoolhub-development-secret-change-me")
}

func hashPassword(password string, salt []byte) string {
	digest := sha256.Sum256(append(salt, []byte(password)...))
	for i := 1; i < passwordRounds; i++ {
		digest = sha256.Sum256(append(digest[:], salt...))
	}
	return hex.EncodeToString(digest[:])
}

func passwordMatches(password string, user User) bool {
	salt, err := hex.DecodeString(user.PasswordSalt)
	if err != nil {
		return false
	}
	expected := hashPassword(password, salt)
	return hmac.Equal([]byte(expected), []byte(user.PasswordHash))
}

func randomBytes(size int) ([]byte, error) {
	result := make([]byte, size)
	_, err := rand.Read(result)
	return result, err
}

func validStudentInput(input studentInput) bool {
	return strings.TrimSpace(input.FullName) != "" && strings.TrimSpace(input.Class) != "" && input.Age > 0 && input.Age < 120 && strings.TrimSpace(input.Email) != ""
}

func normalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

func decodeJSON(w http.ResponseWriter, r *http.Request, target any) bool {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(target); err != nil {
		writeError(w, http.StatusBadRequest, "некорректный JSON")
		return false
	}
	return true
}

func jsonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		next.ServeHTTP(w, r)
	})
}

func handleNotFound(w http.ResponseWriter, r *http.Request) {
	writeError(w, http.StatusNotFound, "маршрут не найден")
}

// handleHome показывает доступные маршруты при открытии localhost в браузере.
func handleHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		handleNotFound(w, r)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"message": "SchoolHub API работает",
		"routes": []string{
			"POST /auth/register",
			"POST /auth/login",
			"GET /me [Bearer token]",
			"GET|POST /students [Bearer token]",
			"GET|PUT|PATCH|DELETE /students/{id} [Bearer token]",
		},
	})
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("ошибка ответа JSON: %v", err)
	}
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, apiError{Error: message})
}
