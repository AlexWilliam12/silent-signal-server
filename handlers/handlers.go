package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/AlexWilliam12/silent-signal/auth"
	"github.com/AlexWilliam12/silent-signal/client"
	"github.com/AlexWilliam12/silent-signal/configs"
	"github.com/AlexWilliam12/silent-signal/database"
	"github.com/gorilla/websocket"
)

// Handler to process the client authentication request
func HandleLogin(w http.ResponseWriter, r *http.Request) {
	logger := configs.NewLogger("handlers")

	var user client.Auth
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	if _, err := database.FindUser(user); err != nil {
		if strings.Contains(err.Error(), "record not found") {
			http.Error(w, "Invalid credentials", http.StatusNotFound)
		} else {
			logger.Err(err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	token, err := auth.GenerateToken(user.Username)
	if err != nil {
		logger.Err(err)
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(client.JWTToken{Token: token})
}

// Handler to process a request to register a user
func HandleRegister(w http.ResponseWriter, r *http.Request) {
	logger := configs.NewLogger("handlers")

	var user client.Auth
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	if _, err := database.CreateUser(user); err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			http.Error(w, "Username is already in use", http.StatusBadRequest)
		} else {
			logger.Err(err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// Upgrade the HTTP protocol to Web Socket protocol
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Handler to process bidirectional requests on private chats
func HandlePrivateChat(w http.ResponseWriter, r *http.Request) {
	logger := configs.NewLogger("handlers")

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Debugf("Failed to upgrade to web socket: %v", err)
	}
	defer conn.Close()

}

// Handler to process bidirectional requests on group chats
func HandleGroupChat(w http.ResponseWriter, r *http.Request) {
	logger := configs.NewLogger("handlers")

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Debugf("Failed to upgrade to web socket: %v", err)
	}
	defer conn.Close()
}
