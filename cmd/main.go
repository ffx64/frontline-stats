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
		log.Println("[main] .env não encontrado, usando variáveis padrão")
	}

	loc, err := time.LoadLocation("America/Sao_Paulo")
	if err != nil {
		log.Fatalf("[main] erro ao definir timezone: %v", err)
	}
	time.Local = loc

	db := database.Connect()
	log.Println("[main] conexão com o banco de dados estabelecida")

	rdb := database.NewRedisClient()
	log.Println("[main] conexão com o Redis estabelecida")

	playersRepo := repositories.NewPlayersRepository(db)
	playersStatsRepo := repositories.NewPlayersStatsRepository(db)
	serversRepo := repositories.NewServersRepository(db)
	roundsRepo := repositories.NewRoundsRepository(db)
	roundsStatsRepo := repositories.NewRoundsStatsRepository(db)
	killsRepo := repositories.NewKillsRepository(db)
	log.Println("[main] repositórios inicializados")

	playersService := services.NewPlayersService(playersRepo, playersStatsRepo)
	playersStatsService := services.NewPlayersStatsService(playersStatsRepo)
	serversService := services.NewServersService(serversRepo)
	roundsService := services.NewRoundsService(roundsRepo, roundsStatsRepo)
	killsService := services.NewKillsService(killsRepo, roundsRepo, serversRepo, playersRepo, rdb)
	log.Println("[main] serviços inicializados")

	s := scheduler.NewScheduler(playersStatsRepo, roundsStatsRepo)
	s.Start()
	defer s.Stop()
	log.Println("[main] scheduler iniciado")

	gin.NewGinServer(playersService, serversService, roundsService, killsService, playersStatsService).Start()
}
