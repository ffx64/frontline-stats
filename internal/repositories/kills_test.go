package repositories

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/google/uuid"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/ffx64/gamestats-backend/internal/entities"
)

// newTestDB creates an in-memory SQLite database for testing.
func newTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		log.Printf("error: falha ao abrir banco de dados em memória: %v", err)
		t.Fatalf("falha ao abrir banco de dados em memória: %v", err)
	}

	if err := db.AutoMigrate(&entities.Kills{}); err != nil {
		log.Printf("error: falha ao migrar schema: %v", err)
		t.Fatalf("falha ao migrar schema: %v", err)
	}

	log.Println("info: banco de dados de teste criado com sucesso")
	return db
}

// newTestKill creates a fake kill log for testing purposes.
func newTestKill() *entities.Kills {
	return &entities.Kills{
		ID:               uuid.New(),
		ServerID:         uuid.New(),
		RoundID:          uuid.New(),
		KillerID:         uuid.New(),
		VictimID:         uuid.New(),
		VictimWeaponName: "MX",
		VictimWeaponType: "Rifle",
		KillerWeaponName: "MX SW",
		KillerWeaponType: "LMG",
		HitZone:          "Head",
		Distance:         120.5,
		IsHeadshot:       true,
		IsFriendly:       false,
		IsVehicle:        false,
		KillerTeam:       "BLUFOR",
		VictimTeam:       "OPFOR",
		Timestamp:        time.Now(),
	}
}

func TestKillsRepository_Save(t *testing.T) {
	log.Println("info: iniciando teste - salvar log de kill")

	db := newTestDB(t)
	repo := NewKillsRepository(db)
	ctx := context.Background()

	kill := newTestKill()
	if err := repo.Save(ctx, kill); err != nil {
		log.Printf("error: falha ao salvar log de kill: %v", err)
		t.Fatalf("falha ao salvar log de kill: %v", err)
	}

	var count int64
	db.Model(&entities.Kills{}).Count(&count)
	if count != 1 {
		log.Printf("error: esperado 1 registro, recebido %d", count)
		t.Fatalf("esperado 1 registro, recebido %d", count)
	}

	log.Println("info: teste finalizado com sucesso - salvar log de kill")
}

func TestKillsRepository_GetKillsByPlayerID(t *testing.T) {
	log.Println("info: iniciando teste - obter kills por jogador")

	db := newTestDB(t)
	repo := NewKillsRepository(db)
	ctx := context.Background()

	kill := newTestKill()
	repo.Save(ctx, kill)

	kills, err := repo.GetKillsByPlayerID(ctx, kill.KillerID)
	if err != nil {
		log.Printf("error: falha na query: %v", err)
		t.Fatalf("falha na query: %v", err)
	}

	if len(kills) != 1 {
		log.Printf("error: esperado 1 kill, recebido %d", len(kills))
		t.Fatalf("esperado 1 kill, recebido %d", len(kills))
	}

	log.Println("info: teste finalizado com sucesso - obter kills por jogador")
}

func TestKillsRepository_GetDeathsByPlayerID(t *testing.T) {
	log.Println("info: iniciando teste - obter mortes por jogador")

	db := newTestDB(t)
	repo := NewKillsRepository(db)
	ctx := context.Background()

	kill := newTestKill()
	repo.Save(ctx, kill)

	deaths, err := repo.GetDeathsByPlayerID(ctx, kill.VictimID)
	if err != nil {
		log.Printf("error: falha na query: %v", err)
		t.Fatalf("falha na query: %v", err)
	}

	if len(deaths) != 1 {
		log.Printf("error: esperado 1 morte, recebido %d", len(deaths))
		t.Fatalf("esperado 1 morte, recebido %d", len(deaths))
	}

	log.Println("info: teste finalizado com sucesso - obter mortes por jogador")
}

func TestKillsRepository_GetKillsForPlayerByServerID(t *testing.T) {
	log.Println("info: iniciando teste - obter kills por servidor")

	db := newTestDB(t)
	repo := NewKillsRepository(db)
	ctx := context.Background()

	kill := newTestKill()
	repo.Save(ctx, kill)

	kills, err := repo.GetKillsForPlayerByServerID(ctx, kill.KillerID, kill.ServerID)
	if err != nil {
		log.Printf("error: falha na query: %v", err)
		t.Fatalf("falha na query: %v", err)
	}

	if len(kills) != 1 {
		log.Printf("error: esperado 1 kill, recebido %d", len(kills))
		t.Fatalf("esperado 1 kill, recebido %d", len(kills))
	}

	log.Println("info: teste finalizado com sucesso - obter kills por servidor")
}

func TestKillsRepository_GetKillsForPlayerByRoundID(t *testing.T) {
	log.Println("info: iniciando teste - obter kills por rodada")

	db := newTestDB(t)
	repo := NewKillsRepository(db)
	ctx := context.Background()

	kill := newTestKill()
	repo.Save(ctx, kill)

	kills, err := repo.GetKillsForPlayerByRoundID(ctx, kill.KillerID, kill.RoundID)
	if err != nil {
		log.Printf("error: falha na query: %v", err)
		t.Fatalf("falha na query: %v", err)
	}

	if len(kills) != 1 {
		log.Printf("error: esperado 1 kill, recebido %d", len(kills))
		t.Fatalf("esperado 1 kill, recebido %d", len(kills))
	}

	log.Println("info: teste finalizado com sucesso - obter kills por rodada")
}

func TestKillsRepository_GetTop10KillsAndDeathByPlayerID(t *testing.T) {
	log.Println("info: iniciando teste - obter top 10 kills por jogador")

	db := newTestDB(t)
	repo := NewKillsRepository(db)
	ctx := context.Background()

	for i := 0; i < 15; i++ {
		k := newTestKill()
		k.KillerID = uuid.MustParse("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa")
		repo.Save(ctx, k)
	}

	kills, err := repo.GetTop10KillsAndDeathByPlayerID(ctx, uuid.MustParse("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa"))
	if err != nil {
		log.Printf("error: falha na query: %v", err)
		t.Fatalf("falha na query: %v", err)
	}

	if len(kills) != 10 {
		log.Printf("error: esperado 10 kills, recebido %d", len(kills))
		t.Fatalf("esperado 10 kills, recebido %d", len(kills))
	}

	log.Println("info: teste finalizado com sucesso - obter top 10 kills por jogador")
}
