package dtos

import "time"

type ServersDTO struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

type ServersSaveDTO struct {
	Name string `json:"name" validate:"required"`
}
