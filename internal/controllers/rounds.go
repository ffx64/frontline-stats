package controllers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/ffx64/gamestats-backend/internal/dtos"
	"github.com/ffx64/gamestats-backend/internal/errors"
	"github.com/ffx64/gamestats-backend/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type RoundsController struct {
	service services.RoundsService
}

func NewRoundsController(service services.RoundsService) *RoundsController {
	return &RoundsController{service: service}
}

// POST /api/v1/rounds
func (c *RoundsController) SaveRound(ctx *gin.Context) {
	var dto dtos.RoundsCreateDTO
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		log.Printf("[controller:rounds] invalid JSON format: %v", err)
		ctx.JSON(errors.ErrJsonInvalidFormat.Status, errors.ErrJsonInvalidFormat)
		return
	}

	saved, err := c.service.SaveRound(ctx, &dto)
	if err != nil {
		if custom, ok := err.(*errors.AppError); ok {
			log.Printf("[controller:rounds] failed to save round: %v", custom)
			ctx.JSON(custom.Status, custom)
			return
		}
		log.Printf("[controller:rounds] failed to save round: %v", err)
		ctx.JSON(http.StatusInternalServerError, errors.New(err.Error(), http.StatusInternalServerError))
		return
	}

	log.Printf("[controller:rounds] round created: %v", saved.ID)
	ctx.JSON(http.StatusCreated, saved)
}

// GET /api/v1/rounds/:id
func (c *RoundsController) GetRoundByID(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		log.Printf("[controller:rounds] invalid UUID: %v", err)
		ctx.JSON(errors.ErrUUIDError.Status, errors.ErrUUIDError)
		return
	}

	round, err := c.service.GetRoundByID(ctx, id)
	if err != nil {
		if custom, ok := err.(*errors.AppError); ok {
			log.Printf("[controller:rounds] failed to get round: %v", custom.Message)
			ctx.JSON(custom.Status, custom)
			return
		}
		log.Printf("[controller:rounds] failed to get round: %v", err)
		ctx.JSON(http.StatusInternalServerError, errors.New(err.Error(), http.StatusInternalServerError))
		return
	}

	log.Printf("[controller:rounds] round retrieved: %v", round.ID)
	ctx.JSON(http.StatusOK, round)
}

// GET /api/v1/rounds/:id/scoreboard
func (c *RoundsController) GetScoreboardByRoundID(ctx *gin.Context) {
	roundIdParam := ctx.Param("id")
	roundId, err := uuid.Parse(roundIdParam)
	if err != nil {
		log.Printf("[controllers:rounds_stats] invalid round UUID: %v", roundIdParam)
		ctx.JSON(http.StatusBadRequest, errors.ErrUUIDError)
		return
	}

	scoreboard, err := c.service.GetScoreboardByRoundID(ctx, roundId)
	if err != nil {
		log.Printf("[controllers:rounds_stats] failed to get round scoreboard for round_id %v: %v", roundId, err)
		ctx.JSON(http.StatusNotFound, errors.ErrRoundNotFound)
		return
	}

	ctx.JSON(http.StatusOK, scoreboard)
}

// GET /api/v1/rounds/servers/:serverId/players/:playerId?limit=&offset=
func (c *RoundsController) GetAllRoundsByServerIDAndPlayerID(ctx *gin.Context) {
	serverIdParam := ctx.Param("serverId")
	serverId, err := uuid.Parse(serverIdParam)
	if err != nil {
		log.Printf("[controllers:rounds_stats] invalid server UUID: %v", serverIdParam)
		ctx.JSON(http.StatusBadRequest, errors.ErrUUIDError)
		return
	}

	playerIdParam := ctx.Param("playerId")
	playerId, err := uuid.Parse(playerIdParam)
	if err != nil {
		log.Printf("[controllers:rounds_stats] invalid player UUID: %v", playerIdParam)
		ctx.JSON(http.StatusBadRequest, errors.ErrUUIDError)
		return
	}

	limit, err := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
	if err != nil {
		log.Printf("[controllers:rounds_stats] failed to parse pagination params: %v", err)
		ctx.JSON(http.StatusBadRequest, errors.ErrConvertParam)
		return
	}

	offset, err := strconv.Atoi(ctx.DefaultQuery("offset", "0"))
	if err != nil {
		log.Printf("[controllers:rounds_stats] failed to parse pagination params: %v", err)
		ctx.JSON(http.StatusBadRequest, errors.ErrConvertParam)
		return
	}

	statsList, err := c.service.GetAllRoundsByServerIDAndPlayerID(ctx, serverId, playerId, limit, offset)
	if err != nil {
		log.Printf("[controllers:rounds_stats] failed to get rounds for server_id %v and player_id %v: %v", serverId, playerId, err)
		ctx.JSON(http.StatusNotFound, errors.ErrRoundsNotFound)
		return
	}

	ctx.JSON(http.StatusOK, statsList)
}

// PUT /api/v1/rounds/:id/ended
func (c *RoundsController) UpdateRoundEnded(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		log.Printf("[controller:rounds] invalid UUID for update: %v", err)
		ctx.JSON(errors.ErrUUIDError.Status, errors.ErrUUIDError)
		return
	}

	var dto dtos.RoundsUpdatedEndedDTO
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		log.Printf("[controller:rounds] invalid JSON format for update: %v", err)
		ctx.JSON(errors.ErrJsonInvalidFormat.Status, errors.ErrJsonInvalidFormat)
		return
	}

	updated, err := c.service.UpdateRoundEnded(ctx, id, &dto)
	if err != nil {
		if custom, ok := err.(*errors.AppError); ok {
			log.Printf("[controller:rounds] failed to update round: %v", custom)
			ctx.JSON(custom.Status, custom)
			return
		}
		log.Printf("[controller:rounds] failed to update round: %v", err)
		ctx.JSON(http.StatusInternalServerError, errors.New(err.Error(), http.StatusInternalServerError))
		return
	}

	log.Printf("[controller:rounds] round updated: %v", updated.ID)
	ctx.JSON(http.StatusOK, updated)
}
