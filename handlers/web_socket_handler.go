package handlers

import (
	"net/http"

	"github.com/AlexWilliam12/silent-signal/configs"
	"github.com/gorilla/websocket"
)

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
