package services

import (
	"context"

	"github.com/ffx64/gamestats-backend/internal/dtos"
	"github.com/ffx64/gamestats-backend/internal/repositories"
	"github.com/google/uuid"
)

type RoundsStatsService interface {
	// Get all rounds by server ID and player ID
	GetAllRoundsByServerIDAndPlayerID(ctx context.Context, serverId, playerId uuid.UUID, limit, offset int) (dtos.RoundsDTOs, error)
}

type roundsStatsService struct {
	repo repositories.RoundsStatsRepository
}

func NewRoundsStatsService(repo repositories.RoundsStatsRepository) RoundsStatsService {
	return &roundsStatsService{repo: repo}
}

func (s *roundsStatsService) GetAllRoundsByServerIDAndPlayerID(ctx context.Context, serverId, playerId uuid.UUID, limit, offset int) (dtos.RoundsDTOs, error) {
	return dtos.RoundsDTOs{}, nil
}
