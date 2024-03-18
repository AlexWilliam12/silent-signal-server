package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/AlexWilliam12/silent-signal/internal/client"
	"github.com/AlexWilliam12/silent-signal/internal/database/repositories"
)

func HandleUserUpdate(w http.ResponseWriter, r *http.Request) {

	claims := handleAuthorization(w, r)

	var request client.UserRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if _, err := repositories.UpdateUser(claims.Username, &request); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func HandleUserDelete(w http.ResponseWriter, r *http.Request) {

	claims := handleAuthorization(w, r)

	if _, err := repositories.DeleteUserByName(claims.Username); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func HandleSaveContact(w http.ResponseWriter, r *http.Request) {

	claims := handleAuthorization(w, r)

	user, err := repositories.FindUserByName(claims.Username)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	var request client.ContactRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	contact, err := repositories.FindUserByName(request.Contact)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if _, err = repositories.SaveContact(user, contact); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
