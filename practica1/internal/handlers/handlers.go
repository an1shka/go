package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"main/internal/auth"
	"main/internal/models"
	"main/internal/store"
)

type Handler struct {
	Store *store.Store
}

func NewHandler(st *store.Store) *Handler {
	return &Handler{Store: st}
}

// sendJSON — утилита для отправки корректных JSON-ответов
func sendJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}

// sendError — утилита для единой обработки ошибок в JSON
func sendError(w http.ResponseWriter, status int, message string) {
	sendJSON(w, status, models.ErrorResponse{Error: message})
}

// AuthMiddleware проверяет наличие корректного токена авторизации
func (h *Handler) AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := auth.ExtractToken(r.Header.Get("Authorization"))
		if token == "" {
			sendError(w, http.StatusUnauthorized, "Отсутствует или неверный токен авторизации")
			return
		}

		_, found := h.Store.GetUserByToken(token)
		if !found {
			sendError(w, http.StatusUnauthorized, "Недействительный токен")
			return
		}

		next(w, r)
	}
}

// Register — POST /auth/register
func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		sendError(w, http.StatusMethodNotAllowed, "Метод не поддерживается")
		return
	}

	var req models.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, http.StatusBadRequest, "Невалидный JSON запроса")
		return
	}

	if strings.TrimSpace(req.Email) == "" || strings.TrimSpace(req.Password) == "" {
		sendError(w, http.StatusBadRequest, "Email и пароль обязательны")
		return
	}

	user, err := h.Store.CreateUser(req.Email, req.Password)
	if err == store.ErrUserExists {
		sendError(w, http.StatusConflict, err.Error())
		return
	} else if err != nil {
		sendError(w, http.StatusInternalServerError, "Ошибка сервера при создании пользователя")
		return
	}

	sendJSON(w, http.StatusCreated, map[string]interface{}{
		"id":    user.ID,
		"email": user.Email,
	})
}

// Login — POST /auth/login
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		sendError(w, http.StatusMethodNotAllowed, "Метод не поддерживается")
		return
	}

	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, http.StatusBadRequest, "Невалидный JSON запроса")
		return
	}

	user, exists := h.Store.GetUserByEmail(req.Email)
	if !exists || user.Password != req.Password {
		sendError(w, http.StatusUnauthorized, "Неверный email или пароль")
		return
	}

	token, err := auth.GenerateToken()
	if err != nil {
		sendError(w, http.StatusInternalServerError, "Не удалось сгенерировать токен")
		return
	}

	h.Store.SaveToken(token, user.ID)
	sendJSON(w, http.StatusOK, models.AuthResponse{Token: token})
}

// GetMe — GET /me
func (h *Handler) GetMe(w http.ResponseWriter, r *http.Request) {
	token := auth.ExtractToken(r.Header.Get("Authorization"))
	user, _ := h.Store.GetUserByToken(token)

	sendJSON(w, http.StatusOK, map[string]string{
		"id":    user.ID,
		"email": user.Email,
	})
}

// StudentsRouter обрабатывает маршуты /students и /students/{id} на чистом net/http
func (h *Handler) StudentsRouter(w http.ResponseWriter, r *http.Request) {
	// Извлекаем ID из URL-пути (например /students/1 -> "1")
	path := strings.TrimPrefix(r.URL.Path, "/students")
	path = strings.Trim(path, "/")

	if path == "" {
		// Коллекция /students
		switch r.Method {
		case http.MethodGet:
			students := h.Store.GetAllStudents()
			sendJSON(w, http.StatusOK, students)
		case http.MethodPost:
			var st models.Student
			if err := json.NewDecoder(r.Body).Decode(&st); err != nil {
				sendError(w, http.StatusBadRequest, "Невалидный JSON ученика")
				return
			}
			if strings.TrimSpace(st.FullName) == "" {
				sendError(w, http.StatusBadRequest, "ФИО ученика обязательно")
				return
			}
			created := h.Store.CreateStudent(st)
			sendJSON(w, http.StatusCreated, created)
		default:
			sendError(w, http.StatusMethodNotAllowed, "Метод не поддерживается")
		}
		return
	}

	// Элемент /students/{id}
	id := path
	switch r.Method {
	case http.MethodGet:
		st, exists := h.Store.GetStudentByID(id)
		if !exists {
			sendError(w, http.StatusNotFound, "Ученик не найден")
			return
		}
		sendJSON(w, http.StatusOK, st)

	case http.MethodPut, http.MethodPatch:
		var st models.Student
		if err := json.NewDecoder(r.Body).Decode(&st); err != nil {
			sendError(w, http.StatusBadRequest, "Невалидный JSON")
			return
		}
		updated, err := h.Store.UpdateStudent(id, st)
		if err == store.ErrStudentNotFound {
			sendError(w, http.StatusNotFound, "Ученик не найден")
			return
		}
		sendJSON(w, http.StatusOK, updated)

	case http.MethodDelete:
		err := h.Store.DeleteStudent(id)
		if err == store.ErrStudentNotFound {
			sendError(w, http.StatusNotFound, "Ученик не найден")
			return
		}
		sendJSON(w, http.StatusOK, map[string]string{"message": "Ученик успешно удален"})

	default:
		sendError(w, http.StatusMethodNotAllowed, "Метод не поддерживается")
	}
}