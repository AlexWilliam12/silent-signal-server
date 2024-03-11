package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/AlexWilliam12/silent-signal/internal/client"
	"github.com/AlexWilliam12/silent-signal/internal/configs"
	"github.com/AlexWilliam12/silent-signal/internal/database/models"
	"github.com/AlexWilliam12/silent-signal/internal/database/repositories"
)

// Handler to process create group request
func HandleCreateGroup(w http.ResponseWriter, r *http.Request) {
	logger := configs.NewLogger("handlers")

	// Check if token is valid
	claims, err := handleAuthorization(w, r)
	if err != nil {
		return
	}

	// Get client request body
	var group client.GroupRequest
	if err := json.NewDecoder(r.Body).Decode(&group); err != nil {
		logger.Debug(err)
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	// Fetch user on database querying by username
	creator, err := repositories.FindUserByName(group.Creator)
	if err != nil {
		logger.Debug(err)
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Check if the token user is the same as the request
	if creator.Username != claims.Username {
		http.Error(w, "Unauthorized request", http.StatusForbidden)
		return
	}

	// Save group on database
	if _, err = repositories.CreateGroup(&models.Group{Name: group.Name, Description: group.Description, Creator: *creator}); err != nil {
		logger.Debug(err)
		if strings.Contains(err.Error(), "duplicate key") {
			http.Error(w, "Group name is already in use", http.StatusBadRequest)
		} else {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// Handler to process a request to fetch all groups
func HandleFetchAllGroups(w http.ResponseWriter, r *http.Request) {
	logger := configs.NewLogger("handlers")

	// Check if token is valid
	_, err := handleAuthorization(w, r)
	if err != nil {
		return
	}

	// Fetch all groups from database
	groups, err := repositories.FindAllGroups()
	if err != nil {
		logger.Debug(err)
		http.Error(w, "Unable to load groups", http.StatusInternalServerError)
		return
	}

	// Change to serializable struct to encode
	var serializableGroups []client.GroupResponse
	for _, group := range groups {
		serializableGroups = append(serializableGroups, client.GroupResponse{
			Name:        group.Name,
			Description: group.Description,
			Creator:     group.Creator.Username,
			PictureURL:  group.Picture,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&serializableGroups)
}

// Handler to process a request to fetch specific group
func HandleFetchGroup(w http.ResponseWriter, r *http.Request) {
	logger := configs.NewLogger("handlers")

	// Check if token is valid
	_, err := handleAuthorization(w, r)
	if err != nil {
		return
	}

	// Get all query parameters
	queryParams := r.URL.Query()

	// Get user parameter if present
	groupParam := queryParams.Get("group")

	if groupParam == "" {
		http.Error(w, "No group name specified", http.StatusBadRequest)
		return
	}

	// Fetch group on database
	group, err := repositories.FindGroupByName(groupParam)
	if err != nil {
		logger.Debug(err)
		http.Error(w, "Group not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&group)
}

// Handler to process update group requests
func HandleUpdateGroup(w http.ResponseWriter, r *http.Request) {
	logger := configs.NewLogger("handlers")

	// Check if token is valid
	claims, err := handleAuthorization(w, r)
	if err != nil {
		return
	}

	// Get all query parameters
	queryParams := r.URL.Query()

	// Get user parameter if present
	groupParam := queryParams.Get("group")

	if groupParam == "" {
		http.Error(w, "No group name specified", http.StatusBadRequest)
		return
	}

	// Get client request body
	var group client.GroupRequest
	if err := json.NewDecoder(r.Body).Decode(&group); err != nil {
		logger.Debug(err)
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	// Fetch group on database
	fetchedGroup, err := repositories.FindGroupByName(groupParam)
	if err != nil {
		logger.Debug(err)
		http.Error(w, "Group not found", http.StatusNotFound)
		return
	}

	// Check if user token is the same as the request
	if fetchedGroup.Creator.Username != claims.Username {
		http.Error(w, "Unauthorized request", http.StatusForbidden)
		return
	}

	fetchedGroup.Name = group.Name
	fetchedGroup.Description = group.Description

	// Update group on database
	_, err = repositories.UpdateGroup(fetchedGroup)
	if err != nil {
		logger.Debug(err)
		http.Error(w, "Unable to update group", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func HandleDeleteGroup(w http.ResponseWriter, r *http.Request) {
	logger := configs.NewLogger("handlers")

	// Check if token is valid
	claims, err := handleAuthorization(w, r)
	if err != nil {
		return
	}

	// Get all query parameters
	queryParams := r.URL.Query()

	// Get user parameter if present
	groupParam := queryParams.Get("group")

	if groupParam == "" {
		http.Error(w, "No group name specified", http.StatusBadRequest)
		return
	}

	// Fetch group on database querying by name
	group, err := repositories.FindGroupByName(groupParam)
	if err != nil {
		logger.Debug(err)
		http.Error(w, "Group not found", http.StatusNotFound)
		return
	}

	if group.Creator.Username != claims.Username {
		http.Error(w, "Unauthorized request", http.StatusForbidden)
		return
	}

	// Delete group on database
	_, err = repositories.DeleteGroupByName(group.Name)
	if err != nil {
		logger.Debug(err)
		http.Error(w, "Unable to delete group", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
