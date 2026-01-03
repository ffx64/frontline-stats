package helpers_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/ffx64/gamestats-backend/internal/entities"
	"github.com/ffx64/gamestats-backend/internal/helpers"
)

func TestToPlayerDTO_WithLastServerID(t *testing.T) {
	now := time.Now()
	lastServerUUID := uuid.New()

	player := &entities.Players{
		ID:              uuid.New(),
		GUID:            "test-guid",
		Username:        "player1",
		Admin:           1,
		Premium:         1,
		PremiumStartAt:  &now,
		PremiumExpireAt: &now,
		IsActive:        true,
		IsBanned:        false,
		LastServerID:    &lastServerUUID,
		LastLogin:       &now,
		Platform:        "pc",
		UpdatedAt:       now,
		CreatedAt:       now,
	}

	dto := helpers.ToPlayerDTO(player)

	assert.Equal(t, player.ID.String(), dto.ID)
	assert.Equal(t, player.GUID, dto.GUID)
	assert.Equal(t, player.Username, dto.Username)
	assert.NotNil(t, dto.LastServerID)
	assert.Equal(t, lastServerUUID.String(), *dto.LastServerID)
}

func TestToPlayerDTO_WithoutLastServerID(t *testing.T) {
	now := time.Now()

	player := &entities.Players{
		ID:              uuid.New(),
		GUID:            "test-guid-2",
		Username:        "player2",
		Admin:           0,
		Premium:         0,
		PremiumStartAt:  &now,
		PremiumExpireAt: &now,
		IsActive:        true,
		IsBanned:        false,
		LastServerID:    nil,
		LastLogin:       &now,
		Platform:        "pc",
		UpdatedAt:       now,
		CreatedAt:       now,
	}

	dto := helpers.ToPlayerDTO(player)

	assert.Equal(t, player.ID.String(), dto.ID)
	assert.Equal(t, player.GUID, dto.GUID)
	assert.Equal(t, player.Username, dto.Username)
	assert.Nil(t, dto.LastServerID)
}
