package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/AlexWilliam12/silent-signal/auth"
	"github.com/AlexWilliam12/silent-signal/client"
	"github.com/AlexWilliam12/silent-signal/configs"
	"github.com/AlexWilliam12/silent-signal/database/repositories"
)

// Handler to update a user
func HandleUserUpdate(w http.ResponseWriter, r *http.Request) {
	logger := configs.NewLogger("handlers")

	// Validate if authorization header is present
	authorization := r.Header.Get("Authorization")
	if authorization == "" {
		http.Error(w, "Unauthorized request", http.StatusUnauthorized)
		return
	}

	// Check if authorization is valid
	if !strings.Contains(authorization, "Bearer ") {
		http.Error(w, "Invalid authorization request", http.StatusBadRequest)
		return
	}

	// Get and validate token
	token := strings.Replace(authorization, "Bearer ", "", 1)
	claims, err := auth.ValidateToken(token)
	if err != nil {
		logger.Debug(err)
		http.Error(w, "Unauthorized request", http.StatusUnauthorized)
		return
	}

	// Get client request body
	var user client.UserRequest
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		logger.Debug(err)
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	// Fetch user on database
	fetchedUser, err := repositories.FindUserByName(user.Username)
	if err != nil {
		logger.Debug(err)
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Check if the token user is the same as the request
	if claims.Username != fetchedUser.Username {
		logger.Debug(err)
		http.Error(w, "Unauthorized request", http.StatusForbidden)
		return
	}

	fetchedUser.Username = user.Username
	fetchedUser.Password = user.Password

	// Update user on database
	_, err = repositories.UpdateUser(fetchedUser)
	if err != nil {
		logger.Debug(err)
		http.Error(w, "Unable to update the user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Handler to delete a user
func HandleUserDelete(w http.ResponseWriter, r *http.Request) {
	logger := configs.NewLogger("handlers")

	// Validate if authorization header is present
	authorization := r.Header.Get("Authorization")
	if authorization == "" {
		http.Error(w, "Unauthorized request", http.StatusUnauthorized)
		return
	}

	// Check if authorization is valid
	if !strings.Contains(authorization, "Bearer ") {
		http.Error(w, "Invalid authorization request", http.StatusBadRequest)
		return
	}

	// Get and validate token
	token := strings.Replace(authorization, "Bearer ", "", 1)
	claims, err := auth.ValidateToken(token)
	if err != nil {
		logger.Debug(err)
		http.Error(w, "Unauthorized request", http.StatusUnauthorized)
		return
	}

	// Get all query parameters
	queryParams := r.URL.Query()

	// Get user parameter if present
	userParam := queryParams.Get("user")

	if userParam == "" {
		http.Error(w, "No username specified", http.StatusBadRequest)
		return
	}

	// Fetch user on database
	fetchedUser, err := repositories.FindUserByName(userParam)
	if err != nil {
		logger.Debug(err)
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Check if the token user is the same as the request
	if claims.Username != fetchedUser.Username {
		logger.Debug(err)
		http.Error(w, "Unauthorized request", http.StatusForbidden)
		return
	}

	// Delete user on database
	_, err = repositories.DeleteUserByName(fetchedUser.Username)
	if err != nil {
		logger.Debug(err)
		http.Error(w, "Unable to delete user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
