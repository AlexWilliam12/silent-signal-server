package client

// JSON user request object
type UserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// JSON group request object
type GroupRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Creator     string `json:"creator"`
}

// JSON user response object
type UserResponse struct {
	Username   string `json:"username"`
	PictureURL string `json:"pictureURL"`
}

// JSON user response object
type GroupResponse struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Creator     string `json:"creator"`
	PictureURL  string `json:"pictureURL"`
}

// Token response to the client
type JWTToken struct {
	Token string `json:"token"`
}

type PrivateMessage struct {
	Sender   string `json:"sender"`
	Receiver string `json:"receiver"`
	Data     string `json:"data"`
}

type GroupMessage struct {
	Sender string `json:"sender"`
	Group  string `json:"group"`
	Data   string `json:"data"`
}

// type WebSocketError struct {
// 	Err string `json:"error"`
// }
