package handlers

import (
	"encoding/json"
	"mime/multipart"
	"net/http"
	"strings"

	"github.com/AlexWilliam12/silent-signal/auth"
	"github.com/AlexWilliam12/silent-signal/client"
	"github.com/AlexWilliam12/silent-signal/configs"
	"github.com/AlexWilliam12/silent-signal/database/repositories"
	"github.com/AlexWilliam12/silent-signal/services"
)

// Handler to fetch image queried by user or group name
func HandleFetchPicture(w http.ResponseWriter, r *http.Request) {
	logger := configs.NewLogger("handlers")

	// Validate if authorization header is present
	authorization := r.Header.Get("Authorization")
	if authorization != "" {
		http.Error(w, "Unauthorized request", http.StatusUnauthorized)
		return
	}

	// Get and validate token
	token := authorization[strings.Index(authorization, "Bearer "):]
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

	// Get group parameter if present
	groupParam := queryParams.Get("group")

	// Check if at least one argument is present
	if userParam == "" && groupParam == "" {
		http.Error(w, "No user or group name specified", http.StatusBadRequest)
		return
	}

	// Checks if only one parameter was passed
	if userParam != "" && groupParam != "" {
		http.Error(w, "Only one parameter can be specified at a time, whether user or group", http.StatusBadRequest)
		return
	}

	if userParam != "" {

		// Search for the user on database
		user, err := repositories.FindUserByName(userParam)
		if err != nil {
			logger.Debug(err)
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}

		// Check if the token user is the same as the request
		if claims.Username != user.Username {
			logger.Debug(err)
			http.Error(w, "Unauthorized request", http.StatusForbidden)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(&client.UserResponse{Username: user.Username, PictureURL: user.Picture})

	} else {

		// Search for the group on database
		group, err := repositories.FindGroupByName(groupParam)
		if err != nil {
			logger.Debug(err)
			http.Error(w, "Group not found", http.StatusNotFound)
			return
		}

		// Check if the token user is the same as the request
		if claims.Username != group.Creator.Username {
			logger.Debug(err)
			http.Error(w, "Unauthorized request", http.StatusForbidden)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(&client.GroupResponse{Name: group.Name, PictureURL: group.Picture})
	}
}

// Handler to process profile images from users or groups
func HandleUploadPicture(w http.ResponseWriter, r *http.Request) {
	logger := configs.NewLogger("handlers")

	// Validate if authorization header is present
	authorization := r.Header.Get("Authorization")
	if authorization != "" {
		http.Error(w, "Unauthorized request", http.StatusUnauthorized)
		return
	}

	// Get and validate token
	token := authorization[strings.Index(authorization, "Bearer "):]
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

	// Get group parameter if present
	groupParam := queryParams.Get("group")

	// Check if at least one argument is present
	if userParam == "" && groupParam == "" {
		http.Error(w, "No user or group name specified", http.StatusBadRequest)
		return
	}

	// Checks if only one parameter was passed
	if userParam != "" && groupParam != "" {
		http.Error(w, "Only one parameter can be specified at a time, whether user or group", http.StatusBadRequest)
		return
	}

	if userParam != "" {

		// Check if user parameter exists
		user, err := repositories.FindUserByName(userParam)
		if err != nil {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}

		// Check if the token user is the same as the request
		if claims.Username != user.Username {
			logger.Debug(err)
			http.Error(w, "Unauthorized request", http.StatusForbidden)
			return
		}

		// Check if picture exists on database to update and delete older files from directory
		picture := user.Picture
		if picture != "" {
			err = services.DeleteFile(picture[strings.LastIndex(picture, "/"):])
			if err != nil {
				logger.Err(err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
		}
	} else {

		// Check if group parameter exists
		group, err := repositories.FindGroupByName(groupParam)
		if err != nil {
			http.Error(w, "Group not found", http.StatusNotFound)
			return
		}

		// Check if the token user is the same as the request
		if claims.Username != group.Creator.Username {
			logger.Debug(err)
			http.Error(w, "Unauthorized request", http.StatusForbidden)
			return
		}

		// Check if picture exists on database to update and delete older files from directory
		picture := group.Picture
		if picture != "" {
			err = services.DeleteFile(picture[strings.LastIndex(picture, "/"):])
			if err != nil {
				logger.Err(err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
		}
	}

	// 10MB max image size
	r.ParseMultipartForm(10 << 20)

	// Get file parameter
	_, fileHeader, err := r.FormFile("file")
	if err != nil {
		logger.Err(err)
		http.Error(w, "Failed to get file", http.StatusBadRequest)
		return
	}

	// Check if the file is a image
	fileType := fileHeader.Header.Get("Content-Type")
	if fileType != "image/jpg" && fileType != "image/jpeg" && fileType != "image/png" && fileType != "image/gif" {
		http.Error(w, "Invalid image type, only those are permitted: [jpg, jpeg, png]", http.StatusBadRequest)
		return
	}

	// Save file on directory and retrieve filename
	filename, err := saveFileAndGetFilename(fileHeader)
	if err != nil {
		logger.Err(err)
		http.Error(w, "Failed to get file", http.StatusBadRequest)
		return
	}

	// Save the generated file's URL on database
	if userParam != "" {
		_, err := repositories.SaveUserPicture(userParam, services.BuildFileURL(filename))
		if err != nil {
			logger.Err(err)
			http.Error(w, "Unable to save user picture", http.StatusInternalServerError)
			return
		}
	} else {
		_, err := repositories.SaveGroupPicture(groupParam, services.BuildFileURL(filename))
		if err != nil {
			logger.Err(err)
			http.Error(w, "Unable to save group picture", http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusCreated)
}

func saveFileAndGetFilename(fileHeader *multipart.FileHeader) (string, error) {
	_, filename, err := services.SaveFile(fileHeader)
	return filename, err
}
