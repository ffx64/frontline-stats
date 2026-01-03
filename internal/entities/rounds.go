package entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Rounds struct {
	ID            uuid.UUID  `gorm:"type:uuid;primaryKey" json:"id"`
	ServerID      uuid.UUID  `gorm:"type:uuid;not null;index" json:"server_id"`
	CurrentMode   string     `gorm:"type:text" json:"current_mode"`
	MissionHeader string     `gorm:"type:text" json:"mission_header"`
	Status        string     `gorm:"type:text" json:"status"`
	WinnerFaction string     `gorm:"type:text;default:null" json:"winner_faction"`
	EndedAt       *time.Time `gorm:"column:ended_at"`
	StartAt       time.Time  `gorm:"autoCreateTime" json:"start_at"`
	CreatedAt     time.Time  `gorm:"autoCreateTime" json:"created_at"`
}

func (r *Rounds) BeforeCreate(tx *gorm.DB) (err error) {
	if r.ID == uuid.Nil {
		r.ID = uuid.New()
	}
	return
}
