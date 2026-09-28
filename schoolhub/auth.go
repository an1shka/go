package main

import (
	"net/http"
	"strings"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		methodNotAllowed(w)
		return
	}
	var credentials User
	if !decodeJSON(w, r, &credentials) {
		return
	}
	if strings.TrimSpace(credentials.Username) != "admin" || credentials.Password != "schoolhub123" {
		writeError(w, http.StatusUnauthorized, "invalid username or password")
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{
		"token": demoToken,
		"type":  "Bearer",
	})
}
