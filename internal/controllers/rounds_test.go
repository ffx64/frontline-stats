package controllers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ffx64/gamestats-backend/internal/controllers"
	"github.com/ffx64/gamestats-backend/internal/database"
	"github.com/ffx64/gamestats-backend/internal/dtos"
	"github.com/ffx64/gamestats-backend/internal/entities"
	"github.com/ffx64/gamestats-backend/internal/repositories"
	"github.com/ffx64/gamestats-backend/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func setupRoundsController() (*gin.Engine, *controllers.RoundsController) {
	gin.SetMode(gin.TestMode)

	db := database.Connect()
	_ = db.Migrator().DropTable(&entities.Rounds{})
	_ = db.AutoMigrate(&entities.Rounds{})

	repo := repositories.NewRoundsRepository(db)
	statsRepo := repositories.NewRoundsStatsRepository(db)

	service := services.NewRoundsService(repo, statsRepo)
	controller := controllers.NewRoundsController(service)

	r := gin.Default()
	r.POST("/api/v1/rounds", controller.SaveRound)
	r.GET("/api/v1/rounds/:id", controller.GetRoundByID)
	r.PUT("/api/v1/rounds/:id/end", controller.UpdateRoundEnded)

	return r, controller
}

func TestSaveRound(t *testing.T) {
	router, _ := setupRoundsController()
	log.Println("info: iniciando teste - salvar rodada")

	dto := dtos.RoundsCreateDTO{
		ServerID:      uuid.New().String(),
		CurrentMode:   "CTF",
		MissionHeader: "Operation Snakebite",
	}

	body, _ := json.Marshal(dto)
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/rounds", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusCreated {
		log.Printf("error: esperado status 201, recebido %d", resp.Code)
		t.Fatalf("esperava 201, recebeu %d", resp.Code)
	}

	var result dtos.RoundsDTO
	if err := json.Unmarshal(resp.Body.Bytes(), &result); err != nil {
		log.Printf("error: falha ao decodificar resposta: %v", err)
		t.Fatal(err)
	}

	if result.CurrentMode != dto.CurrentMode || result.MissionHeader != dto.MissionHeader {
		log.Printf("error: dados incorretos: esperado %v, recebido %v", dto, result)
		t.Fatalf("dados incorretos")
	}

	log.Println("info: teste finalizado com sucesso - salvar rodada")
}

func TestGetRoundByID(t *testing.T) {
	router, _ := setupRoundsController()
	log.Println("info: iniciando teste - obter rodada por id")

	db := database.Connect()
	repo := repositories.NewRoundsRepository(db)

	ctx := context.TODO()
	round := &entities.Rounds{
		ID:            uuid.New(),
		ServerID:      uuid.New(),
		CurrentMode:   "TDM",
		MissionHeader: "Desert Raid",
		Status:        "in_progress",
		StartAt:       time.Now(),
		CreatedAt:     time.Now(),
	}
	repo.Save(ctx, round)

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/rounds/"+round.ID.String(), nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		log.Printf("error: esperado status 200, recebido %d", resp.Code)
		t.Fatalf("esperava 200, recebeu %d", resp.Code)
	}

	var result dtos.RoundsDTO
	if err := json.Unmarshal(resp.Body.Bytes(), &result); err != nil {
		log.Printf("error: falha ao decodificar resposta: %v", err)
		t.Fatal(err)
	}

	if result.CurrentMode != round.CurrentMode || result.MissionHeader != round.MissionHeader {
		log.Printf("error: dados incorretos: esperado %v, recebido %v", round, result)
		t.Fatalf("dados incorretos")
	}

	log.Println("info: teste finalizado com sucesso - obter rodada por id")
}

func TestGetAllRounds(t *testing.T) {
	router, _ := setupRoundsController()
	log.Println("info: iniciando teste - listar todas as rodadas")

	ctx := context.TODO()
	db := database.Connect()
	repo := repositories.NewRoundsRepository(db)

	repo.Save(ctx, &entities.Rounds{
		ID:            uuid.New(),
		ServerID:      uuid.New(),
		CurrentMode:   "AAS",
		MissionHeader: "Urban Clash",
		Status:        "in_progress",
		StartAt:       time.Now(),
		CreatedAt:     time.Now(),
	})
	repo.Save(ctx, &entities.Rounds{
		ID:            uuid.New(),
		ServerID:      uuid.New(),
		CurrentMode:   "CTF",
		MissionHeader: "Forest Ambush",
		Status:        "in_progress",
		StartAt:       time.Now(),
		CreatedAt:     time.Now(),
	})

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/rounds", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		log.Printf("error: esperado status 200, recebido %d", resp.Code)
		t.Fatalf("esperava 200, recebeu %d", resp.Code)
	}

	var list []dtos.RoundsDTO
	if err := json.Unmarshal(resp.Body.Bytes(), &list); err != nil {
		log.Printf("error: falha ao decodificar resposta: %v", err)
		t.Fatal(err)
	}

	if len(list) < 2 {
		log.Printf("error: esperado pelo menos 2 rodadas, recebidas %d", len(list))
		t.Fatalf("quantidade incorreta de rodadas")
	}

	log.Println("info: teste finalizado com sucesso - listar todas as rodadas")
}

func TestUpdateRoundEnded(t *testing.T) {
	router, _ := setupRoundsController()
	log.Println("info: iniciando teste - finalizar rodada")

	db := database.Connect()
	repo := repositories.NewRoundsRepository(db)
	ctx := context.TODO()

	round := &entities.Rounds{
		ID:            uuid.New(),
		ServerID:      uuid.New(),
		CurrentMode:   "TDM",
		MissionHeader: "Bridge Assault",
		Status:        "in_progress",
		StartAt:       time.Now(),
		CreatedAt:     time.Now(),
	}
	repo.Save(ctx, round)

	updateDTO := dtos.RoundsUpdatedEndedDTO{
		WinnerFaction: "US Army",
	}
	body, _ := json.Marshal(updateDTO)

	req, _ := http.NewRequest(http.MethodPut, "/api/v1/rounds/"+round.ID.String()+"/end", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		log.Printf("error: esperado status 200, recebido %d", resp.Code)
		t.Fatalf("esperava 200, recebeu %d", resp.Code)
	}

	var result dtos.RoundsDTO
	if err := json.Unmarshal(resp.Body.Bytes(), &result); err != nil {
		log.Printf("error: falha ao decodificar resposta: %v", err)
		t.Fatal(err)
	}

	if result.Status != "ended" || result.WinnerFaction != "US Army" || result.EndedAt == nil {
		log.Printf("error: rodada não finalizada corretamente: %+v", result)
		t.Fatalf("erro ao finalizar rodada")
	}

	log.Println("info: teste finalizado com sucesso - finalizar rodada")
}
