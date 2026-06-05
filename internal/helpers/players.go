package helpers

import (
	"github.com/ffx64/frontline-stats/internal/dtos"
	"github.com/ffx64/frontline-stats/internal/entities"
)

func ToPlayerDTO(player *entities.Players) *dtos.PlayerDTO {
	var lastServerID *string
	if player.LastServerID != nil {
		idStr := player.LastServerID.String()
		lastServerID = &idStr
	}

	return &dtos.PlayerDTO{
		ID:              player.ID.String(),
		GUID:            player.GUID,
		Username:        player.Username,
		Admin:           player.Admin,
		Premium:         player.Premium,
		PremiumStartAt:  player.PremiumStartAt,
		PremiumExpireAt: player.PremiumExpireAt,
		IsActive:        player.IsActive,
		IsBanned:        player.IsBanned,
		LastServerID:    lastServerID,
		LastLogin:       player.LastLogin,
		Platform:        player.Platform,
		UpdatedAt:       player.UpdatedAt,
		CreatedAt:       player.CreatedAt,
	}
}

func ToPlayerStatsDTO(player *entities.Players, stats *entities.PlayersStats) *dtos.PlayerStatsDTO {
	var lastServerID *string
	if player.LastServerID != nil {
		idStr := player.LastServerID.String()
		lastServerID = &idStr
	}

	return &dtos.PlayerStatsDTO{
		PlayerDTO: dtos.PlayerDTO{
			ID:              player.ID.String(),
			GUID:            player.GUID,
			Username:        player.Username,
			Admin:           player.Admin,
			Premium:         player.Premium,
			PremiumStartAt:  player.PremiumStartAt,
			PremiumExpireAt: player.PremiumExpireAt,
			IsActive:        player.IsActive,
			IsBanned:        player.IsBanned,
			LastServerID:    lastServerID,
			LastLogin:       player.LastLogin,
			Platform:        player.Platform,
			UpdatedAt:       player.UpdatedAt,
			CreatedAt:       player.CreatedAt,
		},
		Stats: dtos.StatsDTO{
			Level:                stats.Level,
			XP:                   stats.XP,
			Kills:                stats.Kills,
			Deaths:               stats.Deaths,
			GrenadesThrown:       stats.GrenadesThrown,
			FriendlyFireKills:    stats.FriendlyFireKills,
			FriendlyFireDeaths:   stats.FriendlyFireDeaths,
			HeadshotsMade:        stats.HeadshotsMade,
			HeadshotsTaken:       stats.HeadshotsTaken,
			VehicleKills:         stats.VehicleKills,
			VehicleDeaths:        stats.VehicleDeaths,
			LongestKillDistance:  stats.LongestKillDistance,
			AverageKillDistance:  stats.AverageKillDistance,
			AverageDeathDistance: stats.AverageDeathDistance,
			WeaponsMostUsed:      stats.WeaponsMostUsed,
			VehicleMostUsed:      stats.VehicleMostUsed,
			HitZonesMostKilled:   stats.HitZonesMostKilled,
			HitZonesMostDied:     stats.HitZonesMostDied,
			RatioKDR:             stats.RatioKDR,
			RatioHeadshot:        stats.RatioHeadshot,
			RatioFriendlyFire:    stats.RatioFriendlyFire,
			RatioVehicle:         stats.RatioVehicle,
			MaxKillDistance:      stats.MaxKillDistance,
			UpdatedAt:            stats.UpdatedAt.String(),
			CreatedAt:            stats.CreatedAt.String(),
		},
	}
}
