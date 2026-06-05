package entities_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/ffx64/frontline-stats/internal/entities"
)

func TestServerInitialization(t *testing.T) {
	now := time.Now()

	server := entities.Servers{
		ID:        uuid.New(),
		Name:      "SAS Main Server",
		CreatedAt: now,
	}

	assert.NotNil(t, server.ID)
	assert.Equal(t, "SAS Main Server", server.Name)
	assert.WithinDuration(t, now, server.CreatedAt, time.Second)
}

func TestServerJSONMarshalling(t *testing.T) {
	now := time.Now()
	server := entities.Servers{
		ID:        uuid.New(),
		Name:      "Server JSON",
		CreatedAt: now,
	}

	data, err := json.Marshal(server)
	assert.NoError(t, err)
	assert.Contains(t, string(data), `"name":"Server JSON"`)

	var unmarshalled entities.Servers
	err = json.Unmarshal(data, &unmarshalled)
	assert.NoError(t, err)
	assert.Equal(t, server.Name, unmarshalled.Name)
	assert.WithinDuration(t, server.CreatedAt, unmarshalled.CreatedAt, time.Second)
}

func TestServerDefaultValues(t *testing.T) {
	server := entities.Servers{
		Name: "Default Server",
	}

	assert.Equal(t, "Default Server", server.Name)
	assert.Equal(t, uuid.Nil, server.ID, "o ID deve ser definido apenas no BeforeCreate, ainda nil aqui")
	assert.True(t, server.CreatedAt.IsZero(), "CreatedAt deve ser zero até o insert ser feito")
}

func TestServerBeforeCreateHook(t *testing.T) {
	server := entities.Servers{Name: "Hook Test"}

	err := server.BeforeCreate(nil)
	assert.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, server.ID, "o ID deve ser gerado no BeforeCreate")
}
