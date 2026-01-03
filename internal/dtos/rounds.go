package dtos

import (
	"time"
)

type RoundsDTO struct {
	ID            string     `json:"id"`
	ServerID      string     `json:"server_id"`
	CurrentMode   string     `json:"current_mode"`
	MissionHeader string     `json:"mission_header"`
	Status        string     `json:"status"`
	WinnerFaction string     `json:"winner_faction"`
	EndedAt       *time.Time `json:"ended_at,omitempty"`
	StartAt       time.Time  `json:"start_at"`
	CreatedAt     time.Time  `json:"created_at"`
}

type RoundsCreateDTO struct {
	ServerID      string `json:"server_id" validate:"required"`
	CurrentMode   string `json:"current_mode" validate:"required"`
	MissionHeader string `json:"mission_header" validate:"required"`
	Status        string `json:"status" validate:"required"`
	WinnerFaction string `json:"winner_faction"`
}

type RoundsUpdatedEndedDTO struct {
	WinnerFaction string `json:"winner_faction"`
}

type RoundsDTOs struct {
	Total int64       `json:"total"`
	Data  []RoundsDTO `json:"data"`
}
