package client

type AuthRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserRequest struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type GroupRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type ContactRequest struct {
	Contact string `json:"contact"`
}
