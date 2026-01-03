package services

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/ffx64/gamestats-backend/internal/dtos"
	"github.com/ffx64/gamestats-backend/internal/entities"
	"github.com/ffx64/gamestats-backend/internal/errors"
	"github.com/ffx64/gamestats-backend/internal/helpers"
	"github.com/ffx64/gamestats-backend/internal/repositories"
	"github.com/google/uuid"
)

type PlayersService interface {
	// Save Player
	Save(ctx context.Context, dto *dtos.PlayerSaveDTO) (*dtos.PlayerDTO, error)

	// Get Player by GUID
	GetByGUID(ctx context.Context, guid string) (*dtos.PlayerDTO, error)

	// Get Player Stats by GUID
	GetPlayerStatsByGUID(ctx context.Context, guid string) (*dtos.PlayerStatsDTO, error)

	// If Not Exists Create Player
	IfNotExistsCreate(ctx context.Context, username string, guid string, serverLastID string) (*dtos.PlayerDTO, error)

	// Update Player by GUID
	Update(ctx context.Context, guid string, dto *dtos.PlayerUpdateDTO) (*dtos.PlayerDTO, error)
}

type playersService struct {
	repo      repositories.PlayersRepository
	statsRepo repositories.PlayersStatsRepository
}

func NewPlayersService(repo repositories.PlayersRepository, statsRepo repositories.PlayersStatsRepository) PlayersService {
	return &playersService{repo: repo, statsRepo: statsRepo}
}

func (s *playersService) GetPlayerStatsByGUID(ctx context.Context, guid string) (*dtos.PlayerStatsDTO, error) {
	player, err := s.repo.FindByGUID(ctx, guid)
	if err != nil {
		log.Printf("[services:players] jogador não encontrado %s: %v", guid, err)
		return nil, errors.ErrPlayerNotFound
	}
	log.Printf("[services:players] jogador encontrado: %s", guid)

	stats, err := s.statsRepo.FindByPlayerID(ctx, player.ID.String())
	if err != nil {
		log.Printf("[services:stats] erro ao buscar estatísticas do jogador %s: %v", player.GUID, err)
		return nil, errors.ErrPlayerNotFound
	}
	log.Printf("[services:stats] estatísticas do jogador recuperadas: %s", player.GUID)

	stats.PlayerID, err = uuid.Parse(player.GUID)

	if err != nil {
		log.Printf("[services:stats] erro ao parsear PlayerID %s: %v", player.GUID, err)
		return nil, fmt.Errorf("falha ao parsear PlayerID: %w", err)
	}

	return helpers.ToPlayerStatsDTO(player, stats), nil
}

