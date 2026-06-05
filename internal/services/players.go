package services

import (
	"context"
	"log"
	"time"

	"github.com/ffx64/frontline-stats/internal/cache"
	"github.com/ffx64/frontline-stats/internal/dtos"
	"github.com/ffx64/frontline-stats/internal/entities"
	"github.com/ffx64/frontline-stats/internal/errors"
	"github.com/ffx64/frontline-stats/internal/helpers"
	"github.com/ffx64/frontline-stats/internal/repositories"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type PlayersService interface {
	Save(ctx context.Context, dto *dtos.PlayerSaveDTO) (*dtos.PlayerDTO, error)
	GetByGUID(ctx context.Context, guid string) (*dtos.PlayerDTO, error)
	GetPlayerStatsByGUID(ctx context.Context, guid string) (*dtos.PlayerStatsDTO, error)
	IfNotExistsCreate(ctx context.Context, username string, guid string, serverLastID string) (*dtos.PlayerDTO, error)
	Update(ctx context.Context, guid string, dto *dtos.PlayerUpdateDTO) (*dtos.PlayerDTO, error)
}

type playersService struct {
	repo      repositories.PlayersRepository
	statsRepo repositories.PlayersStatsRepository
	rdb       *redis.Client
}

func NewPlayersService(repo repositories.PlayersRepository, statsRepo repositories.PlayersStatsRepository, rdb *redis.Client) PlayersService {
	return &playersService{repo: repo, statsRepo: statsRepo, rdb: rdb}
}

func (s *playersService) GetPlayerStatsByGUID(ctx context.Context, guid string) (*dtos.PlayerStatsDTO, error) {
	key := cache.KeyPlayerStats(guid)
	if cached, err := cache.Get[dtos.PlayerStatsDTO](ctx, s.rdb, key); err != nil {
		log.Printf("[services:players] cache read error for key %s: %v", key, err)
	} else if cached != nil {
		return cached, nil
	}

	player, err := s.repo.FindByGUID(ctx, guid)
	if err != nil {
		log.Printf("[services:players] player not found %s: %v", guid, err)
		return nil, errors.ErrPlayerNotFound
	}
	log.Printf("[services:players] player found: %s", guid)

	stats, err := s.statsRepo.FindByPlayerID(ctx, player.ID.String())
	if err != nil {
		log.Printf("[services:stats] failed to find player stats %s: %v", player.GUID, err)
		return nil, errors.ErrStatsNotFound
	}
	log.Printf("[services:stats] player stats retrieved: %s", player.GUID)

	dto := helpers.ToPlayerStatsDTO(player, stats)
	cache.Set(ctx, s.rdb, key, dto, cache.TTLPlayerStats)
	return dto, nil
}

func (s *playersService) Save(ctx context.Context, dto *dtos.PlayerSaveDTO) (*dtos.PlayerDTO, error) {
	if user, _ := s.repo.FindByGUID(ctx, dto.GUID); user != nil {
		log.Printf("[services:players] player already exists: %s", dto.GUID)
		return nil, errors.ErrPlayerExists
	}

	lastServerId, err := uuid.Parse(dto.LastServerID)
	if err != nil {
		log.Printf("[services:players] failed to parse LastServerID %s: %v", dto.LastServerID, err)
		return nil, errors.ErrUUIDError
	}

	now := time.Now()
	player := entities.Players{
		GUID:            dto.GUID,
		Username:        dto.Username,
		Admin:           0,
		Premium:         0,
		PremiumStartAt:  nil,
		PremiumExpireAt: nil,
		IsActive:        true,
		IsBanned:        false,
		LastServerID:    &lastServerId,
		LastLogin:       &now,
		Platform:        dto.Platform,
		UpdatedAt:       now,
		CreatedAt:       now,
	}
	stats := entities.PlayersStats{Level: 1}

	if err := s.repo.CreateWithStats(ctx, &player, &stats); err != nil {
		log.Printf("[services:players] failed to save player to DB: %v", err)
		return nil, errors.New("failed to save player: "+err.Error(), 500)
	}
	log.Printf("[services:players] player and stats created: %s", player.GUID)

	return helpers.ToPlayerDTO(&player), nil
}

func (s *playersService) IfNotExistsCreate(ctx context.Context, username string, guid string, serverLastID string) (*dtos.PlayerDTO, error) {
	if user, _ := s.repo.FindByGUID(ctx, guid); user != nil {
		log.Printf("[services:players] player already exists: %s", guid)
		return helpers.ToPlayerDTO(user), nil
	}

	lastServerId, err := uuid.Parse(serverLastID)
	if err != nil {
		log.Printf("[services:players] failed to parse LastServerID %s: %v", serverLastID, err)
		return nil, errors.ErrUUIDError
	}

	now := time.Now()
	player := entities.Players{
		GUID:            guid,
		Username:        username,
		Admin:           0,
		Premium:         0,
		PremiumStartAt:  nil,
		PremiumExpireAt: nil,
		IsActive:        true,
		IsBanned:        false,
		LastServerID:    &lastServerId,
		LastLogin:       &now,
		Platform:        "unknown",
		UpdatedAt:       now,
		CreatedAt:       now,
	}
	stats := entities.PlayersStats{Level: 1}

	if err := s.repo.CreateWithStats(ctx, &player, &stats); err != nil {
		log.Printf("[services:players] failed to save player to DB: %v", err)
		return nil, errors.New("failed to save player: "+err.Error(), 500)
	}
	log.Printf("[services:players] player and stats created: %s", player.GUID)

	return helpers.ToPlayerDTO(&player), nil
}

func (s *playersService) GetByGUID(ctx context.Context, guid string) (*dtos.PlayerDTO, error) {
	key := cache.KeyPlayer(guid)
	if cached, err := cache.Get[dtos.PlayerDTO](ctx, s.rdb, key); err != nil {
		log.Printf("[services:players] cache read error for key %s: %v", key, err)
	} else if cached != nil {
		return cached, nil
	}

	player, err := s.repo.FindByGUID(ctx, guid)
	if err != nil {
		log.Printf("[services:players] player not found %s: %v", guid, err)
		return nil, errors.ErrPlayerNotFound
	}
	log.Printf("[services:players] player found: %s", guid)

	dto := helpers.ToPlayerDTO(player)
	cache.Set(ctx, s.rdb, key, dto, cache.TTLPlayer)
	return dto, nil
}

func (s *playersService) Update(ctx context.Context, guid string, dto *dtos.PlayerUpdateDTO) (*dtos.PlayerDTO, error) {
	player, err := s.repo.FindByGUID(ctx, guid)
	if err != nil || player == nil {
		log.Printf("[services:players] player not found to update %s: %v", guid, err)
		return nil, errors.ErrPlayerNotFound
	}

	if dto.Username != "" {
		player.Username = dto.Username
	}

	if player.IsActive != dto.IsActive {
		player.IsActive = dto.IsActive
	}

	if dto.LastServerID != "" {
		lastServerID, err := uuid.Parse(dto.LastServerID)
		if err != nil {
			log.Printf("[services:players] failed to parse LastServerID %s: %v", dto.LastServerID, err)
			return nil, errors.ErrUUIDError
		}
		player.LastServerID = &lastServerID
	}

	if err := s.repo.Update(ctx, player); err != nil {
		log.Printf("[services:players] failed to update player %s: %v", guid, err)
		return nil, errors.New("failed to update player: "+err.Error(), 500)
	}

	cache.Delete(ctx, s.rdb, cache.KeyPlayer(guid), cache.KeyPlayerStats(guid))
	log.Printf("[services:players] player updated: %s", guid)

	return helpers.ToPlayerDTO(player), nil
}
