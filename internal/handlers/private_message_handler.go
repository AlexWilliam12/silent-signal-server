package handlers

import (
	"log"
	"net/http"

	"github.com/AlexWilliam12/silent-signal/internal/client"
	"github.com/AlexWilliam12/silent-signal/internal/database/models"
	"github.com/AlexWilliam12/silent-signal/internal/database/repositories"
	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	privateConns     = make(map[string]*websocket.Conn)
	privateBroadcast = make(chan client.PrivateMessage)
)

func HandlePrivateConnections(w http.ResponseWriter, r *http.Request) {

	claims := handleAuthorization(w, r)

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	if _, ok := privateConns[claims.Username]; !ok {
		privateConns[claims.Username] = conn
	}

	go sendPendingPrivateMessages(claims.Username)

	for {
		var message client.PrivateMessage
		if err := conn.ReadJSON(&message); err != nil {
			log.Println(err)
			conn.Close()
			delete(privateConns, claims.Username)
			break
		}

		privateBroadcast <- message
	}
}

func HandlePrivateMessages() {
	for {
		message := <-privateBroadcast
		if conn, ok := privateConns[message.Receiver]; ok {
			if err := conn.WriteJSON(&message); err != nil {
				log.Println(err)
				conn.Close()
				delete(privateConns, message.Receiver)
				go savePrivateMessages(&message, true)
				break
			}
			go savePrivateMessages(&message, false)
		} else {
			go savePrivateMessages(&message, true)
		}
	}
}

func savePrivateMessages(message *client.PrivateMessage, isPending bool) {
	sender, err := repositories.FindUserByName(message.Sender)
	if err != nil {
		log.Println(err)
		return
	}

	receiver, err := repositories.FindUserByName(message.Receiver)
	if err != nil {
		log.Println(err)
		return
	}

	_, err = repositories.SavePrivateMessage(&models.PrivateMessage{
		Sender:   *sender,
		Receiver: *receiver,
		Type:     "text", Content: message.Message,
		IsPending: isPending,
	})
	if err != nil {
		log.Println(err)
		return
	}
}

func sendPendingPrivateMessages(username string) {
	messages, err := repositories.FetchPendingPrivateMessages(username)
	if err != nil {
		log.Println(err)
		return
	}

	if len(messages) == 0 {
		return
	}

	if conn, ok := privateConns[username]; ok {
		var ids []uint
		for _, message := range messages {
			response := client.PrivateMessage{
				Sender:   message.Sender.Username,
				Receiver: message.Receiver.Username,
				Message:  message.Content,
			}
			if err = conn.WriteJSON(&response); err != nil {
				log.Println(err)
				continue
			}
			ids = append(ids, message.ID)
		}
		go repositories.UpdatePendingSituation(ids)
	}
}
