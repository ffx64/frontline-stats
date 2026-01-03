package services

import (
	"context"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"

	"github.com/ffx64/gamestats-backend/internal/dtos"
	"github.com/ffx64/gamestats-backend/internal/entities"
	"github.com/ffx64/gamestats-backend/internal/errors"
	"github.com/ffx64/gamestats-backend/internal/repositories"
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
		log.Println("[services:kills] aviso: nenhuma kill recebida para salvar")
		return errors.ErrNoKillsReceived
	}

	kills := make([]*entities.Kills, 0, count)
	now := time.Now()

	for _, dto := range killsDTO {
		serverId, err := uuid.Parse(dto.ServerID)
		if err != nil {
			log.Printf("[services:kills] erro ao parsear ServerID %s: %v", dto.ServerID, err)
			return errors.ErrInvalidServerID
		}

		roundId, err := uuid.Parse(dto.RoundID)
		if err != nil {
			log.Printf("[services:kills] erro ao parsear RoundID %s: %v", dto.RoundID, err)
			return errors.ErrInvalidRoundID
		}

		if _, err := uuid.Parse(dto.KillerID); err != nil {
			log.Printf("[services:kills] erro ao parsear KillerID %s: %v", dto.KillerID, err)
			return errors.ErrInvalidKillerID
		}

		if _, err := uuid.Parse(dto.VictimID); err != nil {
			log.Printf("[services:kills] erro ao parsear VictimID %s: %v", dto.VictimID, err)
			return errors.ErrInvalidVictimID
		}

		server, err := s.serverRepo.FindByID(ctx, serverId)
		if err != nil {
			log.Printf("[services:kills] erro ao buscar servidor no DB: %v", err)
			return errors.ErrServerNotFoundDB
		}
		if server == nil {
			log.Printf("[services:kills] servidor não encontrado: %s", serverId)
			return errors.ErrServerNotFound
		}

		round, err := s.roundRepo.FindByID(ctx, roundId)
		if err != nil {
			log.Printf("[services:kills] erro ao buscar rodada no DB: %v", err)
			return errors.ErrRoundNotFoundDB
		}
		if round == nil {
			log.Printf("[services:kills] rodada não encontrada: %s", roundId)
			return errors.ErrRoundNotFound
		}

		killer, err := s.playerRepo.FindByGUID(ctx, dto.KillerID)
		if err != nil {
			log.Printf("[services:kills] falha ao buscar killer no DB: %v", err)
			return errors.ErrPlayerLookupFail
		}
		if killer == nil {
			log.Printf("[services:kills] killer não encontrado: %s", dto.KillerID)
			return errors.ErrPlayerNotFound
		}

		victim, err := s.playerRepo.FindByGUID(ctx, dto.VictimID)
		if err != nil {
			log.Printf("[services:kills] falha ao buscar vítima no DB: %v", err)
			return errors.ErrPlayerLookupFail
		}
		if victim == nil {
			log.Printf("[services:kills] vítima não encontrada: %s", dto.VictimID)
			return errors.ErrPlayerNotFound
		}

		if dto.Timestamp == nil {
			dto.Timestamp = &now
		}

		kill := &entities.Kills{
			ID:               uuid.New(),
			ServerID:         serverId,
			RoundID:          roundId,
			KillerID:         killer.ID,
			VictimID:         victim.ID,
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

	if len(kills) == 0 {
		log.Println("[services:kills] aviso: todas as kills descartadas por erros de validação")
		return errors.ErrNoKillsReceived
	}

	size := min(100, len(kills))

	if err := s.repo.SaveBatch(ctx, kills, size); err != nil {
		log.Printf("[services:kills] erro ao salvar batch de %d kills: %v", count, err)
		return errors.ErrBatchSaveFailed
	}

	log.Printf("[services:kills] sucesso: %d kills salvas em batch de %d", count, size)
	return nil
}
