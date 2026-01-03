package entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RoundsStats struct {
	ID                  uuid.UUID `gorm:"primaryKey;default:gen_random_uuid()" json:"id"`
	RoundID             uuid.UUID `gorm:"round_id;type:uuid;not null" json:"round_id"`
	ServerID            uuid.UUID `gorm:"server_id;type:uuid;not null" json:"server_id"`
	PlayerID            uuid.UUID `gorm:"player_id;type:uuid;not null" json:"player_id"`
	Team                string    `gorm:"team;type:varchar(20);default:''" json:"team"`
	TotalKills          int64     `gorm:"total_kills;default:0" json:"total_kills"`
	TotalDeaths         int64     `gorm:"total_deaths;default:0" json:"total_deaths"`
	TotalSuicides       int64     `gorm:"total_suicides;default:0" json:"total_suicides"`
	AverageKillDistance float64   `gorm:"average_kill_distance;default:0" json:"average_kill_distance"`
	TotalTeamKills      int64     `gorm:"total_team_kills;default:0" json:"total_team_kills"`
	TotalHeadshots      int64     `gorm:"total_headshots;default:0" json:"total_headshots"`
	TotalVehicleKills   int64     `gorm:"total_vehicle_kills;default:0" json:"total_vehicle_kills"`
	MostUsedWeapon      string    `gorm:"most_used_weapon;default:''" json:"most_used_weapon"`
	MostHitZone         string    `gorm:"most_hit_zone;default:''" json:"most_hit_zone"`
	UpdatedAt           time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	CreatedAt           time.Time `gorm:"autoCreateTime" json:"created_at"`
}

func (r *RoundsStats) BeforeCreate(tx *gorm.DB) (err error) {
	if r.ID == uuid.Nil {
		r.ID = uuid.New()
	}
	return
}
