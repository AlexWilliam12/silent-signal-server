package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/AlexWilliam12/silent-signal/internal/client"
	"github.com/AlexWilliam12/silent-signal/internal/database/repositories"
)

func HandleCreateGroup(w http.ResponseWriter, r *http.Request) {

	claims := handleAuthorization(w, r)

	var request client.GroupRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	creator, err := repositories.FindUserByName(claims.Username)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if _, err = repositories.CreateGroup(&request, creator); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func HandleFetchAllGroups(w http.ResponseWriter, r *http.Request) {

	handleAuthorization(w, r)

	groups, err := repositories.FindAllGroups()
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

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

func HandleFetchGroup(w http.ResponseWriter, r *http.Request) {

	handleAuthorization(w, r)

	queryParams := r.URL.Query()

	groupParam := queryParams.Get("group")

	if groupParam == "" {
		http.Error(w, "No group name specified", http.StatusBadRequest)
		return
	}

	group, err := repositories.FindGroupByName(groupParam)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&group)
}

func HandleUpdateGroup(w http.ResponseWriter, r *http.Request) {

	claims := handleAuthorization(w, r)

	queryParams := r.URL.Query()

	groupParam := queryParams.Get("group")

	if groupParam == "" {
		http.Error(w, "No group name specified", http.StatusBadRequest)
		return
	}

	var request client.GroupRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	group, err := repositories.FindGroupByName(groupParam)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if group.Creator.Username != claims.Username {
		http.Error(w, "Unauthorized request", http.StatusForbidden)
		return
	}

	if _, err = repositories.UpdateGroup(&request, group); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func HandleDeleteGroup(w http.ResponseWriter, r *http.Request) {

	claims := handleAuthorization(w, r)

	queryParams := r.URL.Query()

	groupParam := queryParams.Get("group")

	if groupParam == "" {
		http.Error(w, "No group name specified", http.StatusBadRequest)
		return
	}

	group, err := repositories.FindGroupByName(groupParam)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if group.Creator.Username != claims.Username {
		http.Error(w, "Unauthorized request", http.StatusForbidden)
		return
	}

	if _, err = repositories.DeleteGroupByName(group.Name); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
