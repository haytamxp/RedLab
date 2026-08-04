package dto

type LoginResponse struct {
	Token string `json:"token"`
}

type RegisterResponse struct {
	Message string `json:"message"`
}