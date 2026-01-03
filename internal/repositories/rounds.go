package repositories

import (
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/ffx64/gamestats-backend/internal/entities"
)

type RoundsRepository interface {
	Save(ctx context.Context, round *entities.Rounds) (*entities.Rounds, error)
	FindByID(ctx context.Context, id uuid.UUID) (*entities.Rounds, error)
	FindAll(ctx context.Context, limit, offset int) ([]entities.Rounds, int64, error)
	FindLastRoundInProgressByServerID(ctx context.Context, serverID uuid.UUID) (*entities.Rounds, error)
	FindAllRoundsByServerIDAndPlayerID(ctx context.Context, serverId, playerId uuid.UUID, limit, offset int) ([]entities.Rounds, int64, error)
	FindAllRoundsInProgress(ctx context.Context) ([]entities.Rounds, error)
	Update(ctx context.Context, round *entities.Rounds, id uuid.UUID) (*entities.Rounds, error)
}

type roundsRepository struct {
	db *gorm.DB
}

func NewRoundsRepository(db *gorm.DB) RoundsRepository {
	return &roundsRepository{db: db}
}

func (r *roundsRepository) Save(ctx context.Context, round *entities.Rounds) (*entities.Rounds, error) {
	tx := r.db.WithContext(ctx).Create(round)
	if tx.Error != nil {
		log.Printf("[repository:rounds] erro ao criar rodada: %v", tx.Error)
		return nil, fmt.Errorf("erro ao criar rodada: %w", tx.Error)
	}
	log.Printf("[repository:rounds] rodada criada: %v", round.ID)
	return round, nil
}

func (r *roundsRepository) FindByID(ctx context.Context, id uuid.UUID) (*entities.Rounds, error) {
	var round entities.Rounds
	tx := r.db.WithContext(ctx).First(&round, "id = ?", id)
	if tx.Error != nil {
		if tx.Error == gorm.ErrRecordNotFound {
			log.Printf("[repository:rounds] nenhuma rodada encontrada com id: %v", id)
			return nil, nil
		}
		log.Printf("[repository:rounds] erro ao buscar rodada por id %v: %v", id, tx.Error)
		return nil, fmt.Errorf("erro ao buscar rodada por id: %w", tx.Error)
	}
	log.Printf("[repository:rounds] rodada recuperada: %v", round.ID)
	return &round, nil
}

func (r *roundsRepository) FindAll(ctx context.Context, limit, offset int) ([]entities.Rounds, int64, error) {
	var rounds []entities.Rounds
	var total int64

	if err := r.db.WithContext(ctx).Model(&entities.Rounds{}).Count(&total).Error; err != nil {
		log.Printf("[repository:rounds] erro ao contar rodadas: %v", err)
		return nil, 0, fmt.Errorf("erro ao contar rodadas: %w", err)
	}

	tx := r.db.WithContext(ctx).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&rounds)

	if tx.Error != nil {
		log.Printf("[repository:rounds] erro ao listar rodadas: %v", tx.Error)
		return nil, 0, fmt.Errorf("erro ao listar rodadas: %w", tx.Error)
	}

	log.Printf("[repository:rounds] %d rodadas recuperadas (offset=%d, limit=%d, total=%d)",
		len(rounds), offset, limit, total)

	return rounds, total, nil
}

func (r *roundsRepository) FindLastRoundInProgressByServerID(ctx context.Context, serverID uuid.UUID) (*entities.Rounds, error) {
	var round entities.Rounds
	tx := r.db.WithContext(ctx).
		Where("server_id = ? AND status = ?", serverID, "in_progress").
		Order("created_at DESC").
		First(&round)

	if tx.Error != nil {
		if tx.Error == gorm.ErrRecordNotFound {
			log.Printf("[repository:rounds] nenhum round em progresso encontrado para o servidor: %v", serverID)
			return nil, nil
		}
		log.Printf("[repository:rounds] erro ao buscar último round em progresso para o servidor %v: %v", serverID, tx.Error)
		return nil, fmt.Errorf("erro ao buscar último round em progresso: %w", tx.Error)
	}

	log.Printf("[repository:rounds] round em progresso recuperado: %v", round.ID)
	return &round, nil
}

