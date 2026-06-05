package repositories

import (
	"context"
	"fmt"
	"log"

	"github.com/ffx64/gamestats-backend/internal/entities"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ServersRepository interface {
	Save(ctx context.Context, server *entities.Servers) (*entities.Servers, error)
	FindByID(ctx context.Context, id uuid.UUID) (*entities.Servers, error)
	FindByName(ctx context.Context, name string) (*entities.Servers, error)
	FindAll(ctx context.Context) ([]entities.Servers, error)
	Update(ctx context.Context, server *entities.Servers) (*entities.Servers, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type serverRepository struct {
	db *gorm.DB
}

func NewServersRepository(db *gorm.DB) ServersRepository {
	return &serverRepository{db: db}
}

func (r *serverRepository) Save(ctx context.Context, server *entities.Servers) (*entities.Servers, error) {
	tx := r.db.WithContext(ctx).Create(server)
	if tx.Error != nil {
		log.Printf("[repository:servers] failed to create server: %v", tx.Error)
		return nil, fmt.Errorf("failed to create server: %w", tx.Error)
	}
	log.Printf("[repository:servers] server created: %v", server.ID)
	return server, nil
}

func (r *serverRepository) FindByID(ctx context.Context, id uuid.UUID) (*entities.Servers, error) {
	var server entities.Servers
	tx := r.db.WithContext(ctx).First(&server, "id = ?", id)
	if tx.Error != nil {
		if tx.Error == gorm.ErrRecordNotFound {
			log.Printf("[repository:servers] no server found with id: %v", id)
			return nil, nil
		}
		log.Printf("[repository:servers] failed to find server by id %v: %v", id, tx.Error)
		return nil, fmt.Errorf("failed to find server by id: %w", tx.Error)
	}
	log.Printf("[repository:servers] server retrieved: %v", server.ID)
	return &server, nil
}

func (r *serverRepository) FindByName(ctx context.Context, name string) (*entities.Servers, error) {
	var server entities.Servers
	tx := r.db.WithContext(ctx).First(&server, "name = ?", name)
	if tx.Error != nil {
		if tx.Error == gorm.ErrRecordNotFound {
			log.Printf("[repository:servers] no server found with name: %v", name)
			return nil, nil
		}
		log.Printf("[repository:servers] failed to find server by name %v: %v", name, tx.Error)
		return nil, fmt.Errorf("failed to find server by name: %w", tx.Error)
	}
	log.Printf("[repository:servers] server retrieved: %v", server.ID)
	return &server, nil
}

func (r *serverRepository) FindAll(ctx context.Context) ([]entities.Servers, error) {
	var servers []entities.Servers
	tx := r.db.WithContext(ctx).Order("created_at DESC").Find(&servers)
	if tx.Error != nil {
		log.Printf("[repository:servers] failed to list servers: %v", tx.Error)
		return nil, fmt.Errorf("failed to list servers: %w", tx.Error)
	}
	log.Printf("[repository:servers] %d servers retrieved", len(servers))
	return servers, nil
}

func (r *serverRepository) Update(ctx context.Context, server *entities.Servers) (*entities.Servers, error) {
	tx := r.db.WithContext(ctx).Model(&entities.Servers{}).
		Where("id = ?", server.ID).
		Updates(map[string]any{
			"name": server.Name,
		})

	if tx.Error != nil {
		log.Printf("[repository:servers] failed to update server %v: %v", server.ID, tx.Error)
		return nil, fmt.Errorf("failed to update server: %w", tx.Error)
	}

	if tx.RowsAffected == 0 {
		log.Printf("[repository:servers] no server found with id: %v", server.ID)
		return nil, fmt.Errorf("no server found with id %s", server.ID)
	}

	updated, err := r.FindByID(ctx, server.ID)
	if err != nil {
		log.Printf("[repository:servers] failed to find updated server %v: %v", server.ID, err)
		return nil, err
	}

	log.Printf("[repository:servers] server updated: %v", updated.ID)
	return updated, nil
}

func (r *serverRepository) Delete(ctx context.Context, id uuid.UUID) error {
	tx := r.db.WithContext(ctx).Delete(&entities.Servers{}, "id = ?", id)
	if tx.Error != nil {
		log.Printf("[repository:servers] failed to delete server %v: %v", id, tx.Error)
		return fmt.Errorf("failed to delete server: %w", tx.Error)
	}
	if tx.RowsAffected == 0 {
		log.Printf("[repository:servers] no server found to delete with id: %v", id)
		return fmt.Errorf("no server found with id %s", id)
	}
	log.Printf("[repository:servers] server deleted: %v", id)
	return nil
}
