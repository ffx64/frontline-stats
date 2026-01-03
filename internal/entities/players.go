package entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Players struct {
	ID              uuid.UUID  `gorm:"type:uuid;primaryKey" json:"id"`
	GUID            string     `gorm:"uniqueIndex;not null" json:"guid"`
	Username        string     `gorm:"not null" json:"username"`
	Admin           int        `gorm:"default:0" json:"admin"`
	Premium         int        `gorm:"default:0" json:"premium"`
	PremiumStartAt  *time.Time `json:"premium_start_at,omitempty"`
	PremiumExpireAt *time.Time `json:"premium_expire_at,omitempty"`
	IsActive        bool       `gorm:"default:true" json:"is_active"`
	IsBanned        bool       `gorm:"default:false" json:"is_banned"`
	LastServerID    *uuid.UUID `gorm:"type:uuid" json:"last_server_id,omitempty"`
	LastLogin       *time.Time `json:"last_login,omitempty"`
	Platform        string     `gorm:"not null" json:"platform"`
	UpdatedAt       time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
	CreatedAt       time.Time  `gorm:"autoCreateTime" json:"created_at"`
}

func (p *Players) BeforeCreate(tx *gorm.DB) (err error) {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return
}
