package controllers

import (
	"log"
	"net/http"

	"github.com/ffx64/gamestats-backend/internal/dtos"
	"github.com/ffx64/gamestats-backend/internal/errors"
	"github.com/ffx64/gamestats-backend/internal/services"
	"github.com/gin-gonic/gin"
)

type PlayersControllers struct {
	service services.PlayersService
}

func NewPlayersControllers(service services.PlayersService) *PlayersControllers {
	return &PlayersControllers{service: service}
}

func (c *PlayersControllers) SavePlayer(ctx *gin.Context) {
	var dto dtos.PlayerSaveDTO
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		log.Printf("[controller:players] invalid JSON format: %v", err)
		ctx.JSON(errors.ErrJsonInvalidFormat.Status, errors.ErrJsonInvalidFormat)
		return
	}

	player, err := c.service.Save(ctx, &dto)
	if err != nil {
		if custom, ok := err.(*errors.AppError); ok {
			ctx.JSON(custom.Status, custom)
			return
		}
		ctx.JSON(http.StatusInternalServerError, errors.New(err.Error(), http.StatusInternalServerError))
		return
	}

	log.Printf("[controller:players] player created: %v", player.ID)
	ctx.JSON(http.StatusCreated, player)
}

func (c *PlayersControllers) GetPlayerByGUID(ctx *gin.Context) {
	guid := ctx.Param("guid")

	player, err := c.service.GetByGUID(ctx, guid)
	if err != nil {
		if custom, ok := err.(*errors.AppError); ok {
			ctx.JSON(custom.Status, custom)
			return
		}
		ctx.JSON(http.StatusInternalServerError, errors.New(err.Error(), http.StatusInternalServerError))
		return
	}

	log.Printf("[controller:players] player retrieved: %v", player.ID)
	ctx.JSON(http.StatusOK, player)
}

func (c *PlayersControllers) IfNotExistsCreatePlayer(ctx *gin.Context) {
	guid := ctx.Param("guid")
	serverLastID := ctx.Param("serverLastID")
	username := ctx.Param("username")

	player, err := c.service.IfNotExistsCreate(ctx, username, guid, serverLastID)
	if err != nil {
		if custom, ok := err.(*errors.AppError); ok {
			ctx.JSON(custom.Status, custom)
			return
		}
		ctx.JSON(http.StatusInternalServerError, errors.New(err.Error(), http.StatusInternalServerError))
		return
	}

	log.Printf("[controller:players] player retrieved: %v", player.ID)
	ctx.JSON(http.StatusOK, player)
}

func (c *PlayersControllers) GetPlayerStatsByGUID(ctx *gin.Context) {
	guid := ctx.Param("guid")

	stats, err := c.service.GetPlayerStatsByGUID(ctx, guid)
	if err != nil {
		if custom, ok := err.(*errors.AppError); ok {
			ctx.JSON(custom.Status, custom)
			return
		}
		ctx.JSON(http.StatusInternalServerError, errors.New(err.Error(), http.StatusInternalServerError))
		return
	}

	log.Printf("[controller:players_stats] player stats retrieved: %v", stats.GUID)
	ctx.JSON(http.StatusOK, stats)
}

func (c *PlayersControllers) UpdatePlayer(ctx *gin.Context) {
	guid := ctx.Param("guid")
	var dto dtos.PlayerUpdateDTO

	if err := ctx.ShouldBindJSON(&dto); err != nil {
		log.Printf("[controller:players] invalid JSON format for update: %v", err)
		ctx.JSON(errors.ErrJsonInvalidFormat.Status, errors.ErrJsonInvalidFormat)
		return
	}

	updated, err := c.service.Update(ctx, guid, &dto)
	if err != nil {
		if custom, ok := err.(*errors.AppError); ok {
			ctx.JSON(custom.Status, custom)
			return
		}
		ctx.JSON(http.StatusInternalServerError, errors.New(err.Error(), http.StatusInternalServerError))
		return
	}

	log.Printf("[controller:players] player updated: %v", updated.ID)
	ctx.JSON(http.StatusOK, updated)
}
