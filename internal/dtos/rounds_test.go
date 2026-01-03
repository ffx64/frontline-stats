package dtos_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/ffx64/gamestats-backend/internal/dtos"
)

func TestRoundsDTO_JSONSerialization(t *testing.T) {
	now := time.Now()
	id := uuid.New().String()
	serverID := uuid.New().String()

	round := dtos.RoundsDTO{
		ID:            id,
		ServerID:      serverID,
		CurrentMode:   "CTF",
		MissionHeader: "Island War",
		Status:        "active",
		WinnerFaction: "US",
		EndedAt:       &now,
		StartAt:       now,
		CreatedAt:     now,
	}

	data, err := json.Marshal(round)
	assert.NoError(t, err)
	assert.Contains(t, string(data), `"current_mode":"CTF"`)
	assert.Contains(t, string(data), `"mission_header":"Island War"`)
	assert.Contains(t, string(data), `"status":"active"`)
	assert.Contains(t, string(data), `"winner_faction":"US"`)
	assert.Contains(t, string(data), `"server_id":"`+serverID+`"`)
	assert.Contains(t, string(data), `"id":"`+id+`"`)
}

func TestRoundsCreateDTO_Validation_Success(t *testing.T) {
	dto := dtos.RoundsCreateDTO{
		ServerID:      uuid.New().String(),
		CurrentMode:   "TDM",
		MissionHeader: "Desert Strike",
		Status:        "ended",
		WinnerFaction: "CSAT",
	}

	validate := validator.New()
	err := validate.Struct(dto)
	assert.NoError(t, err)
}

func TestRoundsCreateDTO_Validation_Fails(t *testing.T) {
	validate := validator.New()

	dto := dtos.RoundsCreateDTO{}
	err := validate.Struct(dto)
	assert.Error(t, err)

	validationErrors := err.(validator.ValidationErrors)
	assert.Equal(t, 4, len(validationErrors))

	fields := []string{"ServerID", "CurrentMode", "MissionHeader", "Status"}
	for _, field := range fields {
		found := false
		for _, ve := range validationErrors {
			if ve.Field() == field {
				found = true
				break
			}
		}
		assert.True(t, found, "campo %s deveria falhar", field)
	}
}
