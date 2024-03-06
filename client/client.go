package client

import "github.com/gorilla/websocket"

// JSON user request object
type UserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// JSON group request object
type GroupRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// JSON user response object
type UserResponse struct {
	Username   string `json:"username"`
	PictureURL string `json:"pictureURL"`
}

// JSON user response object
type GroupResponse struct {
	Name       string `json:"name"`
	PictureURL string `json:"pictureURL"`
}

// Token response to the client
type JWTToken struct {
	Token string `json:"token"`
}

type Connection struct {
	UserID  string
	GroupID string
	Conn    *websocket.Conn
}
