package client

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
