package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/AlexWilliam12/silent-signal/internal/client"
	"github.com/AlexWilliam12/silent-signal/internal/configs"
	"github.com/AlexWilliam12/silent-signal/internal/database/repositories"
)

// Handler to update a user
func HandleUserUpdate(w http.ResponseWriter, r *http.Request) {
	logger := configs.NewLogger("handlers")

	// Check if token is valid
	claims, err := handleAuthorization(w, r)
	if err != nil {
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

	// Check if user token is the same as the request
	if claims.Username != userParam {
		http.Error(w, "Unauthorized request", http.StatusForbidden)
		return
	}

	// Fetch user on database
	fetchedUser, err := repositories.FindUserByName(userParam)
	if err != nil {
		logger.Debug(err)
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Get client request body
	var user client.UserRequest
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		logger.Debug(err)
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
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

	// Check if token is valid
	claims, err := handleAuthorization(w, r)
	if err != nil {
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

	// Check if user token is the same as the request
	if claims.Username != userParam {
		http.Error(w, "Unauthorized request", http.StatusForbidden)
		return
	}

	// Fetch user on database
	fetchedUser, err := repositories.FindUserByName(userParam)
	if err != nil {
		logger.Debug(err)
		http.Error(w, "User not found", http.StatusNotFound)
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
