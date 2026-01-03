package entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PlayersStats struct {
	ID                   uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	PlayerID             uuid.UUID `gorm:"type:uuid;not null" json:"player_id"`
	Level                int       `gorm:"type:int;not null;default:1" json:"level"`
	XP                   int       `gorm:"type:int;not null;default:0" json:"xp"`
	Kills                int       `gorm:"type:int;not null;default:0" json:"kills"`
	Deaths               int       `gorm:"type:int;not null;default:0" json:"deaths"`
	GrenadesThrown       int       `gorm:"type:int;not null;default:0" json:"grenades_thrown"`
	FriendlyFireKills    int       `gorm:"type:int;not null;default:0" json:"friendly_fire_kills"`
	FriendlyFireDeaths   int       `gorm:"type:int;not null;default:0" json:"friendly_fire_deaths"`
	HeadshotsMade        int       `gorm:"type:int;not null;default:0" json:"headshots_made"`
	HeadshotsTaken       int       `gorm:"type:int;not null;default:0" json:"headshots_taken"`
	VehicleKills         int       `gorm:"type:int;not null;default:0" json:"vehicle_kills"`
	VehicleDeaths        int       `gorm:"type:int;not null;default:0" json:"vehicle_deaths"`
	LongestKillDistance  float32   `gorm:"type:real;not null;default:0" json:"longest_kill_distance"`
	AverageKillDistance  float32   `gorm:"type:real;not null;default:0" json:"average_kill_distance"`
	AverageDeathDistance float32   `gorm:"type:real;not null;default:0" json:"average_death_distance"`
	WeaponsMostUsed      string    `gorm:"type:text;not null;default:''" json:"weapons_most_used"`
	VehicleMostUsed      string    `gorm:"type:text;not null;default:''" json:"vehicle_most_used"`
	HitZonesMostKilled   string    `gorm:"type:text;not null;default:''" json:"hit_zones_most_killed"`
	HitZonesMostDied     string    `gorm:"type:text;not null;default:''" json:"hit_zones_most_died"`
	RatioKDR             float32   `gorm:"type:real;not null;default:0" json:"ratio_kdr"`
	RatioHeadshot        float32   `gorm:"type:real;not null;default:0" json:"ratio_headshot"`
	RatioFriendlyFire    float32   `gorm:"type:real;not null;default:0" json:"ratio_friendly_fire"`
	RatioVehicle         float32   `gorm:"type:real;not null;default:0" json:"ratio_vehicle"`
	MaxKillDistance      float32   `gorm:"type:real;not null;default:0" json:"max_kill_distance"`
	UpdatedAt            time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	CreatedAt            time.Time `gorm:"autoCreateTime" json:"created_at"`
}

func (p *PlayersStats) BeforeCreate(tx *gorm.DB) (err error) {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return
}
