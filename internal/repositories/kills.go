package repositories

import (
	"context"
	"log"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/ffx64/gamestats-backend/internal/entities"
)

type KillsRepository interface {
	Save(ctx context.Context, kill *entities.Kills) error
	SaveBatch(ctx context.Context, kills []*entities.Kills, size int) error
	GetKillsByPlayerID(ctx context.Context, playerID uuid.UUID) ([]entities.Kills, error)
	GetDeathsByPlayerID(ctx context.Context, playerID uuid.UUID) ([]entities.Kills, error)
	GetKillsForPlayerByServerID(ctx context.Context, playerID, serverID uuid.UUID) ([]entities.Kills, error)
	GetDeathsForPlayerByServerID(ctx context.Context, playerID, serverID uuid.UUID) ([]entities.Kills, error)
	GetKillsForPlayerByRoundID(ctx context.Context, playerID, roundID uuid.UUID) ([]entities.Kills, error)
	GetDeathsForPlayerByRoundID(ctx context.Context, playerID, roundID uuid.UUID) ([]entities.Kills, error)
	GetTop10KillsAndDeathByPlayerID(ctx context.Context, playerID uuid.UUID) ([]entities.Kills, error)
}

type killsRepository struct {
	db *gorm.DB
}

func NewKillsRepository(db *gorm.DB) KillsRepository {
	return &killsRepository{db: db}
}

func (r *killsRepository) Save(ctx context.Context, kill *entities.Kills) error {
	if err := r.db.WithContext(ctx).Create(kill).Error; err != nil {
		log.Printf("[repository:kills] failed to save kill: %v", err)
		return err
	}
	log.Printf("[repository:kills] kill saved successfully: %v", kill.ID)
	return nil
}

func (r *killsRepository) GetKillsByPlayerID(ctx context.Context, playerID uuid.UUID) ([]entities.Kills, error) {
	var kills []entities.Kills
	err := r.db.WithContext(ctx).
		Where("killer_id = ?", playerID).
		Find(&kills).Error

	if err != nil {
		log.Printf("[repository:kills] failed to get kills for player %s: %v", playerID, err)
		return nil, err
	}
	log.Printf("[repository:kills] retrieved %d kills for player %s", len(kills), playerID)
	return kills, nil
}

func (r *killsRepository) GetDeathsByPlayerID(ctx context.Context, playerID uuid.UUID) ([]entities.Kills, error) {
	var deaths []entities.Kills
	err := r.db.WithContext(ctx).
		Where("victim_id = ?", playerID).
		Find(&deaths).Error

	if err != nil {
		log.Printf("[repository:kills] failed to get deaths for player %s: %v", playerID, err)
		return nil, err
	}
	log.Printf("[repository:kills] retrieved %d deaths for player %s", len(deaths), playerID)
	return deaths, nil
}

func (r *killsRepository) GetKillsForPlayerByServerID(ctx context.Context, playerID, serverID uuid.UUID) ([]entities.Kills, error) {
	var kills []entities.Kills
	err := r.db.WithContext(ctx).
		Where("(killer_id = ?) AND server_id = ?", playerID, serverID).
		Find(&kills).Error

	if err != nil {
		log.Printf("[repository:kills] failed to get kills by server for player %s: %v", playerID, err)
		return nil, err
	}
	log.Printf("[repository:kills] retrieved %d kills for player %s on server %s", len(kills), playerID, serverID)
	return kills, nil
}

func (r *killsRepository) GetDeathsForPlayerByServerID(ctx context.Context, playerID, serverID uuid.UUID) ([]entities.Kills, error) {
	var kills []entities.Kills
	err := r.db.WithContext(ctx).
		Where("(victim_id = ?) AND server_id = ?", playerID, serverID).
		Find(&kills).Error

	if err != nil {
		log.Printf("[repository:kills] failed to get deaths by server for player %s: %v", playerID, err)
		return nil, err
	}
	log.Printf("[repository:kills] retrieved %d deaths for player %s on server %s", len(kills), playerID, serverID)
	return kills, nil
}

func (r *killsRepository) GetKillsForPlayerByRoundID(ctx context.Context, playerID, roundID uuid.UUID) ([]entities.Kills, error) {
	var kills []entities.Kills
	err := r.db.WithContext(ctx).
		Where("(killer_id = ?) AND round_id = ?", playerID, roundID).
		Find(&kills).Error

	if err != nil {
		log.Printf("[repository:kills] failed to get kills by round for player %s: %v", playerID, err)
		return nil, err
	}
	log.Printf("[repository:kills] retrieved %d kills for player %s in round %s", len(kills), playerID, roundID)
	return kills, nil
}

func (r *killsRepository) GetDeathsForPlayerByRoundID(ctx context.Context, playerID, roundID uuid.UUID) ([]entities.Kills, error) {
	var kills []entities.Kills
	err := r.db.WithContext(ctx).
		Where("(victim_id = ?) AND round_id = ?", playerID, roundID).
		Find(&kills).Error

	if err != nil {
		log.Printf("[repository:kills] failed to get deaths by round for player %s: %v", playerID, err)
		return nil, err
	}
	log.Printf("[repository:kills] retrieved %d deaths for player %s in round %s", len(kills), playerID, roundID)
	return kills, nil
}

func (r *killsRepository) GetTop10KillsAndDeathByPlayerID(ctx context.Context, playerID uuid.UUID) ([]entities.Kills, error) {
	var kills []entities.Kills
	err := r.db.WithContext(ctx).
		Where("killer_id = ?", playerID).
		Order("timestamp DESC").
		Limit(10).
		Find(&kills).Error

	if err != nil {
		log.Printf("[repository:kills] failed to get top 10 kills/deaths for player %s: %v", playerID, err)
		return nil, err
	}
	log.Printf("[repository:kills] retrieved %d top 10 kills/deaths for player %s", len(kills), playerID)
	return kills, nil
}

func (r *killsRepository) SaveBatch(ctx context.Context, kills []*entities.Kills, size int) error {
	if len(kills) == 0 {
		log.Println("[repository:kills] warn: no kills to save in batch")
		return nil
	}

	if size <= 0 {
		size = len(kills)
	}

	if err := r.db.WithContext(ctx).CreateInBatches(kills, size).Error; err != nil {
		log.Printf("[repository:kills] failed to save kill batch: %v", err)
		return err
	}

	log.Printf("[repository:kills] %d kills saved in batch of %d", len(kills), size)
	return nil
}
