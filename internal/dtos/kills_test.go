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

func TestKillsDTO_JSONSerialization(t *testing.T) {
	now := time.Now()
	id := uuid.New().String()

	dto := dtos.KillsDTO{
		ID:               id,
		ServerID:         uuid.New().String(),
		RoundID:          uuid.New().String(),
		KillerID:         uuid.New().String(),
		VictimID:         uuid.New().String(),
		VictimWeaponName: "AK-47",
		VictimWeaponType: "Rifle",
		KillerWeaponName: "M4A1",
		KillerWeaponType: "Rifle",
		HitZone:          "Head",
		Distance:         120.5,
		IsHeadshot:       true,
		IsFriendly:       false,
		IsVehicle:        false,
		KillerTeam:       "Blue",
		VictimTeam:       "Red",
		Timestamp:        now,
		CreatedAt:        now,
	}

	data, err := json.Marshal(dto)
	assert.NoError(t, err)
	assert.Contains(t, string(data), `"victim_weapon_name":"AK-47"`)
	assert.Contains(t, string(data), `"killer_weapon_name":"M4A1"`)
	assert.Contains(t, string(data), `"is_headshot":true`)
}

func TestKillsSaveDTO_Validation_Success(t *testing.T) {
	now := time.Now()

	dto := dtos.KillsSaveDTO{
		ServerID:         uuid.New().String(),
		RoundID:          uuid.New().String(),
		KillerID:         uuid.New().String(),
		VictimID:         uuid.New().String(),
		VictimWeaponName: "AK-47",
		VictimWeaponType: "Rifle",
		KillerWeaponName: "M4A1",
		KillerWeaponType: "Rifle",
		HitZone:          "Head",
		Distance:         150.3,
		IsHeadshot:       true,
		IsFriendly:       false,
		IsVehicle:        false,
		KillerTeam:       "Blue",
		VictimTeam:       "Red",
		Timestamp:        &now,
	}

	validate := validator.New()
	err := validate.Struct(dto)
	assert.NoError(t, err)
}

func TestKillsSaveDTO_Validation_Fails(t *testing.T) {
	validate := validator.New()

	dto := dtos.KillsSaveDTO{}
	err := validate.Struct(dto)
	assert.Error(t, err)

	validationErrors := err.(validator.ValidationErrors)

	// Ajusta esse número se tu for adicionar tags `validate:"required"` nos campos do DTO
	expectedFields := []string{
		"ServerID", "RoundID", "KillerID", "VictimID",
		"VictimWeaponName", "VictimWeaponType",
		"KillerWeaponName", "KillerWeaponType",
		"HitZone", "KillerTeam", "VictimTeam", "Timestamp",
	}

	for _, field := range expectedFields {
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