func (s *playersService) Save(ctx context.Context, dto *dtos.PlayerSaveDTO) (*dtos.PlayerDTO, error) {
	if user, _ := s.repo.FindByGUID(ctx, dto.GUID); user != nil {
		log.Printf("[services:players] jogador já existe: %s", dto.GUID)
		return nil, errors.ErrPlayerExists
	}

	lastServerId, err := uuid.Parse(dto.LastServerID)
	if err != nil {
		log.Printf("[services:players] erro ao parsear LastServerID %s: %v", dto.LastServerID, err)
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

	if err := s.repo.Save(ctx, &player); err != nil {
		log.Printf("[services:players] erro ao salvar jogador no DB: %v", err)
		return nil, errors.New("falha ao salvar jogador: "+err.Error(), 500)
	}
	log.Printf("[services:players] jogador salvo: %s", player.GUID)

	stats := entities.PlayersStats{
		PlayerID:           player.ID,
		Level:              1,
		XP:                 0,
		Kills:              0,
		Deaths:             0,
		GrenadesThrown:     0,
		FriendlyFireKills:  0,
		FriendlyFireDeaths: 0,
		VehicleKills:       0,
		VehicleDeaths:      0,
		RatioKDR:           0.0,
		RatioHeadshot:      0.0,
		RatioVehicle:       0.0,
		WeaponsMostUsed:    "",
		HitZonesMostDied:   "",
		MaxKillDistance:    0.0,
		UpdatedAt:          now,
		CreatedAt:          now,
	}

	if err := s.statsRepo.Save(ctx, &stats); err != nil {
		log.Printf("[services:players] erro ao criar stats do jogador %s: %v", player.GUID, err)
		return nil, errors.New("falha ao salvar stats do jogador: "+err.Error(), 500)
	}
	log.Printf("[services:players] stats criadas para jogador: %s", player.GUID)

	return helpers.ToPlayerDTO(&player), nil
}

func (s *playersService) IfNotExistsCreate(ctx context.Context, username string, guid string, serverLastID string) (*dtos.PlayerDTO, error) {
	if user, _ := s.repo.FindByGUID(ctx, guid); user != nil {
		log.Printf("[services:players] jogador já existe: %s", guid)
		return helpers.ToPlayerDTO(user), nil
	}

	lastServerId, err := uuid.Parse(serverLastID)
	if err != nil {
		log.Printf("[services:players] erro ao parsear LastServerID %s: %v", lastServerId, err)
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
		Platform:        "unkdown",
		UpdatedAt:       now,
		CreatedAt:       now,
	}

	if err := s.repo.Save(ctx, &player); err != nil {
		log.Printf("[services:players] erro ao salvar jogador no DB: %v", err)
		return nil, errors.New("falha ao salvar jogador: "+err.Error(), 500)
	}
	log.Printf("[services:players] jogador salvo: %s", player.GUID)

	stats := entities.PlayersStats{
		PlayerID:           player.ID,
		Level:              1,
		XP:                 0,
		Kills:              0,
		Deaths:             0,
		GrenadesThrown:     0,
		FriendlyFireKills:  0,
		FriendlyFireDeaths: 0,
		VehicleKills:       0,
		VehicleDeaths:      0,
		RatioKDR:           0.0,
		RatioHeadshot:      0.0,
		RatioVehicle:       0.0,
		WeaponsMostUsed:    "",
		HitZonesMostDied:   "",
		MaxKillDistance:    0.0,
		UpdatedAt:          now,
		CreatedAt:          now,
	}

	if err := s.statsRepo.Save(ctx, &stats); err != nil {
		log.Printf("[services:players] erro ao criar stats do jogador %s: %v", player.GUID, err)
		return nil, errors.New("falha ao salvar stats do jogador: "+err.Error(), 500)
	}
	log.Printf("[services:players] stats criadas para jogador: %s", player.GUID)

	return helpers.ToPlayerDTO(&player), nil
}

func (s *playersService) GetByGUID(ctx context.Context, guid string) (*dtos.PlayerDTO, error) {
	player, err := s.repo.FindByGUID(ctx, guid)
	if err != nil {
		log.Printf("[services:players] jogador não encontrado %s: %v", guid, err)
		return nil, errors.ErrPlayerNotFound
	}
	log.Printf("[services:players] jogador encontrado: %s", guid)
	return helpers.ToPlayerDTO(player), nil
}

func (s *playersService) Update(ctx context.Context, guid string, dto *dtos.PlayerUpdateDTO) (*dtos.PlayerDTO, error) {
	player, err := s.repo.FindByGUID(ctx, guid)
	if err != nil || player == nil {
		log.Printf("[services:players] jogador não encontrado para atualizar %s: %v", guid, err)
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
			log.Printf("[services:players] erro ao parsear LastServerID %s: %v", dto.LastServerID, err)
			return nil, errors.ErrUUIDError
		}
		player.LastServerID = &lastServerID
	}

	if err := s.repo.Update(ctx, player); err != nil {
		log.Printf("[services:players] erro ao atualizar jogador %s: %v", guid, err)
		return nil, errors.New("falha ao atualizar jogador: "+err.Error(), 500)
	}
	log.Printf("[services:players] jogador atualizado: %s", guid)

	return helpers.ToPlayerDTO(player), nil
}
