package dtos_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/ffx64/frontline-stats/internal/dtos"
)

func TestPlayerDTO_JSONSerialization(t *testing.T) {
	now := time.Now()
	id := uuid.New().String()
	lastServerID := uuid.New().String()

	player := dtos.PlayerDTO{
		ID:              id,
		GUID:            "test-guid",
		Username:        "testuser",
		Admin:           1,
		Premium:         1,
		PremiumStartAt:  &now,
		PremiumExpireAt: &now,
		IsActive:        true,
		IsBanned:        false,
		LastServerID:    &lastServerID,
		LastLogin:       &now,
		Platform:        "pc",
		UpdatedAt:       now,
		CreatedAt:       now,
	}

	data, err := json.Marshal(player)
	assert.NoError(t, err)
	assert.Contains(t, string(data), `"guid":"test-guid"`)
	assert.Contains(t, string(data), `"username":"testuser"`)
}

func TestPlayerSaveDTO_Validation_Success(t *testing.T) {
	dto := dtos.PlayerSaveDTO{
		GUID:         "some-guid",
		Username:     "user",
		LastServerID: uuid.New().String(),
		Platform:     "pc",
	}

	validate := validator.New()
	err := validate.Struct(dto)
	assert.NoError(t, err)
}

func TestPlayerSaveDTO_Validation_Fails(t *testing.T) {
	validate := validator.New()

	dto := dtos.PlayerSaveDTO{}
	err := validate.Struct(dto)
	assert.Error(t, err)

	validationErrors := err.(validator.ValidationErrors)
	assert.Equal(t, 7, len(validationErrors))

	fields := []string{
		"GUID", "Username", "LastServerID", "Platform",
		"MachineProfileName", "MachineName", "MachineAdapterName",
	}

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
