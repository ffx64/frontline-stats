package helpers

import (
	"github.com/ffx64/gamestats-backend/internal/dtos"
	"github.com/ffx64/gamestats-backend/internal/entities"
)

func ToServersDTO(server *entities.Servers) *dtos.ServersDTO {
	return &dtos.ServersDTO{
		ID:        server.ID.String(),
		Name:      server.Name,
		CreatedAt: server.CreatedAt,
	}
}
