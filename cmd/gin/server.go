package gin

import (
	"log"
	"os"

	"github.com/ffx64/gamestats-backend/internal/controllers"
	"github.com/ffx64/gamestats-backend/internal/middleware"
	"github.com/ffx64/gamestats-backend/internal/services"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type GinInitializer struct {
	engine *gin.Engine
	port   string
}

func NewGinServer(
	playersService services.PlayersService,
	serversService services.ServersService,
	roundsService services.RoundsService,
	killsService services.KillsService,
	playersStatsService services.PlayersStatsService) *GinInitializer {

	appEnv := os.Getenv("APP_ENV")
	switch appEnv {
	case "dev":
		gin.SetMode(gin.DebugMode)
	case "test":
		gin.SetMode(gin.TestMode)
	default:
		gin.SetMode(gin.ReleaseMode)
	}
	log.Printf("[gin:server] environment set: %s", appEnv)

	port := os.Getenv("GIN_PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("[gin:server] port configured: %s", port)

	router := gin.Default()

	key := os.Getenv("API_KEY")
	if key != "" {
		log.Printf("[gin:server] api key configured: %s", key)
		router.Use(middleware.Auth(key))
		log.Printf("[gin:server] auth middleware loaded")
	}

	router.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
		AllowMethods:    []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:    []string{"Origin", "Content-Type", "Authorization"},
	}))
	log.Println("[gin:server] CORS middleware loaded")

	playersController := controllers.NewPlayersControllers(playersService)
	playersStatsController := controllers.NewPlayersStatsController(playersStatsService)
	serversController := controllers.NewServersController(serversService)
	roundsController := controllers.NewRoundsController(roundsService)
	killsController := controllers.NewKillsController(killsService)
	log.Println("[gin:server] controllers initialized")

	v1 := router.Group("/api/v1")

	players := v1.Group("/players")
	servers := v1.Group("/servers")
	rounds := v1.Group("/rounds")
	events := v1.Group("/events")

	leaderboard := v1.Group("/leaderboard")
	leaderboard.GET("", playersStatsController.GetLeaderboard)
	leaderboard.GET("/vehicles", playersStatsController.GetLeaderboardVehicle)
	leaderboard.GET("/headshots", playersStatsController.GetLeaderboardHeadshots)

	players.GET("/:guid", playersController.GetPlayerByGUID)
	players.POST("", playersController.SavePlayer)
	players.PUT("/:guid", playersController.UpdatePlayer)
	players.GET("/:guid/stats", playersController.GetPlayerStatsByGUID)

	players.GET("/if-not-exists-create/:username/:guid/:serverLastID", playersController.IfNotExistsCreatePlayer)

	servers.POST("", serversController.SaveServer)
	servers.GET("", serversController.GetAllServers)
	servers.GET("/:id", serversController.GetServerByID)
	servers.PUT("/:id", serversController.UpdateServer)
	servers.DELETE("/:id", serversController.DeleteServer)

	rounds.POST("", roundsController.SaveRound)
	rounds.GET("/:id", roundsController.GetRoundByID)
	rounds.PUT("/:id/ended", roundsController.UpdateRoundEnded)
	rounds.GET("/:id/scoreboard", roundsController.GetScoreboardByRoundID)
	rounds.GET("/server/:serverId/player/:playerId", roundsController.GetAllRoundsByServerIDAndPlayerID)

	events.POST("/kill", killsController.SaveKills)
	log.Println("[gin:server] routes registered")

	return &GinInitializer{engine: router, port: port}
}

func (g *GinInitializer) Start() {
	log.Printf("[gin:server] server started on port %s", g.port)
	if err := g.engine.Run(":" + g.port); err != nil {
		log.Fatalf("[gin:server] failed to start server: %v", err)
	}
}