func (r *roundsRepository) FindAllRoundsInProgress(ctx context.Context) ([]entities.Rounds, error) {
	var rounds []entities.Rounds
	tx := r.db.WithContext(ctx).Where("status = in_progress").Find(&rounds)
	if tx.Error != nil {
		log.Printf("[repository:rounds] erro ao buscar rounds em progresso: %v", tx.Error)
		return nil, fmt.Errorf("erro ao buscar rounds em progresso: %w", tx.Error)
	}

	log.Printf("[repository:rounds] rounds em progresso recuperados: %d registros", len(rounds))
	return rounds, nil
}

func (r *roundsRepository) FindAllRoundsByServerIDAndPlayerID(ctx context.Context, serverId, playerId uuid.UUID, limit, offset int) ([]entities.Rounds, int64, error) {
	var stats []entities.RoundsStats
	var total int64

	// RESOLVER DAQUI A POUCO
	if err := r.db.WithContext(ctx).Model(&entities.RoundsStats{}).Where("server_id = ? AND player_id = ?", serverId, playerId).Count(&total).Error; err != nil {
		log.Printf("[repository:rounds_stats] erro ao contar stats do player: %v", err)
		return nil, 0, fmt.Errorf("erro ao contar stats: %w", err)
	}

	if err := r.db.WithContext(ctx).Model(&entities.RoundsStats{}).Where("server_id = ? AND player_id = ?", serverId, playerId).Limit(limit).Offset(offset).Find(&stats).Error; err != nil {
		log.Printf("[repository:rounds_stats] erro ao buscar stats do player: %v", err)
		return nil, 0, fmt.Errorf("erro ao buscar stats: %w", err)
	}

	if len(stats) == 0 {
		return []entities.Rounds{}, total, nil
	}

	roundIDs := make([]uuid.UUID, 0, len(stats))
	for _, s := range stats {
		roundIDs = append(roundIDs, s.RoundID)
	}

	var rounds []entities.Rounds
	if err := r.db.WithContext(ctx).Where("id IN ?", roundIDs).Find(&rounds).Error; err != nil {
		log.Printf("[repository:rounds] erro ao buscar rounds dos stats: %v", err)
		return nil, 0, fmt.Errorf("erro ao buscar rounds: %w", err)
	}

	log.Printf("[repository:rounds] rounds recuperados: %d registros", len(rounds))

	return rounds, total, nil
}

func (r *roundsRepository) Update(ctx context.Context, round *entities.Rounds, id uuid.UUID) (*entities.Rounds, error) {
	tx := r.db.WithContext(ctx).Model(&entities.Rounds{}).
		Where("id = ?", id).
		Updates(map[string]any{
			"status":         round.Status,
			"winner_faction": round.WinnerFaction,
			"ended_at":       round.EndedAt,
		})

	if tx.Error != nil {
		log.Printf("[repository:rounds] erro ao atualizar round %v: %v", id, tx.Error)
		return nil, fmt.Errorf("erro ao atualizar round: %w", tx.Error)
	}

	if tx.RowsAffected == 0 {
		log.Printf("[repository:rounds] nenhum round encontrado com id: %v", id)
		return nil, fmt.Errorf("nenhum round encontrado com o id %s", id)
	}

	updated, err := r.FindByID(ctx, id)
	if err != nil {
		log.Printf("[repository:rounds] erro ao buscar round atualizado %v: %v", id, err)
		return nil, fmt.Errorf("erro ao buscar round atualizado: %w", err)
	}

	log.Printf("[repository:rounds] round atualizado: %v", updated.ID)
	return updated, nil
}
