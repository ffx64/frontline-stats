package dtos

import "github.com/google/uuid"

type StatsDTO struct {
	ID                   uuid.UUID `json:"id"`
	Level                int       `json:"level"`
	XP                   int       `json:"xp"`
	Kills                int       `json:"kills"`
	Deaths               int       `json:"deaths"`
	GrenadesThrown       int       `json:"grenades_thrown"`
	FriendlyFireKills    int       `json:"friendly_fire_kills"`
	FriendlyFireDeaths   int       `json:"friendly_fire_deaths"`
	HeadshotsMade        int       `json:"headshots_made"`
	HeadshotsTaken       int       `json:"headshots_taken"`
	VehicleKills         int       `json:"vehicle_kills"`
	VehicleDeaths        int       `json:"vehicle_deaths"`
	LongestKillDistance  float32   `json:"longest_kill_distance"`
	AverageKillDistance  float32   `json:"average_kill_distance"`
	AverageDeathDistance float32   `json:"average_death_distance"`
	WeaponsMostUsed      string    `json:"weapons_most_used"`
	VehicleMostUsed      string    `json:"vehicle_most_used"`
	HitZonesMostKilled   string    `json:"hit_zones_most_killed"`
	HitZonesMostDied     string    `json:"hit_zones_most_died"`
	RatioKDR             float32   `json:"ratio_kdr"`
	RatioHeadshot        float32   `json:"ratio_headshot"`
	RatioFriendlyFire    float32   `json:"ratio_friendly_fire"`
	RatioVehicle         float32   `json:"ratio_vehicle"`
	MaxKillDistance      float32   `json:"max_kill_distance"`
	UpdatedAt            string    `json:"updated_at"`
	CreatedAt            string    `json:"created_at"`
}
