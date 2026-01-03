package dtos

import "time"

type PlayerDTO struct {
	ID              string     `json:"id"` // uuid como string
	GUID            string     `json:"guid"`
	Username        string     `json:"username"`
	Admin           int        `json:"admin"`
	Premium         int        `json:"premium"`
	PremiumStartAt  *time.Time `json:"premium_start_at,omitempty"`
	PremiumExpireAt *time.Time `json:"premium_expire_at,omitempty"`
	IsActive        bool       `json:"is_active"`
	IsBanned        bool       `json:"is_banned"`
	LastServerID    *string    `json:"last_server_id,omitempty"` // uuid como string
	LastLogin       *time.Time `json:"last_login,omitempty"`
	Platform        string     `json:"platform"`
	UpdatedAt       time.Time  `json:"updated_at"`
	CreatedAt       time.Time  `json:"created_at"`
}

type PlayerSaveDTO struct {
	GUID         string `json:"guid" binding:"required" validate:"required"`
	Username     string `json:"username" binding:"required" validate:"required"`
	LastServerID string `json:"last_server_id" binding:"required" validate:"required"` // uuid como string
	Platform     string `json:"platform" binding:"required" validate:"required"`
}

type PlayerUpdateDTO struct {
	Username     string `json:"username"`
	LastServerID string `json:"last_server_id"`
	IsActive     bool   `json:"is_active"`
}

type PlayerStatsDTO struct {
	PlayerDTO `json:",inline"`
	Stats     StatsDTO `json:"stats"`
}
