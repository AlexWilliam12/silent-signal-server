package handlers

import (
	"log"
	"net/http"
	"strings"

	"github.com/AlexWilliam12/silent-signal/internal/auth"
)

func handleAuthorization(w http.ResponseWriter, r *http.Request) *auth.CustomClaims {

	authorization := r.Header.Get("Authorization")
	if authorization == "" {
		http.Error(w, "Unauthorized request", http.StatusUnauthorized)
		return nil
	}

	if !strings.Contains(authorization, "Bearer ") {
		http.Error(w, "Invalid authorization request", http.StatusBadRequest)
		return nil
	}

	token := strings.Replace(authorization, "Bearer ", "", 1)
	claims, err := auth.ValidateToken(token)
	if err != nil {
		log.Println(err)
		http.Error(w, "Unauthorized request", http.StatusUnauthorized)
		return nil
	}

	return claims
}
