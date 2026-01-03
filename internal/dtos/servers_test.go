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

func TestServersDTO_JSONSerialization(t *testing.T) {
	now := time.Now()
	id := uuid.New().String()

	server := dtos.ServersDTO{
		ID:        id,
		Name:      "SAS Server #1",
		CreatedAt: now,
	}

	data, err := json.Marshal(server)
	assert.NoError(t, err)
	assert.Contains(t, string(data), `"name":"SAS Server #1"`)
	assert.Contains(t, string(data), `"id":"`+id+`"`)
	assert.Contains(t, string(data), `"created_at"`)
}

func TestServersSaveDTO_Validation_Success(t *testing.T) {
	dto := dtos.ServersSaveDTO{
		Name: "Servidor Principal",
	}

	validate := validator.New()
	err := validate.Struct(dto)
	assert.NoError(t, err)
}

func TestServersSaveDTO_Validation_Fails(t *testing.T) {
	validate := validator.New()

	dto := dtos.ServersSaveDTO{}
	err := validate.Struct(dto)
	assert.Error(t, err)

	validationErrors := err.(validator.ValidationErrors)
	assert.Equal(t, 1, len(validationErrors))

	field := validationErrors[0].Field()
	assert.Equal(t, "Name", field)
}
