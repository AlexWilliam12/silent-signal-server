package handlers

import (
	"log"
	"net/http"

	"github.com/AlexWilliam12/silent-signal/internal/client"
	"github.com/AlexWilliam12/silent-signal/internal/database/models"
	"github.com/AlexWilliam12/silent-signal/internal/database/repositories"
	"github.com/gorilla/websocket"
)

var groupConns = make(map[string][]*client.GroupUser)

func HandleGroupMessages(w http.ResponseWriter, r *http.Request) {

	claims := handleAuthorization(w, r)

	queryParams := r.URL.Query()

	groupParam := queryParams.Get("name")

	if groupParam == "" {
		http.Error(w, "Group name not specified", http.StatusBadRequest)
		return
	}

	group, err := repositories.FindGroupByName(groupParam)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	if conns, ok := groupConns[groupParam]; ok {
		var isPresent bool
		for _, groupConn := range conns {
			if groupConn.Conn == conn {
				isPresent = true
				break
			}
		}
		if !isPresent {
			groupConns[groupParam] = append(groupConns[groupParam], &client.GroupUser{Username: claims.Username, Conn: conn})
		}
	} else {
		groupConns[groupParam] = append(groupConns[groupParam], &client.GroupUser{Username: claims.Username, Conn: conn})
	}

	go sendPendingGroupMessages(claims.Username, conn)

	for {
		var message client.GroupMessage
		if err := conn.ReadJSON(&message); err != nil {
			log.Println(err)
			break
		}

		if conns, ok := groupConns[message.Group]; ok {
			for _, groupConn := range conns {
				if groupConn.Conn != conn {
					if err := groupConn.Conn.WriteJSON(&message); err != nil {
						log.Println(err)
						groupConn.Conn.Close()
						delete(privateConns, message.Group)
						break
					}
				}
			}
			go saveGroupMessages(&message, group)
		}
	}
}

func saveGroupMessages(message *client.GroupMessage, group *models.Group) {

	sender, err := repositories.FindUserByName(message.Sender)
	if err != nil {
		log.Println(err)
		return
	}

	var usernames []string
	for _, conn := range groupConns[group.Name] {
		usernames = append(usernames, conn.Username)
	}

	users, err := repositories.FetchAllByUserames(usernames)
	if err != nil {
		log.Println(err)
		return
	}

	_, err = repositories.SaveGroupMessage(&models.GroupMessage{
		Sender:  *sender,
		Group:   *group,
		Type:    "text",
		Content: message.Message.Content,
		SeenBy:  users,
	})
	if err != nil {
		log.Println(err)
		return
	}
}

func sendPendingGroupMessages(username string, conn *websocket.Conn) {

	messages, err := repositories.FetchPendingGroupMessages(username)
	if err != nil {
		log.Println(err)
		return
	}

	for _, message := range messages {
		err := conn.WriteJSON(client.GroupMessage{
			Sender: message.Sender.Username,
			Group:  message.Group.Name,
			Message: client.Message{
				Type:    "text",
				Content: message.Content,
			},
		})
		if err != nil {
			log.Println(err)
			continue
		}
	}
}
