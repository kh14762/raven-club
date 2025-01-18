package types

import "github.com/gorilla/sessions"

type RegisterResponse struct {
	Success bool              `json:"success"`
	Message string            `json:"message"`
	Session *sessions.Session `json:"session"`
}

type LoginResponse struct {
	Success bool              `json:"success"`
	Message string            `json:"message"`
	Session *sessions.Session `json:"session"`
}
