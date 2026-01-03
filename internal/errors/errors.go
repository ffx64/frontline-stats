package errors

import (
	"time"
)

type AppError struct {
	Status    int       `json:"status"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
}

func (e *AppError) Error() string {
	return e.Message
}

func New(message string, status int) *AppError {
	return &AppError{
		Message:   message,
		Status:    status,
		Timestamp: time.Now().UTC(),
	}
}
