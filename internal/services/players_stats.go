package services

import (
	"context"
	"fmt"
	"log"

	"github.com/ffx64/gamestats-backend/internal/entities"
	"github.com/ffx64/gamestats-backend/internal/repositories"
)

type PlayersStatsService interface {
	// Refresh Player Stats
	RefreshPlayerStats(ctx context.Context) error
	GetLeaderboard(ctx context.Context) ([]*entities.PlayersStats, error)
	GetLeaderboardHeadshots(ctx context.Context) ([]*entities.PlayersStats, error)
	GetLeaderboardVehicle(ctx context.Context) ([]*entities.PlayersStats, error)
}

type playersStatsService struct {
	repo repositories.PlayersStatsRepository
}

func NewPlayersStatsService(repo repositories.PlayersStatsRepository) PlayersStatsService {
	return &playersStatsService{repo: repo}
}

func (s *playersStatsService) RefreshPlayerStats(ctx context.Context) error {
	if err := s.repo.UpdatePlayerStats(ctx); err != nil {
		log.Printf("[services:stats] erro ao atualizar estatísticas dos jogadores: %v", err)
		return fmt.Errorf("falha ao atualizar estatísticas dos jogadores: %w", err)
	}

	log.Println("[services:stats] estatísticas dos jogadores atualizadas com sucesso")
	return nil
}

func (s *playersStatsService) GetLeaderboard(ctx context.Context) ([]*entities.PlayersStats, error) {
	var stats []*entities.PlayersStats

	stats, err := s.repo.GetLeaderboard(ctx)

	if err != nil {
		log.Printf("[services:stats] erro ao obter leaderboard: %v", err)
		return nil, err
	}

	return stats, nil
}

func (s *playersStatsService) GetLeaderboardHeadshots(ctx context.Context) ([]*entities.PlayersStats, error) {
	var stats []*entities.PlayersStats

	stats, err := s.repo.GetLeaderboardHeadshots(ctx)

	if err != nil {
		log.Printf("[services:stats] erro ao obter leaderboard: %v", err)
		return nil, err
	}

	return stats, nil
}

func (s *playersStatsService) GetLeaderboardVehicle(ctx context.Context) ([]*entities.PlayersStats, error) {
	var stats []*entities.PlayersStats

	stats, err := s.repo.GetLeaderboardVehicle(ctx)

	if err != nil {
		log.Printf("[services:stats] erro ao obter leaderboard: %v", err)
		return nil, err
	}

	return stats, nil
}
