package repositories

import (
	"context"
	"fmt"
	"log"

	"github.com/ffx64/frontline-stats/internal/entities"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PlayersRepository interface {
	Save(ctx context.Context, player *entities.Players) error
	CreateWithStats(ctx context.Context, player *entities.Players, stats *entities.PlayersStats) error
	FindByID(ctx context.Context, id uuid.UUID) (*entities.Players, error)
	FindByGUID(ctx context.Context, guid string) (*entities.Players, error)
	FindAll(ctx context.Context) ([]entities.Players, error)
	Update(ctx context.Context, player *entities.Players) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type playerRepository struct {
	db *gorm.DB
}

func NewPlayersRepository(db *gorm.DB) PlayersRepository {
	return &playerRepository{db: db}
}

func (r *playerRepository) Save(ctx context.Context, player *entities.Players) error {
	if err := r.db.WithContext(ctx).Create(player).Error; err != nil {
		log.Printf("[repository:players] failed to create player: %v", err)
		return fmt.Errorf("failed to create player: %w", err)
	}
	log.Printf("[repository:players] player created: %v", player.ID)
	return nil
}

func (r *playerRepository) CreateWithStats(ctx context.Context, player *entities.Players, stats *entities.PlayersStats) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(player).Error; err != nil {
			log.Printf("[repository:players] failed to create player in transaction: %v", err)
			return fmt.Errorf("failed to create player: %w", err)
		}
		stats.PlayerID = player.ID
		if err := tx.Create(stats).Error; err != nil {
			log.Printf("[repository:players] failed to create player stats in transaction: %v", err)
			return fmt.Errorf("failed to create player stats: %w", err)
		}
		log.Printf("[repository:players] player and stats created in transaction: %v", player.ID)
		return nil
	})
}

func (r *playerRepository) FindByID(ctx context.Context, id uuid.UUID) (*entities.Players, error) {
	var player entities.Players
	if err := r.db.WithContext(ctx).First(&player, "id = ?", id).Error; err != nil {
		log.Printf("[repository:players] failed to find player by ID %v: %v", id, err)
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("player not found")
		}
		return nil, fmt.Errorf("failed to find player: %w", err)
	}
	log.Printf("[repository:players] player retrieved by ID: %v", player.ID)
	return &player, nil
}

func (r *playerRepository) FindByGUID(ctx context.Context, guid string) (*entities.Players, error) {
	var player entities.Players
	if err := r.db.WithContext(ctx).First(&player, "guid = ?", guid).Error; err != nil {
		log.Printf("[repository:players] failed to find player by GUID %v: %v", guid, err)
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("player not found")
		}
		return nil, fmt.Errorf("failed to find player: %w", err)
	}
	log.Printf("[repository:players] player retrieved by GUID: %v", player.GUID)
	return &player, nil
}

func (r *playerRepository) FindAll(ctx context.Context) ([]entities.Players, error) {
	var players []entities.Players
	if err := r.db.WithContext(ctx).Order("created_at DESC").Find(&players).Error; err != nil {
		log.Printf("[repository:players] failed to list players: %v", err)
		return nil, fmt.Errorf("failed to list players: %w", err)
	}
	log.Printf("[repository:players] %d players retrieved", len(players))
	return players, nil
}

func (r *playerRepository) Update(ctx context.Context, player *entities.Players) error {
	if err := r.db.WithContext(ctx).Save(player).Error; err != nil {
		log.Printf("[repository:players] failed to update player %v: %v", player.ID, err)
		return fmt.Errorf("failed to update player: %w", err)
	}
	log.Printf("[repository:players] player updated: %v", player.ID)
	return nil
}

func (r *playerRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if err := r.db.WithContext(ctx).Delete(&entities.Players{}, id).Error; err != nil {
		log.Printf("[repository:players] failed to delete player %v: %v", id, err)
		return fmt.Errorf("failed to delete player: %w", err)
	}
	log.Printf("[repository:players] player deleted: %v", id)
	return nil
}
