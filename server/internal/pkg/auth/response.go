package auth

import "github.com/gorilla/sessions"

type Response struct {
	Success bool              `json:"success"`
	Message string            `json:"message"`
	Session *sessions.Session `json:"session"`
}
