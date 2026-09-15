package auth

import (
	"crypto/rand"
	"encoding/hex"
	"strings"
)

// GenerateToken создает случайный защищенный токен для сессии
func GenerateToken() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// ExtractToken извлекает токен из заголовка Authorization (формат: "Bearer <token>")
func ExtractToken(authHeader string) string {
	parts := strings.Split(authHeader, " ")
	if len(parts) == 2 && strings.ToLower(parts[0]) == "bearer" {
		return parts[1]
	}
	return ""
}