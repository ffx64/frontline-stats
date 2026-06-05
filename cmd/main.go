package main

import (
	"log"
	"time"

	"github.com/ffx64/gamestats-backend/cmd/gin"
	"github.com/ffx64/gamestats-backend/internal/database"
	"github.com/ffx64/gamestats-backend/internal/repositories"
	"github.com/ffx64/gamestats-backend/internal/scheduler"
	"github.com/ffx64/gamestats-backend/internal/services"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("[main] .env not found, using default environment variables")
	}

	loc, err := time.LoadLocation("America/Sao_Paulo")
	if err != nil {
		log.Fatalf("[main] failed to set timezone: %v", err)
	}
	time.Local = loc

	db := database.Connect()
	log.Println("[main] database connection established")

	rdb := database.NewRedisClient()
	log.Println("[main] redis connection established")

	playersRepo := repositories.NewPlayersRepository(db)
	playersStatsRepo := repositories.NewPlayersStatsRepository(db)
	serversRepo := repositories.NewServersRepository(db)
	roundsRepo := repositories.NewRoundsRepository(db)
	roundsStatsRepo := repositories.NewRoundsStatsRepository(db)
	killsRepo := repositories.NewKillsRepository(db)
	log.Println("[main] repositories initialized")

	playersService := services.NewPlayersService(playersRepo, playersStatsRepo, rdb)
	playersStatsService := services.NewPlayersStatsService(playersStatsRepo, rdb)
	serversService := services.NewServersService(serversRepo, rdb)
	roundsService := services.NewRoundsService(roundsRepo, roundsStatsRepo, rdb)
	killsService := services.NewKillsService(killsRepo, roundsRepo, serversRepo, playersRepo, rdb)
	log.Println("[main] services initialized")

	s := scheduler.NewScheduler(playersStatsRepo, roundsStatsRepo, rdb)
	s.Start()
	defer s.Stop()
	log.Println("[main] scheduler started")

	gin.NewGinServer(playersService, serversService, roundsService, killsService, playersStatsService).Start()
}
