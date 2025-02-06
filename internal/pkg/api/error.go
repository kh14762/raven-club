package api

import (
	"fmt"
	"net/http"
)

// Error is used to define a common structured error response for all API calls.
type Error struct {
	Message  string `json:"message"`
	HttpCode int    `json:"httpCode"`
}

func NewError(msg string) Error {
	return Error{msg, http.StatusInternalServerError}
}

// NewErrorf Wraps fmt.Spritnf for convenience
func NewErrorf(msg string, args ...any) Error {
	return Error{fmt.Sprintf(msg, args...), http.StatusInternalServerError}
}

// WrapError convenience method to convert Error to string
func WrapError(err error) Error {
	return Error{err.Error(), http.StatusInternalServerError}
}
