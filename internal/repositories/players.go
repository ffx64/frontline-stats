package repositories

import (
	"context"
	"fmt"
	"log"

	"github.com/ffx64/gamestats-backend/internal/entities"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PlayersRepository interface {
	Save(ctx context.Context, player *entities.Players) error
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
		log.Printf("[repository:players] erro ao criar jogador: %v", err)
		return fmt.Errorf("erro ao criar jogador: %w", err)
	}
	log.Printf("[repository:players] jogador criado: %v", player.ID)
	return nil
}

func (r *playerRepository) FindByID(ctx context.Context, id uuid.UUID) (*entities.Players, error) {
	var player entities.Players
	if err := r.db.WithContext(ctx).First(&player, "id = ?", id).Error; err != nil {
		log.Printf("[repository:players] erro ao buscar jogador por ID %v: %v", id, err)
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("jogador não encontrado")
		}
		return nil, fmt.Errorf("erro ao buscar jogador: %w", err)
	}
	log.Printf("[repository:players] jogador recuperado por ID: %v", player.ID)
	return &player, nil
}

func (r *playerRepository) FindByGUID(ctx context.Context, guid string) (*entities.Players, error) {
	var player entities.Players
	if err := r.db.WithContext(ctx).First(&player, "guid = ?", guid).Error; err != nil {
		log.Printf("[repository:players] erro ao buscar jogador por GUID %v: %v", guid, err)
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("jogador não encontrado")
		}
		return nil, fmt.Errorf("erro ao buscar jogador: %w", err)
	}
	log.Printf("[repository:players] jogador recuperado por GUID: %v", player.GUID)
	return &player, nil
}

func (r *playerRepository) FindAll(ctx context.Context) ([]entities.Players, error) {
	var players []entities.Players
	if err := r.db.WithContext(ctx).Order("created_at DESC").Find(&players).Error; err != nil {
		log.Printf("[repository:players] erro ao listar jogadores: %v", err)
		return nil, fmt.Errorf("erro ao listar jogadores: %w", err)
	}
	log.Printf("[repository:players] %d jogadores recuperados", len(players))
	return players, nil
}

func (r *playerRepository) Update(ctx context.Context, player *entities.Players) error {
	if err := r.db.WithContext(ctx).Save(player).Error; err != nil {
		log.Printf("[repository:players] erro ao atualizar jogador %v: %v", player.ID, err)
		return fmt.Errorf("erro ao atualizar jogador: %w", err)
	}
	log.Printf("[repository:players] jogador atualizado: %v", player.ID)
	return nil
}

func (r *playerRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if err := r.db.WithContext(ctx).Delete(&entities.Players{}, id).Error; err != nil {
		log.Printf("[repository:players] erro ao deletar jogador %v: %v", id, err)
		return fmt.Errorf("erro ao deletar jogador: %w", err)
	}
	log.Printf("[repository:players] jogador deletado: %v", id)
	return nil
}
