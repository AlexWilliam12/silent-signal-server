package client

import "github.com/gorilla/websocket"

type Message struct {
	Type    string  `json:"type"`
	Content string  `json:"content"`
	Link    *string `json:"link,omitempty"`
}

type PrivateMessage struct {
	Sender   string  `json:"sender"`
	Receiver string  `json:"receiver"`
	Message  Message `json:"message"`
}

type GroupMessage struct {
	Sender  string  `json:"sender"`
	Group   string  `json:"group"`
	Message Message `json:"message"`
}

type GroupUser struct {
	Username string
	Conn     *websocket.Conn
}
