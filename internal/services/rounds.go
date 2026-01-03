package services

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/ffx64/gamestats-backend/internal/dtos"
	"github.com/ffx64/gamestats-backend/internal/entities"
	"github.com/ffx64/gamestats-backend/internal/errors"
	"github.com/ffx64/gamestats-backend/internal/helpers"
	"github.com/ffx64/gamestats-backend/internal/repositories"
	"github.com/google/uuid"
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
}

func NewRoundsService(repo repositories.RoundsRepository, roundsStatsRepo repositories.RoundsStatsRepository) RoundsService {
	return &roundsService{repo: repo, statsRepo: roundsStatsRepo}
}

func (s *roundsService) SaveRound(ctx context.Context, dto *dtos.RoundsCreateDTO) (*dtos.RoundsDTO, error) {
	serverId, err := uuid.Parse(dto.ServerID)
	if err != nil {
		log.Printf("[services:rounds] erro ao parsear ServerID %s: %v", dto.ServerID, err)
		return nil, errors.ErrUUIDError
	}

	round := entities.Rounds{
		ServerID:      serverId,
		CurrentMode:   dto.CurrentMode,
		MissionHeader: dto.MissionHeader,
		Status:        "in_progress",
		WinnerFaction: "",
		EndedAt:       nil,
		StartAt:       time.Now(),
		CreatedAt:     time.Now(),
	}

	lastRound, err := s.repo.FindLastRoundInProgressByServerID(ctx, serverId)

	if err != nil {
		log.Printf("[services:rounds] erro ao buscar última rodada em progresso para o servidor %s: %v", dto.ServerID, err)
		return nil, errors.New("falha ao verificar rodadas em progresso: "+err.Error(), http.StatusInternalServerError)
	}

	if lastRound != nil {
		log.Printf("[services:rounds] atualizando o status da última rodada em progresso para 'crash': roundID=%s", lastRound.ID)
		lastRound.Status = "crash"
		s.repo.Update(ctx, lastRound, lastRound.ID)
	}

	saved, err := s.repo.Save(ctx, &round)
	if err != nil {
		log.Printf("[services:rounds] erro ao salvar rodada no DB: %v", err)
		return nil, errors.New("falha ao salvar rodadas: "+err.Error(), http.StatusInternalServerError)
	}

	log.Printf("[services:rounds] rodada salva com sucesso: serverID=%s startAt=%s", dto.ServerID, round.StartAt)
	return helpers.ToRoundsDTO(saved), nil
}

func (s *roundsService) GetRoundByID(ctx context.Context, id uuid.UUID) (*dtos.RoundsDTO, error) {
	if id == uuid.Nil {
		log.Println("[services:rounds] erro: RoundID é nulo")
		return nil, errors.ErrUUIDError
	}

	round, err := s.repo.FindByID(ctx, id)
	if err != nil {
		log.Printf("[services:rounds] erro ao buscar rodada no DB: %v", err)
		return nil, err
	}

	if round == nil {
		log.Printf("[services:rounds] rodada não encontrada: %s", id)
		return nil, errors.ErrRoundNotFound
	}

	log.Printf("[services:rounds] rodada encontrada: %s", id)
	return helpers.ToRoundsDTO(round), nil
}

func (s *roundsService) GetScoreboardByRoundID(ctx context.Context, roundId uuid.UUID) ([]entities.RoundsStats, error) {
	if roundId == uuid.Nil {
		log.Printf("[services:rounds_stats] UUID da rodada inválido: %v", roundId)
		return nil, errors.ErrUUIDError
	}

	scoreboard, err := s.statsRepo.FindScoreboardByRoundID(ctx, roundId)
	if err != nil {
		log.Printf("[services:rounds_stats] erro ao buscar estatísticas da rodada por round_id %v: %v", roundId, err)
		return nil, errors.New("falha ao buscar estatísticas das rodadas: "+err.Error(), 500)
	}

	log.Printf("[services:rounds_stats] estatísticas da rodada recuperadas com sucesso para round_id: %v", roundId)
	return scoreboard, nil
}

func (s *roundsService) GetAllRoundsByServerIDAndPlayerID(ctx context.Context, serverId, playerId uuid.UUID, limit, offset int) (dtos.RoundsDTOs, error) {
	if serverId == uuid.Nil {
		log.Printf("[services:rounds_stats] UUID do servidor inválido: %v", serverId)
		return dtos.RoundsDTOs{}, errors.ErrUUIDError
	}

	if playerId == uuid.Nil {
		log.Printf("[services:rounds_stats] UUID do jogador inválido: %v", playerId)
		return dtos.RoundsDTOs{}, errors.ErrUUIDError
	}

	stats, total, err := s.repo.FindAllRoundsByServerIDAndPlayerID(ctx, serverId, playerId, limit, offset)
	if err != nil {
		log.Printf("[services:rounds_stats] erro ao buscar as rodadas para server_id %v e player_id %v: %v", serverId, playerId, err)
		return dtos.RoundsDTOs{}, errors.New("falha ao buscar estatísticas das rodadas: "+err.Error(), 500)
	}

	log.Printf("[services:rounds_stats] rodadas recuperadas com sucesso para server_id: %v e player_id: %v", serverId, playerId)

	return helpers.ToRoundsDTOs(total, stats), nil
}

func (s *roundsService) UpdateRoundEnded(ctx context.Context, id uuid.UUID, dto *dtos.RoundsUpdatedEndedDTO) (*dtos.RoundsDTO, error) {
	if id == uuid.Nil {
		log.Println("[services:rounds] erro: RoundID é nulo")
		return nil, errors.ErrUUIDError
	}

	round, err := s.repo.FindByID(ctx, id)
	if err != nil {
		log.Printf("[services:rounds] erro ao buscar rodada para atualizar: %v", err)
		return nil, errors.New("falha ao atualizar rodada: "+err.Error(), http.StatusInternalServerError)
	}

	if round == nil {
		log.Printf("[services:rounds] rodada não encontrada para finalizar: %s", id)
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
		log.Printf("[services:rounds] erro ao finalizar rodada %s: %v", id, err)
		return nil, errors.New("falha ao finalizar partida: "+err.Error(), http.StatusInternalServerError)
	}

	log.Printf("[services:rounds] rodada finalizada com sucesso: %s winnerFaction=%s", id, round.WinnerFaction)
	return helpers.ToRoundsDTO(updated), nil
}
