package entities_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/ffx64/frontline-stats/internal/entities"
)

func TestPlayersStatsInitialization(t *testing.T) {
	now := time.Now()
	stats := entities.PlayersStats{
		ID:                   uuid.New(),
		PlayerID:             uuid.New(),
		Level:                5,
		XP:                   1200,
		Kills:                15,
		Deaths:               3,
		GrenadesThrown:       7,
		FriendlyFireKills:    2,
		FriendlyFireDeaths:   1,
		HeadshotsMade:        8,
		HeadshotsTaken:       6,
		VehicleKills:         4,
		VehicleDeaths:        3,
		LongestKillDistance:  1200.45,
		AverageKillDistance:  350.23,
		AverageDeathDistance: 210.78,
		WeaponsMostUsed:      "AK-47",
		VehicleMostUsed:      "Tank",
		HitZonesMostKilled:   "head",
		HitZonesMostDied:     "chest",
		RatioKDR:             5.0,
		RatioHeadshot:        0.4,
		RatioFriendlyFire:    0.2,
		RatioVehicle:         0.1,
		MaxKillDistance:      250.4,
		UpdatedAt:            now,
		CreatedAt:            now,
	}

	assert.NotNil(t, stats.ID)
	assert.Equal(t, 5, stats.Level)
	assert.Equal(t, 1200, stats.XP)
	assert.Equal(t, 15, stats.Kills)
	assert.Equal(t, 3, stats.Deaths)
	assert.Equal(t, 7, stats.GrenadesThrown)
	assert.Equal(t, 2, stats.FriendlyFireKills)
	assert.Equal(t, 1, stats.FriendlyFireDeaths)
	assert.Equal(t, 8, stats.HeadshotsMade)
	assert.Equal(t, 6, stats.HeadshotsTaken)
	assert.Equal(t, 4, stats.VehicleKills)
	assert.Equal(t, 3, stats.VehicleDeaths)
	assert.Equal(t, "AK-47", stats.WeaponsMostUsed)
	assert.Equal(t, "Tank", stats.VehicleMostUsed)
	assert.Equal(t, "head", stats.HitZonesMostKilled)
	assert.Equal(t, "chest", stats.HitZonesMostDied)
	assert.Equal(t, float32(5.0), stats.RatioKDR)
	assert.WithinDuration(t, now, stats.UpdatedAt, time.Second)
	assert.WithinDuration(t, now, stats.CreatedAt, time.Second)
}

func TestPlayersStatsJSONMarshalling(t *testing.T) {
	now := time.Now()
	stats := entities.PlayersStats{
		ID:                 uuid.New(),
		PlayerID:           uuid.New(),
		Level:              10,
		XP:                 2000,
		Kills:              40,
		Deaths:             8,
		FriendlyFireKills:  3,
		FriendlyFireDeaths: 2,
		VehicleKills:       7,
		VehicleDeaths:      5,
		RatioKDR:           5.0,
		RatioHeadshot:      0.25,
		RatioVehicle:       0.3,
		WeaponsMostUsed:    "M4A1",
		HitZonesMostDied:   "chest",
		MaxKillDistance:    380.7,
		UpdatedAt:          now,
		CreatedAt:          now,
	}

	data, err := json.Marshal(stats)
	assert.NoError(t, err)
	assert.Contains(t, string(data), `"level":10`)
	assert.Contains(t, string(data), `"kills":40`)
	assert.Contains(t, string(data), `"ratio_kdr":5`)

	var unmarshalled entities.PlayersStats
	err = json.Unmarshal(data, &unmarshalled)
	assert.NoError(t, err)
	assert.Equal(t, stats.Level, unmarshalled.Level)
	assert.Equal(t, stats.Kills, unmarshalled.Kills)
	assert.Equal(t, stats.RatioKDR, unmarshalled.RatioKDR)
	assert.Equal(t, stats.WeaponsMostUsed, unmarshalled.WeaponsMostUsed)
	assert.Equal(t, stats.HitZonesMostDied, unmarshalled.HitZonesMostDied)
}

func TestPlayersStatsDefaultValues(t *testing.T) {
	now := time.Now()
	stats := entities.PlayersStats{
		PlayerID:  uuid.New(),
		CreatedAt: now,
		UpdatedAt: now,
	}

	assert.Equal(t, 0, stats.Level)
	assert.Equal(t, 0, stats.XP)
	assert.Equal(t, 0, stats.Kills)
	assert.Equal(t, 0, stats.Deaths)
	assert.Equal(t, 0, stats.GrenadesThrown)
	assert.Equal(t, 0, stats.FriendlyFireKills)
	assert.Equal(t, 0, stats.HeadshotsMade)
	assert.Equal(t, 0, stats.VehicleKills)
	assert.Equal(t, "", stats.WeaponsMostUsed)
	assert.Equal(t, "", stats.HitZonesMostKilled)
	assert.Equal(t, float32(0), stats.RatioKDR)
	assert.WithinDuration(t, now, stats.CreatedAt, time.Second)
	assert.WithinDuration(t, now, stats.UpdatedAt, time.Second)
}
