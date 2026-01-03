package dtos

type RoundsStatsDTO struct {
	ID                  string  `json:"id"`
	RoundID             string  `json:"round_id"`
	ServerID            string  `json:"server_id"`
	PlayerID            string  `json:"player_id"`
	Team                string  `json:"team"`
	TotalKills          int64   `json:"total_kills"`
	TotalDeaths         int64   `json:"total_deaths"`
	TotalSuicides       int64   `json:"total_suicides"`
	AverageKillDistance float64 `json:"average_kill_distance"`
	TotalTeamKills      int64   `json:"total_team_kills"`
	TotalHeadshots      int64   `json:"total_headshots"`
	TotalVehicleKills   int64   `json:"total_vehicle_kills"`
	MostUsedWeapon      string  `json:"most_used_weapon"`
	MostHitZone         string  `json:"most_hit_zone"`
	CreatedAt           int64   `json:"created_at"`
	UpdatedAt           int64   `json:"updated_at"`
}
