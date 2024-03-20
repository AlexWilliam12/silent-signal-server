package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/AlexWilliam12/silent-signal/internal/auth"
	"github.com/AlexWilliam12/silent-signal/internal/client"
	"github.com/AlexWilliam12/silent-signal/internal/database/repositories"
)

func HandleLogin(w http.ResponseWriter, r *http.Request) {

	var request client.AuthRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := repositories.FindUserByCredentials(&request)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	token, err := auth.GenerateToken(user.Username)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(client.JWTToken{Token: token})
}

func HandleRegister(w http.ResponseWriter, r *http.Request) {

	var user client.UserRequest
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if _, err := repositories.CreateUser(user); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
