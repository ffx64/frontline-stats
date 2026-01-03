package helpers_test

import (
	"testing"
	"time"

	"github.com/ffx64/gamestats-backend/internal/entities"
	"github.com/ffx64/gamestats-backend/internal/helpers"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestToServersDTO_WithoutLastServerID(t *testing.T) {
	now := time.Now()

	server := &entities.Servers{
		ID:        uuid.New(),
		Name:      "Server test",
		CreatedAt: now,
	}

	dto := helpers.ToServersDTO(server)

	assert.Equal(t, server.ID.String(), dto.ID)
	assert.Equal(t, server.Name, dto.Name)
	assert.Equal(t, server.CreatedAt, dto.CreatedAt)
}
