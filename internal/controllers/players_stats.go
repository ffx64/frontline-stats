package controllers

import (
	"github.com/ffx64/frontline-stats/internal/services"
	"github.com/gin-gonic/gin"
)

type PlayersStatsControllers struct {
	service services.PlayersStatsService
}

func NewPlayersStatsController(playersStatsService services.PlayersStatsService) *PlayersStatsControllers {
	return &PlayersStatsControllers{
		service: playersStatsService,
	}
}

func (c *PlayersStatsControllers) GetLeaderboard(ctx *gin.Context) {
	stats, err := c.service.GetLeaderboard(ctx)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(200, stats)
}

func (c *PlayersStatsControllers) GetLeaderboardVehicle(ctx *gin.Context) {
	stats, err := c.service.GetLeaderboardVehicle(ctx)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(200, stats)
}

func (c *PlayersStatsControllers) GetLeaderboardHeadshots(ctx *gin.Context) {
	stats, err := c.service.GetLeaderboardHeadshots(ctx)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(200, stats)
}
