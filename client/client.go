package client

import "github.com/gorilla/websocket"

// Client authentication request
type Auth struct {
	Username string `json:"username"`
	Password string `json:"password"`
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
