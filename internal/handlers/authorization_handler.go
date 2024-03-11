package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/AlexWilliam12/silent-signal/internal/auth"
	"github.com/AlexWilliam12/silent-signal/internal/configs"
)

// Handler to process token claims
func handleAuthorization(w http.ResponseWriter, r *http.Request) (*auth.CustomClaims, error) {
	logger := configs.NewLogger("handlers")

	// Validate if authorization header is present
	authorization := r.Header.Get("Authorization")
	if authorization == "" {
		http.Error(w, "Unauthorized request", http.StatusUnauthorized)
		return nil, fmt.Errorf("unauthorized request")
	}

	// Check if authorization is valid
	if !strings.Contains(authorization, "Bearer ") {
		http.Error(w, "Invalid authorization request", http.StatusBadRequest)
		return nil, fmt.Errorf("invalid authorization request")
	}

	// Get and validate token
	token := strings.Replace(authorization, "Bearer ", "", 1)
	claims, err := auth.ValidateToken(token)
	if err != nil {
		logger.Debug(err)
		http.Error(w, "Unauthorized request", http.StatusUnauthorized)
		return nil, fmt.Errorf("unauthorized request")
	}

	return claims, nil
}
