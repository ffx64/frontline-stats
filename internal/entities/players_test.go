package entities_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/ffx64/frontline-stats/internal/entities"
)

func TestPlayerInitialization(t *testing.T) {
	now := time.Now()
	serverID := uuid.New()
	player := entities.Players{
		ID:              uuid.New(),
		GUID:            "GUID_123",
		Username:        "test_player",
		Admin:           1,
		Premium:         1,
		PremiumStartAt:  &now,
		PremiumExpireAt: &now,
		IsActive:        true,
		IsBanned:        false,
		LastServerID:    &serverID,
		LastLogin:       &now,
		Platform:        "steam",
		UpdatedAt:       now,
		CreatedAt:       now,
	}

	assert.NotNil(t, player.ID)
	assert.Equal(t, "test_player", player.Username)
	assert.Equal(t, "steam", player.Platform)
	assert.True(t, player.IsActive)
	assert.False(t, player.IsBanned)
	assert.NotNil(t, player.PremiumStartAt)
	assert.NotNil(t, player.LastServerID)
}

func TestPlayerJSONMarshalling(t *testing.T) {
	now := time.Now()
	player := entities.Players{
		ID:             uuid.New(),
		GUID:           "GUID_ABC",
		Username:       "json_test",
		Admin:          0,
		Premium:        0,
		PremiumStartAt: &now,
		IsActive:       true,
		IsBanned:       false,
		Platform:       "epic",
		UpdatedAt:      now,
		CreatedAt:      now,
	}

	data, err := json.Marshal(player)
	assert.NoError(t, err)
	assert.Contains(t, string(data), `"username":"json_test"`)

	var unmarshalled entities.Players
	err = json.Unmarshal(data, &unmarshalled)
	assert.NoError(t, err)
	assert.Equal(t, player.Username, unmarshalled.Username)
	assert.Equal(t, player.Platform, unmarshalled.Platform)
}

func TestPlayerDefaultValues(t *testing.T) {
	player := entities.Players{
		GUID:      "GUID_DEFAULT",
		Username:  "default_user",
		IsActive:  true,
		IsBanned:  false,
		Platform:  "console",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	assert.Equal(t, "default_user", player.Username)
	assert.Equal(t, "console", player.Platform)
	assert.True(t, player.IsActive)
	assert.False(t, player.IsBanned)
}
