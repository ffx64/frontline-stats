package controllers

import (
	"log"
	"net/http"

	"github.com/ffx64/gamestats-backend/internal/dtos"
	"github.com/ffx64/gamestats-backend/internal/errors"
	"github.com/ffx64/gamestats-backend/internal/services"
	"github.com/gin-gonic/gin"
)

type KillsController struct {
	service services.KillsService
}

func NewKillsController(service services.KillsService) *KillsController {
	return &KillsController{service: service}
}

func (c *KillsController) SaveKills(ctx *gin.Context) {
	var dto []dtos.KillsSaveDTO

	if err := ctx.ShouldBindJSON(&dto); err != nil {
		log.Printf("[controllers/kills] failed to decode JSON body: %v", err)
		ctx.JSON(errors.ErrJsonInvalidFormat.Status, errors.ErrJsonInvalidFormat)
		return
	}

	if err := c.service.SaveKills(ctx, dto); err != nil {
		if custom, ok := err.(*errors.AppError); ok {
			log.Printf("[controllers/kills] failed to save kills: %v", custom.Message)
			ctx.JSON(custom.Status, custom)
			return
		}
		log.Printf("[controllers/kills] internal error saving kills: %v", err)
		ctx.JSON(http.StatusInternalServerError, errors.New(err.Error(), http.StatusInternalServerError))
		return
	}

	log.Printf("[controllers/kills] kills saved successfully - count=%d", len(dto))
	ctx.JSON(http.StatusCreated, gin.H{
		"message": "kills saved successfully",
		"count":   len(dto),
	})
}
