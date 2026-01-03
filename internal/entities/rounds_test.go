package entities_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/ffx64/gamestats-backend/internal/entities"
)

func TestRoundsInitialization(t *testing.T) {
	now := time.Now()
	serverID := uuid.New()

	round := entities.Rounds{
		ID:            uuid.New(),
		ServerID:      serverID,
		CurrentMode:   "CTI",
		MissionHeader: "Operation Storm",
		Status:        "active",
		WinnerFaction: "NATO",
		StartAt:       now,
		CreatedAt:     now,
	}

	assert.NotNil(t, round.ID)
	assert.Equal(t, serverID, round.ServerID)
	assert.Equal(t, "CTI", round.CurrentMode)
	assert.Equal(t, "Operation Storm", round.MissionHeader)
	assert.Equal(t, "active", round.Status)
	assert.Equal(t, "NATO", round.WinnerFaction)
	assert.WithinDuration(t, now, round.CreatedAt, time.Second)
}

func TestRoundsJSONMarshalling(t *testing.T) {
	now := time.Now()
	serverID := uuid.New()

	round := entities.Rounds{
		ID:            uuid.New(),
		ServerID:      serverID,
		CurrentMode:   "AAS",
		MissionHeader: "Night Strike",
		Status:        "ended",
		WinnerFaction: "CSAT",
		StartAt:       now,
		CreatedAt:     now,
	}

	data, err := json.Marshal(round)
	assert.NoError(t, err)
	assert.Contains(t, string(data), `"mission_header":"Night Strike"`)

	var unmarshalled entities.Rounds
	err = json.Unmarshal(data, &unmarshalled)
	assert.NoError(t, err)
	assert.Equal(t, round.MissionHeader, unmarshalled.MissionHeader)
	assert.Equal(t, round.Status, unmarshalled.Status)
	assert.WithinDuration(t, round.CreatedAt, unmarshalled.CreatedAt, time.Second)
}

func TestRoundsDefaultValues(t *testing.T) {
	round := entities.Rounds{
		ServerID: uuid.New(),
	}

	assert.Equal(t, uuid.Nil, round.ID, "o ID deve ser definido apenas no BeforeCreate, ainda nil aqui")
	assert.True(t, round.CreatedAt.IsZero(), "CreatedAt deve ser zero até o insert ser feito")
	assert.True(t, round.StartAt.IsZero(), "StartAt deve ser zero até o insert ser feito")
	assert.Nil(t, round.EndedAt, "EndedAt deve ser nil por padrão")
}

func TestRoundsBeforeCreateHook(t *testing.T) {
	round := entities.Rounds{
		ServerID: uuid.New(),
		Status:   "waiting",
	}

	err := round.BeforeCreate(nil)
	assert.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, round.ID, "o ID deve ser gerado no BeforeCreate")
}
