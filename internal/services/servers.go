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

type ServersService interface {
	SaveServer(ctx context.Context, dto *dtos.ServersSaveDTO) (*dtos.ServersDTO, error)
	GetServerByID(ctx context.Context, id uuid.UUID) (*dtos.ServersDTO, error)
	GetAllServers(ctx context.Context) ([]dtos.ServersDTO, error)
	UpdateServer(ctx context.Context, id uuid.UUID, server *dtos.ServersSaveDTO) (*dtos.ServersDTO, error)
	DeleteServer(ctx context.Context, id uuid.UUID) (bool, error)
}

type serversService struct {
	repo repositories.ServersRepository
	rdb  *redis.Client
}

func NewServersService(repo repositories.ServersRepository, rdb *redis.Client) ServersService {
	return &serversService{repo: repo, rdb: rdb}
}

func (s *serversService) SaveServer(ctx context.Context, dto *dtos.ServersSaveDTO) (*dtos.ServersDTO, error) {
	if server, err := s.repo.FindByName(ctx, dto.Name); err != nil {
		log.Printf("[services:servers] failed to find server by name %s: %v", dto.Name, err)
		return nil, errors.New("failed to save server: "+err.Error(), http.StatusInternalServerError)
	} else if server != nil {
		log.Printf("[services:servers] server already exists: %s", dto.Name)
		return nil, errors.ErrServerExists
	}

	now := time.Now()
	server := entities.Servers{
		Name:      dto.Name,
		CreatedAt: now,
	}

	created, err := s.repo.Save(ctx, &server)
	if err != nil {
		log.Printf("[services:servers] failed to save server to DB: %v", err)
		return nil, errors.New("failed to save server: "+err.Error(), http.StatusInternalServerError)
	}

	log.Printf("[services:servers] server created successfully: %s", dto.Name)
	return helpers.ToServersDTO(created), nil
}

func (s *serversService) GetServerByID(ctx context.Context, id uuid.UUID) (*dtos.ServersDTO, error) {
	key := cache.KeyServer(id.String())
	if cached, err := cache.Get[dtos.ServersDTO](ctx, s.rdb, key); err != nil {
		log.Printf("[services:servers] cache read error for key %s: %v", key, err)
	} else if cached != nil {
		return cached, nil
	}

	server, err := s.repo.FindByID(ctx, id)
	if err != nil {
		log.Printf("[services:servers] failed to find server %s: %v", id, err)
		return nil, errors.New("failed to find server: "+err.Error(), http.StatusInternalServerError)
	}

	if server == nil {
		log.Printf("[services:servers] server not found: %s", id)
		return nil, errors.ErrServerNotFound
	}

	dto := helpers.ToServersDTO(server)
	cache.Set(ctx, s.rdb, key, dto, cache.TTLServer)
	log.Printf("[services:servers] server found: %s", id)
	return dto, nil
}

func (s *serversService) GetAllServers(ctx context.Context) ([]dtos.ServersDTO, error) {
	servers, err := s.repo.FindAll(ctx)
	if err != nil {
		log.Printf("[services:servers] failed to list servers: %v", err)
		return nil, errors.New("failed to list servers: "+err.Error(), http.StatusInternalServerError)
	}

	if servers == nil {
		log.Println("[services:servers] no servers found")
		return nil, errors.ErrServerNotFound
	}

	serverDTOs := make([]dtos.ServersDTO, len(servers))
	for i, d := range servers {
		serverDTOs[i] = *helpers.ToServersDTO(&d)
	}

	log.Printf("[services:servers] %d servers listed successfully", len(servers))
	return serverDTOs, nil
}

func (s *serversService) UpdateServer(ctx context.Context, id uuid.UUID, dto *dtos.ServersSaveDTO) (*dtos.ServersDTO, error) {
	server, err := s.repo.FindByID(ctx, id)
	if err != nil {
		log.Printf("[services:servers] failed to find server %s to update: %v", id, err)
		return nil, errors.New("failed to update server: "+err.Error(), http.StatusInternalServerError)
	}

	if server == nil {
		log.Printf("[services:servers] server not found to update: %s", id)
		return nil, errors.ErrServerNotFound
	}

	server.Name = dto.Name

	updated, err := s.repo.Update(ctx, server)
	if err != nil {
		log.Printf("[services:servers] failed to update server %s: %v", id, err)
		return nil, errors.New("failed to update server: "+err.Error(), http.StatusInternalServerError)
	}

	cache.Delete(ctx, s.rdb, cache.KeyServer(id.String()))
	log.Printf("[services:servers] server updated successfully: %s", id)
	return helpers.ToServersDTO(updated), nil
}

func (s *serversService) DeleteServer(ctx context.Context, id uuid.UUID) (bool, error) {
	server, err := s.repo.FindByID(ctx, id)
	if err != nil {
		log.Printf("[services:servers] failed to find server %s to delete: %v", id, err)
		return false, errors.New("failed to delete server: "+err.Error(), http.StatusInternalServerError)
	}

	if server == nil {
		log.Printf("[services:servers] server not found to delete: %s", id)
		return false, errors.ErrServerNotFound
	}

	err = s.repo.Delete(ctx, id)
	if err != nil {
		log.Printf("[services:servers] failed to delete server %s: %v", id, err)
		return false, errors.New("failed to delete server: "+err.Error(), http.StatusInternalServerError)
	}

	log.Printf("[services:servers] server deleted successfully: %s", id)
	return true, nil
}
