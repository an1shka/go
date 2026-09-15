package models

// User — учётная запись пользователя в системе
type User struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Password string `json:"-"` // Поле "-" скрывает пароль при сериализации в JSON
}

// Student — сущность ученика
type Student struct {
	ID        string `json:"id"`
	FullName  string `json:"full_name"`  // ФИО
	Grade     string `json:"grade"`      // Класс (например, 10А)
	Age       int    `json:"age"`        // Возраст
	Email     string `json:"email"`      // Email ученика
	UserID    string `json:"user_id"`    // Связь с профилем учетной записи User
}

// RegisterRequest — структура входящего запроса на регистрацию
type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginRequest — структура входящего запроса на вход
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// AuthResponse — структура успешного ответа с токеном
type AuthResponse struct {
	Token string `json:"token"`
}

// ErrorResponse — единый формат ответов с ошибками
type ErrorResponse struct {
	Error string `json:"error"`
}