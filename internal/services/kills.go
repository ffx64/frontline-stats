package services

import (
	"context"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"

	"github.com/ffx64/frontline-stats/internal/dtos"
	"github.com/ffx64/frontline-stats/internal/entities"
	"github.com/ffx64/frontline-stats/internal/errors"
	"github.com/ffx64/frontline-stats/internal/repositories"
)

type KillsService interface {
	SaveKills(ctx context.Context, killsDTO []dtos.KillsSaveDTO) error
}

type killsService struct {
	repo       repositories.KillsRepository
	roundRepo  repositories.RoundsRepository
	serverRepo repositories.ServersRepository
	playerRepo repositories.PlayersRepository
	rdb        *redis.Client
}

func NewKillsService(
	repo repositories.KillsRepository,
	roundRepo repositories.RoundsRepository,
	serverRepo repositories.ServersRepository,
	playerRepo repositories.PlayersRepository,
	rdb *redis.Client,
) KillsService {
	return &killsService{
		repo:       repo,
		roundRepo:  roundRepo,
		serverRepo: serverRepo,
		playerRepo: playerRepo,
		rdb:        rdb,
	}
}

func (s *killsService) SaveKills(ctx context.Context, killsDTO []dtos.KillsSaveDTO) error {
	count := len(killsDTO)
	if count == 0 {
		log.Println("[services:kills] no kills received to save")
		return errors.ErrNoKillsReceived
	}

	serverCache := make(map[uuid.UUID]*entities.Servers)
	roundCache := make(map[uuid.UUID]*entities.Rounds)
	playerCache := make(map[string]*entities.Players)

	kills := make([]*entities.Kills, 0, count)
	now := time.Now()

	for _, dto := range killsDTO {
		serverId, err := uuid.Parse(dto.ServerID)
		if err != nil {
			log.Printf("[services:kills] failed to parse ServerID %s: %v", dto.ServerID, err)
			return errors.ErrInvalidServerID
		}

		roundId, err := uuid.Parse(dto.RoundID)
		if err != nil {
			log.Printf("[services:kills] failed to parse RoundID %s: %v", dto.RoundID, err)
			return errors.ErrInvalidRoundID
		}

		if _, err := uuid.Parse(dto.KillerID); err != nil {
			log.Printf("[services:kills] failed to parse KillerID %s: %v", dto.KillerID, err)
			return errors.ErrInvalidKillerID
		}

		if _, err := uuid.Parse(dto.VictimID); err != nil {
			log.Printf("[services:kills] failed to parse VictimID %s: %v", dto.VictimID, err)
			return errors.ErrInvalidVictimID
		}

		if _, ok := serverCache[serverId]; !ok {
			server, err := s.serverRepo.FindByID(ctx, serverId)
			if err != nil {
				log.Printf("[services:kills] failed to find server in DB: %v", err)
				return errors.ErrServerNotFoundDB
			}
			if server == nil {
				log.Printf("[services:kills] server not found: %s", serverId)
				return errors.ErrServerNotFound
			}
			serverCache[serverId] = server
		}

		if _, ok := roundCache[roundId]; !ok {
			round, err := s.roundRepo.FindByID(ctx, roundId)
			if err != nil {
				log.Printf("[services:kills] failed to find round in DB: %v", err)
				return errors.ErrRoundNotFoundDB
			}
			if round == nil {
				log.Printf("[services:kills] round not found: %s", roundId)
				return errors.ErrRoundNotFound
			}
			roundCache[roundId] = round
		}

		if _, ok := playerCache[dto.KillerID]; !ok {
			killer, err := s.playerRepo.FindByGUID(ctx, dto.KillerID)
			if err != nil {
				log.Printf("[services:kills] failed to find killer in DB: %v", err)
				return errors.ErrPlayerLookupFail
			}
			if killer == nil {
				log.Printf("[services:kills] killer not found: %s", dto.KillerID)
				return errors.ErrPlayerNotFound
			}
			playerCache[dto.KillerID] = killer
		}

		if _, ok := playerCache[dto.VictimID]; !ok {
			victim, err := s.playerRepo.FindByGUID(ctx, dto.VictimID)
			if err != nil {
				log.Printf("[services:kills] failed to find victim in DB: %v", err)
				return errors.ErrPlayerLookupFail
			}
			if victim == nil {
				log.Printf("[services:kills] victim not found: %s", dto.VictimID)
				return errors.ErrPlayerNotFound
			}
			playerCache[dto.VictimID] = victim
		}

		if dto.Timestamp == nil {
			dto.Timestamp = &now
		}

		kill := &entities.Kills{
			ID:               uuid.New(),
			ServerID:         serverId,
			RoundID:          roundId,
			KillerID:         playerCache[dto.KillerID].ID,
			VictimID:         playerCache[dto.VictimID].ID,
			VictimWeaponName: dto.VictimWeaponName,
			VictimWeaponType: dto.VictimWeaponType,
			KillerWeaponName: dto.KillerWeaponName,
			KillerWeaponType: dto.KillerWeaponType,
			HitZone:          dto.HitZone,
			Distance:         dto.Distance,
			IsHeadshot:       dto.IsHeadshot,
			IsFriendly:       dto.IsFriendly,
			IsVehicle:        dto.IsVehicle,
			KillerTeam:       dto.KillerTeam,
			VictimTeam:       dto.VictimTeam,
			Timestamp:        *dto.Timestamp,
			CreatedAt:        now,
		}
		kills = append(kills, kill)
	}

	size := min(100, len(kills))

	if err := s.repo.SaveBatch(ctx, kills, size); err != nil {
		log.Printf("[services:kills] failed to save batch of %d kills: %v", count, err)
		return errors.ErrBatchSaveFailed
	}

	log.Printf("[services:kills] %d kills saved in batch of %d", count, size)
	return nil
}
