package services

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/ffx64/gamestats-backend/internal/cache"
	"github.com/ffx64/gamestats-backend/internal/dtos"
	"github.com/ffx64/gamestats-backend/internal/entities"
	"github.com/ffx64/gamestats-backend/internal/errors"
	"github.com/ffx64/gamestats-backend/internal/helpers"
	"github.com/ffx64/gamestats-backend/internal/repositories"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type RoundsService interface {
	// Create Round
	SaveRound(ctx context.Context, round *dtos.RoundsCreateDTO) (*dtos.RoundsDTO, error)

	// Get Round by ID
	GetRoundByID(ctx context.Context, id uuid.UUID) (*dtos.RoundsDTO, error)

	// Get Scoreboard by RoundID
	GetScoreboardByRoundID(ctx context.Context, roundId uuid.UUID) ([]entities.RoundsStats, error)

	// Get Rounds Stats by ServerID and PlayerID
	GetAllRoundsByServerIDAndPlayerID(ctx context.Context, serverId, playerId uuid.UUID, limit, offset int) (dtos.RoundsDTOs, error)

	// Update Round to Ended
	UpdateRoundEnded(ctx context.Context, id uuid.UUID, round *dtos.RoundsUpdatedEndedDTO) (*dtos.RoundsDTO, error)
}

type roundsService struct {
	repo      repositories.RoundsRepository
	statsRepo repositories.RoundsStatsRepository
	rdb       *redis.Client
}

func NewRoundsService(repo repositories.RoundsRepository, roundsStatsRepo repositories.RoundsStatsRepository, rdb *redis.Client) RoundsService {
	return &roundsService{repo: repo, statsRepo: roundsStatsRepo, rdb: rdb}
}

func (s *roundsService) SaveRound(ctx context.Context, dto *dtos.RoundsCreateDTO) (*dtos.RoundsDTO, error) {
	serverId, err := uuid.Parse(dto.ServerID)
	if err != nil {
		log.Printf("[services:rounds] failed to parse ServerID %s: %v", dto.ServerID, err)
		return nil, errors.ErrUUIDError
	}

	now := time.Now()
	round := entities.Rounds{
		ServerID:      serverId,
		CurrentMode:   dto.CurrentMode,
		MissionHeader: dto.MissionHeader,
		Status:        "in_progress",
		WinnerFaction: "",
		EndedAt:       nil,
		StartAt:       now,
		CreatedAt:     now,
	}

	lastRound, err := s.repo.FindLastRoundInProgressByServerID(ctx, serverId)
	if err != nil {
		log.Printf("[services:rounds] failed to find last in-progress round for server %s: %v", dto.ServerID, err)
		return nil, errors.New("failed to check in-progress rounds: "+err.Error(), http.StatusInternalServerError)
	}

	if lastRound != nil {
		log.Printf("[services:rounds] updating previous in-progress round to 'crash': roundID=%s", lastRound.ID)
		lastRound.Status = "crash"
		if _, err := s.repo.Update(ctx, lastRound, lastRound.ID); err != nil {
			log.Printf("[services:rounds] failed to mark previous round as crash: %v", err)
		}
	}

	saved, err := s.repo.Save(ctx, &round)
	if err != nil {
		log.Printf("[services:rounds] failed to save round to DB: %v", err)
		return nil, errors.New("failed to save round: "+err.Error(), http.StatusInternalServerError)
	}

	log.Printf("[services:rounds] round saved successfully: serverID=%s startAt=%s", dto.ServerID, round.StartAt)
	return helpers.ToRoundsDTO(saved), nil
}

func (s *roundsService) GetRoundByID(ctx context.Context, id uuid.UUID) (*dtos.RoundsDTO, error) {
	if id == uuid.Nil {
		log.Println("[services:rounds] error: RoundID is nil")
		return nil, errors.ErrUUIDError
	}

	key := cache.KeyRound(id.String())
	if cached, err := cache.Get[dtos.RoundsDTO](ctx, s.rdb, key); err != nil {
		log.Printf("[services:rounds] cache read error for key %s: %v", key, err)
	} else if cached != nil {
		return cached, nil
	}

	round, err := s.repo.FindByID(ctx, id)
	if err != nil {
		log.Printf("[services:rounds] failed to find round in DB: %v", err)
		return nil, err
	}

	if round == nil {
		log.Printf("[services:rounds] round not found: %s", id)
		return nil, errors.ErrRoundNotFound
	}

	dto := helpers.ToRoundsDTO(round)
	cache.Set(ctx, s.rdb, key, dto, cache.TTLRound)
	log.Printf("[services:rounds] round found: %s", id)
	return dto, nil
}

