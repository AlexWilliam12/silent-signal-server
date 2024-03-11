package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/AlexWilliam12/silent-signal/internal/auth"
	"github.com/AlexWilliam12/silent-signal/internal/client"
	"github.com/AlexWilliam12/silent-signal/internal/configs"
	"github.com/AlexWilliam12/silent-signal/internal/database/repositories"
)

// Handler to process the client authentication request
func HandleLogin(w http.ResponseWriter, r *http.Request) {
	logger := configs.NewLogger("handlers")

	// Get client request body
	var user client.UserRequest
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	// Fetch user on database querying by user credentials
	if _, err := repositories.FindUserByCredentials(user); err != nil {
		if strings.Contains(err.Error(), "record not found") {
			http.Error(w, "Invalid credentials", http.StatusNotFound)
		} else {
			logger.Err(err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	// Generate JWT token for user
	token, err := auth.GenerateToken(user.Username)
	if err != nil {
		logger.Err(err)
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(client.JWTToken{Token: token})
}

// Handler to process a request to register user
func HandleRegister(w http.ResponseWriter, r *http.Request) {
	logger := configs.NewLogger("handlers")

	// Get client request body
	var user client.UserRequest
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	// Save a user on database
	if _, err := repositories.CreateUser(user); err != nil {
		logger.Err(err)
		if strings.Contains(err.Error(), "duplicate key") {
			http.Error(w, "Username is already in use", http.StatusBadRequest)
		} else {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusCreated)
}
