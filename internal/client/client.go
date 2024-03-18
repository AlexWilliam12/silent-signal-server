package client

import "github.com/gorilla/websocket"

type UserRequest struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type GroupRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type UserResponse struct {
	Email      string `json:"email"`
	Username   string `json:"username"`
	PictureURL string `json:"pictureURL"`
}

type GroupResponse struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Creator     string `json:"creator"`
	PictureURL  string `json:"pictureURL"`
}

type JWTToken struct {
	Token string `json:"token"`
}

type PrivateMessage struct {
	Sender   string `json:"sender"`
	Receiver string `json:"receiver"`
	Message  string `json:"message"`
}

type GroupMessage struct {
	Sender  string `json:"sender"`
	Group   string `json:"group"`
	Message string `json:"message"`
}

type ContactRequest struct {
	Contact string `json:"contact"`
}

type GroupUser struct {
	Username string
	Conn     *websocket.Conn
}
