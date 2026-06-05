package helpers_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/ffx64/frontline-stats/internal/entities"
	"github.com/ffx64/frontline-stats/internal/helpers"
)

func TestToRoundsDTO_WithEndedAt(t *testing.T) {
	now := time.Now()
	serverID := uuid.New()
	roundID := uuid.New()

	endedAt := now.Add(10 * time.Minute)
	round := &entities.Rounds{
		ID:            roundID,
		ServerID:      serverID,
		CurrentMode:   "CTI",
		MissionHeader: "Operation Thunder",
		Status:        "active",
		WinnerFaction: "NATO",
		EndedAt:       &endedAt,
		StartAt:       now,
		CreatedAt:     now,
	}

	dto := helpers.ToRoundsDTO(round)

	assert.NotNil(t, dto)
	assert.Equal(t, roundID.String(), dto.ID)
	assert.Equal(t, serverID.String(), dto.ServerID)
	assert.Equal(t, "CTI", dto.CurrentMode)
	assert.Equal(t, "Operation Thunder", dto.MissionHeader)
	assert.Equal(t, "active", dto.Status)
	assert.Equal(t, "NATO", dto.WinnerFaction)
	assert.NotNil(t, dto.EndedAt)
	assert.WithinDuration(t, endedAt, *dto.EndedAt, time.Second)
	assert.WithinDuration(t, now, dto.StartAt, time.Second)
	assert.WithinDuration(t, now, dto.CreatedAt, time.Second)
}

func TestToRoundsDTO_EndedAtNil(t *testing.T) {
	now := time.Now()
	serverID := uuid.New()
	roundID := uuid.New()

	round := &entities.Rounds{
		ID:            roundID,
		ServerID:      serverID,
		CurrentMode:   "AAS",
		MissionHeader: "Night Raid",
		Status:        "waiting",
		WinnerFaction: "",
		EndedAt:       nil, // EndedAt nulo
		StartAt:       now,
		CreatedAt:     now,
	}

	dto := helpers.ToRoundsDTO(round)

	assert.NotNil(t, dto)
	assert.Equal(t, roundID.String(), dto.ID)
	assert.Equal(t, serverID.String(), dto.ServerID)
	assert.Nil(t, dto.EndedAt, "EndedAt deve ser nil")
	assert.Equal(t, "AAS", dto.CurrentMode)
	assert.Equal(t, "Night Raid", dto.MissionHeader)
	assert.Equal(t, "waiting", dto.Status)
	assert.Equal(t, "", dto.WinnerFaction)
	assert.WithinDuration(t, now, dto.StartAt, time.Second)
	assert.WithinDuration(t, now, dto.CreatedAt, time.Second)
}
