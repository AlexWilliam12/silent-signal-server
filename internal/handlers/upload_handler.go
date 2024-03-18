package handlers

import (
	"encoding/json"
	"log"
	"mime/multipart"
	"net/http"
	"strings"

	"github.com/AlexWilliam12/silent-signal/internal/client"
	"github.com/AlexWilliam12/silent-signal/internal/database/models"
	"github.com/AlexWilliam12/silent-signal/internal/database/repositories"
	"github.com/AlexWilliam12/silent-signal/internal/services"
)

func HandleFetchPicture(w http.ResponseWriter, r *http.Request) {

	claims := handleAuthorization(w, r)

	queryParams := r.URL.Query()

	groupParam := queryParams.Get("group")

	if groupParam != "" {

		group, err := repositories.FindGroupByName(groupParam)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		if claims.Username != group.Creator.Username {
			http.Error(w, "Unauthorized request", http.StatusForbidden)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(&client.GroupResponse{Name: group.Name, PictureURL: group.Picture})

	} else {

		user, err := repositories.FindUserByName(claims.Username)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(&client.UserResponse{Username: user.Username, PictureURL: user.Picture})
	}
}

func HandleUploadPicture(w http.ResponseWriter, r *http.Request) {

	claims := handleAuthorization(w, r)

	queryParams := r.URL.Query()

	groupParam := queryParams.Get("group")

	if groupParam != "" {

		group, err := repositories.FindGroupByName(groupParam)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		if claims.Username != group.Creator.Username {
			http.Error(w, "Unauthorized request", http.StatusForbidden)
			return
		}

		if group.Picture != "" {
			if err = services.DeleteFile(group.Picture[strings.LastIndex(group.Picture, "/"):]); err != nil {
				log.Println(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

	} else {

		user, err := repositories.FindUserByName(claims.Username)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		if user.Picture != "" {
			if err = services.DeleteFile(user.Picture[strings.LastIndex(user.Picture, "/"):]); err != nil {
				log.Println(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

	}

	// 10MB max image size
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, fileHeader, err := r.FormFile("file")
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fileType := fileHeader.Header.Get("Content-Type")
	if fileType != "image/jpg" && fileType != "image/jpeg" && fileType != "image/png" && fileType != "image/gif" {
		http.Error(w, "Invalid image type, only those are permitted: [jpg, jpeg, png]", http.StatusBadRequest)
		return
	}

	filename, err := saveFileAndGetFilename(fileHeader)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if groupParam != "" {
		if _, err := repositories.SaveGroupPicture(groupParam, services.BuildFileURL(filename)); err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		if _, err := repositories.SaveUserPicture(claims.Username, services.BuildFileURL(filename)); err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	if r.Method == "POST" {
		w.WriteHeader(http.StatusCreated)
	} else {
		w.WriteHeader(http.StatusNoContent)
	}
}

func HandleChatUpload(w http.ResponseWriter, r *http.Request) {

	claims := handleAuthorization(w, r)

	queryParams := r.URL.Query()

	groupParam := queryParams.Get("group")

	receiverParam := queryParams.Get("receiver")

	if receiverParam == "" || groupParam == "" {
		http.Error(w, "No receiver or group name specified", http.StatusBadRequest)
		return
	}

	if receiverParam != "" && groupParam != "" {
		http.Error(w, "Only one parameter can be specified at a time, whether receiver or group", http.StatusBadRequest)
		return
	}

	var receiver *models.User
	var group *models.Group

	if groupParam != "" {
		fetchedGroup, err := repositories.FindGroupByName(groupParam)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		group = fetchedGroup
	} else {
		user, err := repositories.FindUserByName(receiverParam)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		receiver = user
	}

	// 100MB max file size
	if err := r.ParseMultipartForm(100 << 20); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	filename, err := saveFileAndGetFilename(fileHeader)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fileType := fileHeader.Header.Get("Content-Type")

	if groupParam != "" {
		sender, err := repositories.FindUserByName(claims.Username)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		fileUrl := services.BuildFileURL(filename)
		_, err = repositories.SaveGroupMessage(&models.GroupMessage{
			Sender:  *sender,
			Group:   *group,
			Type:    fileType,
			Content: fileUrl,
		})
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var encodedFile string

		switch fileType {
		case
			"image/jpg",
			"image/jpeg",
			"image/gif",
			"image/png",
			"image/webp",
			"image/svg+xml",
			"image/bmp",
			"image/tiff":
			if encodedFile, err = services.EncodedImgThumbnail(&file); err != nil {
				log.Println(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		case
			"video/mp4",
			"video/x-msvideo",
			"video/quicktime",
			"video/x-ms-wmv",
			"video/x-matroska",
			"video/x-flv",
			"video/webm",
			"video/3gpp":
			if encodedFile, err = services.EncodedVideoThumbnail(&file, fileType[strings.Index(fileType, "/"):]); err != nil {
				log.Println(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		if conns, ok := groupConns[groupParam]; ok {
			for _, conn := range conns {
				if conn.Username != claims.Username {
					conn.Conn.WriteJSON(client.GroupMessage{
						Sender: sender.Username,
						Group:  group.Name,
						Message: client.Message{
							Type:    fileType,
							Content: encodedFile,
							Link:    &fileUrl,
						},
					})
				}
			}
		}

	} else {
		sender, err := repositories.FindUserByName(claims.Username)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		_, err = repositories.SavePrivateMessage(&models.PrivateMessage{
			Sender:   *sender,
			Receiver: *receiver,
			Type:     fileType,
			Content:  services.BuildFileURL(filename),
		})
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusCreated)
}

func saveFileAndGetFilename(fileHeader *multipart.FileHeader) (string, error) {
	filename, err := services.SaveFile(fileHeader)
	return filename, err
}
