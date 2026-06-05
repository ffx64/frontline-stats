package services_test

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/ffx64/frontline-stats/internal/database"
	"github.com/ffx64/frontline-stats/internal/dtos"
	"github.com/ffx64/frontline-stats/internal/entities"
	"github.com/ffx64/frontline-stats/internal/repositories"
	"github.com/ffx64/frontline-stats/internal/services"
	"github.com/google/uuid"
)

func setupRoundsDB(t *testing.T) (repositories.RoundsRepository, repositories.RoundsStatsRepository) {
	t.Helper()

	db := database.Connect()

	_ = db.Migrator().DropTable(&entities.Rounds{})
	if err := db.AutoMigrate(&entities.Rounds{}); err != nil {
		t.Fatalf("erro ao migrar entidades: %v", err)
	}

	repo := repositories.NewRoundsRepository(db)
	statsRepo := repositories.NewRoundsStatsRepository(db)

	return repo, statsRepo
}

func TestRoundsService_SaveRound(t *testing.T) {
	log.Printf("info: iniciando teste - salvar rodada")

	repo, statsRepo := setupRoundsDB(t)
	service := services.NewRoundsService(repo, statsRepo)
	ctx := context.Background()

	serverID := uuid.New()
	dto := &dtos.RoundsCreateDTO{
		ServerID:      serverID.String(),
		CurrentMode:   "CTI",
		MissionHeader: "Operação Falcoes de Moscow",
	}

	created, err := service.SaveRound(ctx, dto)
	if err != nil {
		t.Fatalf("falha ao salvar rodada: %v", err)
	}

	if created.ServerID != serverID.String() {
		t.Errorf("esperava server_id %s, recebeu %s", serverID, created.ServerID)
	}

	if created.Status != "in_progress" {
		t.Errorf("esperava status 'in_progress', recebeu '%s'", created.Status)
	}
}

func TestRoundsService_GetRoundByID(t *testing.T) {
	log.Printf("info: iniciando teste - buscar rodada por ID")

	repo, statsRepo := setupRoundsDB(t)
	service := services.NewRoundsService(repo, statsRepo)
	ctx := context.Background()

	round := &entities.Rounds{
		ID:            uuid.New(),
		ServerID:      uuid.New(),
		CurrentMode:   "AAS",
		MissionHeader: "Mission Bravo",
		Status:        "in_progress",
		StartAt:       time.Now(),
		CreatedAt:     time.Now(),
	}

	repo.Save(ctx, round)

	found, err := service.GetRoundByID(ctx, round.ID)
	if err != nil {
		t.Fatalf("falha ao buscar rodada: %v", err)
	}

	if found.ID != round.ID.String() {
		t.Errorf("esperava id %s, recebeu %s", round.ID, found.ID)
	}
}

func TestRoundsService_UpdateRoundEnded(t *testing.T) {
	log.Printf("info: iniciando teste - atualizar rodada para 'ended'")

	repo, statsRepo := setupRoundsDB(t)
	service := services.NewRoundsService(repo, statsRepo)
	ctx := context.Background()

	round := &entities.Rounds{
		ID:            uuid.New(),
		ServerID:      uuid.New(),
		CurrentMode:   "CTI",
		MissionHeader: "Mission to End",
		Status:        "in_progress",
		StartAt:       time.Now(),
		CreatedAt:     time.Now(),
	}

	repo.Save(ctx, round)

	dto := &dtos.RoundsUpdatedEndedDTO{
		WinnerFaction: "BLUFOR",
	}

	updated, err := service.UpdateRoundEnded(ctx, round.ID, dto)
	if err != nil {
		t.Fatalf("falha ao atualizar rodada: %v", err)
	}

	if updated.Status != "ended" {
		t.Errorf("esperava status 'ended', recebeu '%s'", updated.Status)
	}

	if updated.WinnerFaction != "BLUFOR" {
		t.Errorf("esperava WinnerFaction 'BLUFOR', recebeu '%s'", updated.WinnerFaction)
	}

	if updated.EndedAt == nil {
		t.Errorf("esperava EndedAt não nulo, mas recebeu nil")
	}
}
