package entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Kills struct {
	ID               uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	ServerID         uuid.UUID `gorm:"type:uuid;not null" json:"server_id"`
	RoundID          uuid.UUID `gorm:"type:uuid;not null" json:"round_id"`
	KillerID         uuid.UUID `gorm:"type:uuid;not null" json:"killer_id"`
	VictimID         uuid.UUID `gorm:"type:uuid;not null" json:"victim_id"`
	VictimWeaponName string    `gorm:"type:varchar(100)" json:"victim_weapon_name"`
	VictimWeaponType string    `gorm:"type:varchar(50)" json:"victim_weapon_type"`
	KillerWeaponName string    `gorm:"type:varchar(100)" json:"killer_weapon_name"`
	KillerWeaponType string    `gorm:"type:varchar(50)" json:"killer_weapon_type"`
	HitZone          string    `gorm:"type:varchar(50)" json:"hit_zone"`
	Distance         float64   `gorm:"type:float" json:"distance"`
	IsHeadshot       bool      `gorm:"default:false" json:"is_headshot"`
	IsFriendly       bool      `gorm:"default:false" json:"is_friendly"`
	IsVehicle        bool      `gorm:"default:false" json:"is_vehicle"`
	KillerTeam       string    `gorm:"type:varchar(50)" json:"killer_team"`
	VictimTeam       string    `gorm:"type:varchar(50)" json:"victim_team"`
	Timestamp        time.Time `gorm:"not null" json:"timestamp"`
	CreatedAt        time.Time `gorm:"autoCreateTime" json:"created_at"`
}

func (p *Kills) BeforeCreate(tx *gorm.DB) (err error) {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return
}
