package services

import (
	"context"
	"fmt"
	"log"

	"github.com/ffx64/gamestats-backend/internal/cache"
	"github.com/ffx64/gamestats-backend/internal/entities"
	"github.com/ffx64/gamestats-backend/internal/repositories"
	"github.com/redis/go-redis/v9"
)

type PlayersStatsService interface {
	RefreshPlayerStats(ctx context.Context) error
	GetLeaderboard(ctx context.Context) ([]*entities.PlayersStats, error)
	GetLeaderboardHeadshots(ctx context.Context) ([]*entities.PlayersStats, error)
	GetLeaderboardVehicle(ctx context.Context) ([]*entities.PlayersStats, error)
}

type playersStatsService struct {
	repo repositories.PlayersStatsRepository
	rdb  *redis.Client
}

func NewPlayersStatsService(repo repositories.PlayersStatsRepository, rdb *redis.Client) PlayersStatsService {
	return &playersStatsService{repo: repo, rdb: rdb}
}

func (s *playersStatsService) RefreshPlayerStats(ctx context.Context) error {
	if err := s.repo.UpdatePlayerStats(ctx); err != nil {
		log.Printf("[services:stats] failed to update player stats: %v", err)
		return fmt.Errorf("failed to update player stats: %w", err)
	}

	log.Println("[services:stats] player stats updated successfully")
	return nil
}

func (s *playersStatsService) GetLeaderboard(ctx context.Context) ([]*entities.PlayersStats, error) {
	if cached, err := cache.Get[[]*entities.PlayersStats](ctx, s.rdb, cache.KeyLeaderboardKills); err != nil {
		log.Printf("[services:stats] cache read error for leaderboard:kills: %v", err)
	} else if cached != nil {
		return *cached, nil
	}

	stats, err := s.repo.GetLeaderboard(ctx)
	if err != nil {
		log.Printf("[services:stats] failed to get leaderboard: %v", err)
		return nil, err
	}
	cache.Set(ctx, s.rdb, cache.KeyLeaderboardKills, stats, cache.TTLLeaderboard)
	return stats, nil
}

func (s *playersStatsService) GetLeaderboardHeadshots(ctx context.Context) ([]*entities.PlayersStats, error) {
	if cached, err := cache.Get[[]*entities.PlayersStats](ctx, s.rdb, cache.KeyLeaderboardHeadshots); err != nil {
		log.Printf("[services:stats] cache read error for leaderboard:headshots: %v", err)
	} else if cached != nil {
		return *cached, nil
	}

	stats, err := s.repo.GetLeaderboardHeadshots(ctx)
	if err != nil {
		log.Printf("[services:stats] failed to get headshots leaderboard: %v", err)
		return nil, err
	}
	cache.Set(ctx, s.rdb, cache.KeyLeaderboardHeadshots, stats, cache.TTLLeaderboard)
	return stats, nil
}

func (s *playersStatsService) GetLeaderboardVehicle(ctx context.Context) ([]*entities.PlayersStats, error) {
	if cached, err := cache.Get[[]*entities.PlayersStats](ctx, s.rdb, cache.KeyLeaderboardVehicles); err != nil {
		log.Printf("[services:stats] cache read error for leaderboard:vehicles: %v", err)
	} else if cached != nil {
		return *cached, nil
	}

	stats, err := s.repo.GetLeaderboardVehicle(ctx)
	if err != nil {
		log.Printf("[services:stats] failed to get vehicles leaderboard: %v", err)
		return nil, err
	}
	cache.Set(ctx, s.rdb, cache.KeyLeaderboardVehicles, stats, cache.TTLLeaderboard)
	return stats, nil
}
