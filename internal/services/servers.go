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

type ServersService interface {
	SaveServer(ctx context.Context, dto *dtos.ServersSaveDTO) (*dtos.ServersDTO, error)
	GetServerByID(ctx context.Context, id uuid.UUID) (*dtos.ServersDTO, error)
	GetAllServers(ctx context.Context) ([]dtos.ServersDTO, error)
	UpdateServer(ctx context.Context, id uuid.UUID, server *dtos.ServersSaveDTO) (*dtos.ServersDTO, error)
	DeleteServer(ctx context.Context, id uuid.UUID) (bool, error)
}

type serversService struct {
	repo repositories.ServersRepository
}

func NewServersService(repo repositories.ServersRepository) ServersService {
	return &serversService{repo: repo}
}

func (s *serversService) SaveServer(ctx context.Context, dto *dtos.ServersSaveDTO) (*dtos.ServersDTO, error) {
	if server, err := s.repo.FindByName(ctx, dto.Name); err != nil {
		log.Printf("[services:servers] erro ao buscar servidor pelo nome %s: %v", dto.Name, err)
		return nil, errors.New("falha ao salvar servidor: "+err.Error(), http.StatusInternalServerError)
	} else if server != nil {
		log.Printf("[services:servers] servidor já existe: %s", dto.Name)
		return nil, errors.ErrServerExists
	}

	now := time.Now()
	server := entities.Servers{
		Name:      dto.Name,
		CreatedAt: now,
	}

	created, err := s.repo.Save(ctx, &server)
	if err != nil {
		log.Printf("[services:servers] erro ao salvar servidor no DB: %v", err)
		return nil, errors.New("falha ao salvar servidor: "+err.Error(), http.StatusInternalServerError)
	}

	log.Printf("[services:servers] servidor criado com sucesso: %s", dto.Name)
	return helpers.ToServersDTO(created), nil
}

func (s *serversService) GetServerByID(ctx context.Context, id uuid.UUID) (*dtos.ServersDTO, error) {
	server, err := s.repo.FindByID(ctx, id)
	if err != nil {
		log.Printf("[services:servers] erro ao buscar servidor %s: %v", id, err)
		return nil, errors.New("falha ao buscar servidor: "+err.Error(), http.StatusInternalServerError)
	}

	if server == nil {
		log.Printf("[services:servers] servidor não encontrado: %s", id)
		return nil, errors.ErrServerNotFound
	}

	log.Printf("[services:servers] servidor encontrado: %s", id)
	return helpers.ToServersDTO(server), nil
}

func (s *serversService) GetAllServers(ctx context.Context) ([]dtos.ServersDTO, error) {
	servers, err := s.repo.FindAll(ctx)
	if err != nil {
		log.Printf("[services:servers] erro ao listar servidores: %v", err)
		return nil, errors.New("falha ao buscar servidores: "+err.Error(), http.StatusInternalServerError)
	}

	if servers == nil {
		log.Println("[services:servers] nenhum servidor encontrado")
		return nil, errors.ErrServerNotFound
	}

	serverDTOs := make([]dtos.ServersDTO, len(servers))
	for i, d := range servers {
		serverDTOs[i] = *helpers.ToServersDTO(&d)
	}

	log.Printf("[services:servers] %d servidores listados com sucesso", len(servers))
	return serverDTOs, nil
}

func (s *serversService) UpdateServer(ctx context.Context, id uuid.UUID, dto *dtos.ServersSaveDTO) (*dtos.ServersDTO, error) {
	server, err := s.repo.FindByID(ctx, id)
	if err != nil {
		log.Printf("[services:servers] erro ao buscar servidor %s para atualizar: %v", id, err)
		return nil, errors.New("falha ao atualizar servidor: "+err.Error(), http.StatusInternalServerError)
	}

	if server == nil {
		log.Printf("[services:servers] servidor não encontrado para atualizar: %s", id)
		return nil, errors.ErrServerNotFound
	}

	server.Name = dto.Name

	updated, err := s.repo.Update(ctx, server)
	if err != nil {
		log.Printf("[services:servers] erro ao atualizar servidor %s: %v", id, err)
		return nil, errors.New("falha ao atualizar servidor: "+err.Error(), http.StatusInternalServerError)
	}

	log.Printf("[services:servers] servidor atualizado com sucesso: %s", id)
	return helpers.ToServersDTO(updated), nil
}

func (s *serversService) DeleteServer(ctx context.Context, id uuid.UUID) (bool, error) {
	server, err := s.repo.FindByID(ctx, id)
	if err != nil {
		log.Printf("[services:servers] erro ao buscar servidor %s para deletar: %v", id, err)
		return false, errors.New("falha ao deletar servidor: "+err.Error(), http.StatusInternalServerError)
	}

	if server == nil {
		log.Printf("[services:servers] servidor não encontrado para deletar: %s", id)
		return false, errors.ErrServerNotFound
	}

	err = s.repo.Delete(ctx, id)
	if err != nil {
		log.Printf("[services:servers] erro ao deletar servidor %s: %v", id, err)
		return false, errors.New("falha ao deletar servidor: "+err.Error(), http.StatusInternalServerError)
	}

	log.Printf("[services:servers] servidor deletado com sucesso: %s", id)
	return true, nil
}
