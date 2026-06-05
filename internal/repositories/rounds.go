package repositories

import (
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/ffx64/frontline-stats/internal/entities"
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
		log.Printf("[repository:rounds] failed to create round: %v", tx.Error)
		return nil, fmt.Errorf("failed to create round: %w", tx.Error)
	}
	log.Printf("[repository:rounds] round created: %v", round.ID)
	return round, nil
}

func (r *roundsRepository) FindByID(ctx context.Context, id uuid.UUID) (*entities.Rounds, error) {
	var round entities.Rounds
	tx := r.db.WithContext(ctx).First(&round, "id = ?", id)
	if tx.Error != nil {
		if tx.Error == gorm.ErrRecordNotFound {
			log.Printf("[repository:rounds] no round found with id: %v", id)
			return nil, nil
		}
		log.Printf("[repository:rounds] failed to find round by id %v: %v", id, tx.Error)
		return nil, fmt.Errorf("failed to find round by id: %w", tx.Error)
	}
	log.Printf("[repository:rounds] round retrieved: %v", round.ID)
	return &round, nil
}

func (r *roundsRepository) FindAll(ctx context.Context, limit, offset int) ([]entities.Rounds, int64, error) {
	var rounds []entities.Rounds
	var total int64

	if err := r.db.WithContext(ctx).Model(&entities.Rounds{}).Count(&total).Error; err != nil {
		log.Printf("[repository:rounds] failed to count rounds: %v", err)
		return nil, 0, fmt.Errorf("failed to count rounds: %w", err)
	}

	tx := r.db.WithContext(ctx).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&rounds)

	if tx.Error != nil {
		log.Printf("[repository:rounds] failed to list rounds: %v", tx.Error)
		return nil, 0, fmt.Errorf("failed to list rounds: %w", tx.Error)
	}

	log.Printf("[repository:rounds] %d rounds retrieved (offset=%d, limit=%d, total=%d)",
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
			log.Printf("[repository:rounds] no in-progress round found for server: %v", serverID)
			return nil, nil
		}
		log.Printf("[repository:rounds] failed to find last in-progress round for server %v: %v", serverID, tx.Error)
		return nil, fmt.Errorf("failed to find last in-progress round: %w", tx.Error)
	}

	log.Printf("[repository:rounds] in-progress round retrieved: %v", round.ID)
	return &round, nil
}

func (r *roundsRepository) FindAllRoundsInProgress(ctx context.Context) ([]entities.Rounds, error) {
	var rounds []entities.Rounds
	tx := r.db.WithContext(ctx).Where("status = ?", "in_progress").Find(&rounds)
	if tx.Error != nil {
		log.Printf("[repository:rounds] failed to find rounds in progress: %v", tx.Error)
		return nil, fmt.Errorf("failed to find rounds in progress: %w", tx.Error)
	}

	log.Printf("[repository:rounds] rounds in progress retrieved: %d records", len(rounds))
	return rounds, nil
}

func (r *roundsRepository) FindAllRoundsByServerIDAndPlayerID(ctx context.Context, serverId, playerId uuid.UUID, limit, offset int) ([]entities.Rounds, int64, error) {
	var total int64
	var rounds []entities.Rounds

	base := r.db.WithContext(ctx).
		Model(&entities.Rounds{}).
		Joins("INNER JOIN rounds_stats ON rounds_stats.round_id = rounds.id").
		Where("rounds_stats.server_id = ? AND rounds_stats.player_id = ?", serverId, playerId)

	if err := base.Count(&total).Error; err != nil {
		log.Printf("[repository:rounds_stats] failed to count rounds for player: %v", err)
		return nil, 0, fmt.Errorf("failed to count rounds: %w", err)
	}

	if err := base.Order("rounds.created_at DESC").Limit(limit).Offset(offset).Find(&rounds).Error; err != nil {
		log.Printf("[repository:rounds_stats] failed to find rounds for player: %v", err)
		return nil, 0, fmt.Errorf("failed to fetch rounds: %w", err)
	}

	log.Printf("[repository:rounds] rounds retrieved: %d records", len(rounds))
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
		log.Printf("[repository:rounds] failed to update round %v: %v", id, tx.Error)
		return nil, fmt.Errorf("failed to update round: %w", tx.Error)
	}

	if tx.RowsAffected == 0 {
		log.Printf("[repository:rounds] no round found with id: %v", id)
		return nil, fmt.Errorf("no round found with id %s", id)
	}

	updated, err := r.FindByID(ctx, id)
	if err != nil {
		log.Printf("[repository:rounds] failed to find updated round %v: %v", id, err)
		return nil, fmt.Errorf("failed to find updated round: %w", err)
	}

	log.Printf("[repository:rounds] round updated: %v", updated.ID)
	return updated, nil
}
