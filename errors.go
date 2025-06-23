package main

import "github.com/danielgtaylor/huma/v2"

type Error struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func (e *Error) Error() string {
	return e.Message
}

func (e *Error) GetStatus() int {
	return e.Status
}

func setErrorModel(_ huma.API) {
	huma.NewError = func(status int, message string, _ ...error) huma.StatusError {
		return &Error{
			Status:  status,
			Message: message,
		}
	}
}