func (s *roundsService) GetScoreboardByRoundID(ctx context.Context, roundId uuid.UUID) ([]entities.RoundsStats, error) {
	if roundId == uuid.Nil {
		log.Printf("[services:rounds_stats] invalid round UUID: %v", roundId)
		return nil, errors.ErrUUIDError
	}

	key := cache.KeyScoreboard(roundId.String())
	if cached, err := cache.Get[[]entities.RoundsStats](ctx, s.rdb, key); err != nil {
		log.Printf("[services:rounds_stats] cache read error for key %s: %v", key, err)
	} else if cached != nil {
		return *cached, nil
	}

	scoreboard, err := s.statsRepo.FindScoreboardByRoundID(ctx, roundId)
	if err != nil {
		log.Printf("[services:rounds_stats] failed to find round stats for round_id %v: %v", roundId, err)
		return nil, errors.New("failed to fetch round stats: "+err.Error(), 500)
	}

	cache.Set(ctx, s.rdb, key, scoreboard, cache.TTLScoreboard)
	log.Printf("[services:rounds_stats] round stats retrieved for round_id: %v", roundId)
	return scoreboard, nil
}

func (s *roundsService) GetAllRoundsByServerIDAndPlayerID(ctx context.Context, serverId, playerId uuid.UUID, limit, offset int) (dtos.RoundsDTOs, error) {
	if serverId == uuid.Nil {
		log.Printf("[services:rounds_stats] invalid server UUID: %v", serverId)
		return dtos.RoundsDTOs{}, errors.ErrUUIDError
	}

	if playerId == uuid.Nil {
		log.Printf("[services:rounds_stats] invalid player UUID: %v", playerId)
		return dtos.RoundsDTOs{}, errors.ErrUUIDError
	}

	stats, total, err := s.repo.FindAllRoundsByServerIDAndPlayerID(ctx, serverId, playerId, limit, offset)
	if err != nil {
		log.Printf("[services:rounds_stats] failed to find rounds for server_id %v and player_id %v: %v", serverId, playerId, err)
		return dtos.RoundsDTOs{}, errors.New("failed to fetch round stats: "+err.Error(), 500)
	}

	log.Printf("[services:rounds_stats] rounds retrieved for server_id: %v and player_id: %v", serverId, playerId)
	return helpers.ToRoundsDTOs(total, stats), nil
}

func (s *roundsService) UpdateRoundEnded(ctx context.Context, id uuid.UUID, dto *dtos.RoundsUpdatedEndedDTO) (*dtos.RoundsDTO, error) {
	if id == uuid.Nil {
		log.Println("[services:rounds] error: RoundID is nil")
		return nil, errors.ErrUUIDError
	}

	round, err := s.repo.FindByID(ctx, id)
	if err != nil {
		log.Printf("[services:rounds] failed to find round to update: %v", err)
		return nil, errors.New("failed to update round: "+err.Error(), http.StatusInternalServerError)
	}

	if round == nil {
		log.Printf("[services:rounds] round not found to end: %s", id)
		return nil, errors.ErrRoundNotFound
	}

	if dto.WinnerFaction != "" {
		round.WinnerFaction = dto.WinnerFaction
	}

	now := time.Now()
	round.Status = "ended"
	round.EndedAt = &now

	updated, err := s.repo.Update(ctx, round, id)
	if err != nil {
		log.Printf("[services:rounds] failed to end round %s: %v", id, err)
		return nil, errors.New("failed to end round: "+err.Error(), http.StatusInternalServerError)
	}

	cache.Delete(ctx, s.rdb, cache.KeyRound(id.String()), cache.KeyScoreboard(id.String()))
	log.Printf("[services:rounds] round ended successfully: %s winnerFaction=%s", id, round.WinnerFaction)
	return helpers.ToRoundsDTO(updated), nil
}
