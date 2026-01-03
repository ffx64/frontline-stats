package controllers

import (
	"log"
	"net/http"

	"github.com/ffx64/gamestats-backend/internal/dtos"
	"github.com/ffx64/gamestats-backend/internal/errors"
	"github.com/ffx64/gamestats-backend/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ServersController struct {
	service services.ServersService
}

func NewServersController(service services.ServersService) *ServersController {
	return &ServersController{service: service}
}

// POST /api/v1/servers
func (c *ServersController) SaveServer(ctx *gin.Context) {
	var dto dtos.ServersSaveDTO
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		log.Printf("[controller:servers] formato de JSON inválido: %v", err)
		ctx.JSON(errors.ErrJsonInvalidFormat.Status, errors.ErrJsonInvalidFormat)
		return
	}

	saved, err := c.service.SaveServer(ctx, &dto)
	if err != nil {
		if custom, ok := err.(*errors.AppError); ok {
			ctx.JSON(custom.Status, custom)
			return
		}
		log.Printf("[controller:servers] erro ao salvar server: %v", err)
		ctx.JSON(http.StatusInternalServerError, errors.New(err.Error(), http.StatusInternalServerError))
		return
	}

	log.Printf("[controller:servers] server criado: %v", saved.ID)
	ctx.JSON(http.StatusCreated, saved)
}

// GET /api/v1/servers/:id
func (c *ServersController) GetServerByID(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		log.Printf("[controller:servers] UUID inválido: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "id inválido"})
		return
	}

	server, err := c.service.GetServerByID(ctx, id)
	if err != nil {
		if custom, ok := err.(*errors.AppError); ok {
			ctx.JSON(custom.Status, custom)
			return
		}
		log.Printf("[controller:servers] erro ao obter server: %v", err)
		ctx.JSON(http.StatusInternalServerError, errors.New(err.Error(), http.StatusInternalServerError))
		return
	}

	log.Printf("[controller:servers] server recuperado: %v", server.ID)
	ctx.JSON(http.StatusOK, server)
}

// GET /api/v1/servers
func (c *ServersController) GetAllServers(ctx *gin.Context) {
	servers, err := c.service.GetAllServers(ctx)
	if err != nil {
		if custom, ok := err.(*errors.AppError); ok {
			ctx.JSON(custom.Status, custom)
			return
		}
		log.Printf("[controller:servers] erro ao obter servers: %v", err)
		ctx.JSON(http.StatusInternalServerError, errors.New(err.Error(), http.StatusInternalServerError))
		return
	}

	log.Printf("[controller:servers] %d servers recuperados", len(servers))
	ctx.JSON(http.StatusOK, servers)
}

// PUT /api/v1/servers/:id
func (c *ServersController) UpdateServer(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		log.Printf("[controller:servers] UUID inválido para atualização: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "id inválido"})
		return
	}

	var dto dtos.ServersSaveDTO
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		log.Printf("[controller:servers] formato de JSON inválido para atualização: %v", err)
		ctx.JSON(errors.ErrJsonInvalidFormat.Status, errors.ErrJsonInvalidFormat)
		return
	}

	updated, err := c.service.UpdateServer(ctx, id, &dto)
	if err != nil {
		if custom, ok := err.(*errors.AppError); ok {
			ctx.JSON(custom.Status, custom)
			return
		}
		log.Printf("[controller:servers] erro ao atualizar server: %v", err)
		ctx.JSON(http.StatusInternalServerError, errors.New(err.Error(), http.StatusInternalServerError))
		return
	}

	log.Printf("[controller:servers] server atualizado: %v", updated.ID)
	ctx.JSON(http.StatusOK, updated)
}

// DELETE /api/v1/servers/:id
func (c *ServersController) DeleteServer(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		log.Printf("[controller:servers] UUID inválido para delete: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "id inválido"})
		return
	}

	deleted, err := c.service.DeleteServer(ctx, id)
	if err != nil {
		if custom, ok := err.(*errors.AppError); ok {
			ctx.JSON(custom.Status, custom)
			return
		}
		log.Printf("[controller:servers] erro ao deletar server: %v", err)
		ctx.JSON(http.StatusInternalServerError, errors.New(err.Error(), http.StatusInternalServerError))
		return
	}

	log.Printf("[controller:servers] server deletado: %v", id)
	ctx.JSON(http.StatusOK, deleted)
}
