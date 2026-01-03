package helpers

import (
	"github.com/ffx64/gamestats-backend/internal/dtos"
	"github.com/ffx64/gamestats-backend/internal/entities"
)

func ToRoundsDTO(server *entities.Rounds) *dtos.RoundsDTO {
	return &dtos.RoundsDTO{
		ID:            server.ID.String(),
		ServerID:      server.ServerID.String(),
		CurrentMode:   server.CurrentMode,
		MissionHeader: server.MissionHeader,
		Status:        server.Status,
		WinnerFaction: server.WinnerFaction,
		EndedAt:       server.EndedAt,
		StartAt:       server.StartAt,
		CreatedAt:     server.CreatedAt,
	}
}

func ToRoundsDTOs(total int64, rounds []entities.Rounds) dtos.RoundsDTOs {
	roundsDTOs := make([]dtos.RoundsDTO, len(rounds))
	for i, round := range rounds {
		roundsDTOs[i] = *ToRoundsDTO(&round)
	}

	return dtos.RoundsDTOs{Total: total, Data: roundsDTOs}
}
