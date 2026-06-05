package entities_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/ffx64/frontline-stats/internal/entities"
)

func TestKillInitialization(t *testing.T) {
	now := time.Now()
	kill := entities.Kills{
		ID:               uuid.New(),
		ServerID:         uuid.New(),
		RoundID:          uuid.New(),
		KillerID:         uuid.New(),
		VictimID:         uuid.New(),
		VictimWeaponName: "AK-47",
		VictimWeaponType: "rifle",
		KillerWeaponName: "M4A1",
		KillerWeaponType: "rifle",
		HitZone:          "head",
		Distance:         150.5,
		IsHeadshot:       true,
		IsFriendly:       false,
		IsVehicle:        false,
		KillerTeam:       "blue",
		VictimTeam:       "red",
		Timestamp:        now,
		CreatedAt:        now,
	}

	assert.NotNil(t, kill.ID)
	assert.Equal(t, "AK-47", kill.VictimWeaponName)
	assert.Equal(t, "rifle", kill.KillerWeaponType)
	assert.Equal(t, "head", kill.HitZone)
	assert.True(t, kill.IsHeadshot)
	assert.False(t, kill.IsFriendly)
	assert.Equal(t, "blue", kill.KillerTeam)
	assert.Equal(t, "red", kill.VictimTeam)
	assert.WithinDuration(t, now, kill.Timestamp, time.Second)
}

func TestKillJSONMarshalling(t *testing.T) {
	now := time.Now()
	kill := entities.Kills{
		ID:               uuid.New(),
		ServerID:         uuid.New(),
		RoundID:          uuid.New(),
		KillerID:         uuid.New(),
		VictimID:         uuid.New(),
		VictimWeaponName: "AKS-74U",
		VictimWeaponType: "smg",
		KillerWeaponName: "G36",
		KillerWeaponType: "rifle",
		HitZone:          "chest",
		Distance:         75.2,
		IsHeadshot:       false,
		IsFriendly:       false,
		IsVehicle:        false,
		KillerTeam:       "red",
		VictimTeam:       "blue",
		Timestamp:        now,
		CreatedAt:        now,
	}

	data, err := json.Marshal(kill)
	assert.NoError(t, err)
	assert.Contains(t, string(data), `"killer_weapon_name":"G36"`)

	var unmarshalled entities.Kills
	err = json.Unmarshal(data, &unmarshalled)
	assert.NoError(t, err)
	assert.Equal(t, kill.KillerWeaponName, unmarshalled.KillerWeaponName)
	assert.Equal(t, kill.VictimWeaponType, unmarshalled.VictimWeaponType)
	assert.Equal(t, kill.HitZone, unmarshalled.HitZone)
}

func TestKillDefaultValues(t *testing.T) {
	now := time.Now()
	kill := entities.Kills{
		ServerID:  uuid.New(),
		RoundID:   uuid.New(),
		KillerID:  uuid.New(),
		VictimID:  uuid.New(),
		Timestamp: now,
		CreatedAt: now,
	}

	assert.False(t, kill.IsHeadshot)
	assert.False(t, kill.IsFriendly)
	assert.False(t, kill.IsVehicle)
	assert.Empty(t, kill.KillerTeam)
	assert.Empty(t, kill.VictimTeam)
	assert.WithinDuration(t, now, kill.Timestamp, time.Second)
}
