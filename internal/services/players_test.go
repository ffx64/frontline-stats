package services_test

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/ffx64/gamestats-backend/internal/database"
	"github.com/ffx64/gamestats-backend/internal/dtos"
	"github.com/ffx64/gamestats-backend/internal/entities"
	"github.com/ffx64/gamestats-backend/internal/repositories"
	"github.com/ffx64/gamestats-backend/internal/services"
)

func setupPlayersDB(t *testing.T) *gorm.DB {
	log.Printf("info: conectando ao banco de dados em memória para testes de players service")
	db := database.Connect()

	log.Printf("info: limpando tabela players e players_stats antes dos testes")
	if err := db.Migrator().DropTable(&entities.Players{}, &entities.PlayersStats{}); err != nil {
		log.Printf("error: falha ao dropar tabela players e players_stats: %v", err)
		t.Fatalf("erro ao dropar tabela players e players_stats: %v", err)
	}

	log.Printf("info: migrando entidades players e players_stats")
	if err := db.AutoMigrate(&entities.Players{}, &entities.PlayersStats{}); err != nil {
		log.Printf("error: falha ao migrar entidades: %v", err)
		t.Fatalf("erro ao migrar entidades: %v", err)
	}
	log.Printf("info: banco de dados preparado com sucesso")
	return db
}

func TestPlayersService_SavePlayer(t *testing.T) {
	log.Printf("info: iniciando teste - salvar player")
	db := setupPlayersDB(t)

	repo := repositories.NewPlayersRepository(db)
	statsRepo := repositories.NewPlayersStatsRepository(db)
	service := services.NewPlayersService(repo, statsRepo)
	ctx := context.Background()

	dto := &dtos.PlayerSaveDTO{
		GUID:         uuid.NewString(),
		Username:     "CriarTest",
		LastServerID: uuid.NewString(),
		Platform:     "pc",
	}

	log.Printf("info: salvando jogador %s", dto.Username)
	created, err := service.Save(ctx, dto)
	if err != nil {
		log.Printf("error: falha ao criar jogador: %v", err)
		t.Fatalf("falha ao criar jogador: %v", err)
	}
	log.Printf("info: jogador criado com sucesso, GUID=%s", created.GUID)

	if created.GUID == "" {
		log.Printf("error: GUID do jogador está vazio após criação")
		t.Fatal("esperava GUID preenchido, mas veio vazio")
	}
	if created.Username != dto.Username {
		log.Printf("error: username incorreto, esperado '%s', obtido '%s'", dto.Username, created.Username)
		t.Errorf("esperava username %s, recebeu %s", dto.Username, created.Username)
	}
}

func TestPlayersService_GetPlayerByGUID(t *testing.T) {
	log.Printf("info: iniciando teste - buscar player por GUID")
	db := setupPlayersDB(t)

	repo := repositories.NewPlayersRepository(db)
	statsRepo := repositories.NewPlayersStatsRepository(db)
	service := services.NewPlayersService(repo, statsRepo)
	ctx := context.Background()

	player := &entities.Players{
		ID:        uuid.New(),
		GUID:      uuid.NewString(),
		Username:  "BuscarTest",
		Platform:  "pc",
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	log.Printf("info: salvando jogador para teste, username=%s", player.Username)
	repo.Save(ctx, player)

	log.Printf("info: buscando jogador pelo GUID=%s", player.GUID)
	found, err := service.GetByGUID(ctx, player.GUID)
	if err != nil {
		log.Printf("error: falha ao buscar jogador pelo GUID: %v", err)
		t.Fatalf("falha ao buscar jogador pelo GUID: %v", err)
	}
	log.Printf("info: jogador encontrado, username=%s", found.Username)

	if found.Username != player.Username {
		log.Printf("error: username divergente, esperado '%s', obtido '%s'", player.Username, found.Username)
		t.Errorf("esperava username %s, recebeu %s", player.Username, found.Username)
	}
}

func TestPlayersService_UpdatePlayer(t *testing.T) {
	log.Printf("info: iniciando teste - atualizar player")
	db := setupPlayersDB(t)

	repo := repositories.NewPlayersRepository(db)
	statsRepo := repositories.NewPlayersStatsRepository(db)
	service := services.NewPlayersService(repo, statsRepo)
	ctx := context.Background()

	player := &entities.Players{
		ID:        uuid.New(),
		GUID:      uuid.NewString(),
		Username:  "AntigoNome",
		Platform:  "pc",
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	log.Printf("info: salvando jogador para teste de atualização, username=%s", player.Username)
	if err := repo.Save(ctx, player); err != nil {
		log.Printf("error: falha ao criar jogador inicial: %v", err)
		t.Fatalf("falha ao criar jogador inicial: %v", err)
	}

	updateDTO := &dtos.PlayerUpdateDTO{
		Username:     "NovoNome",
		LastServerID: uuid.NewString(),
	}

	log.Printf("info: atualizando jogador GUID=%s", player.GUID)
	updated, err := service.Update(ctx, player.GUID, updateDTO)
	if err != nil {
		log.Printf("error: falha ao atualizar jogador: %v", err)
		t.Fatalf("falha ao atualizar jogador: %v", err)
	}
	log.Printf("info: jogador atualizado com sucesso, novo username=%s", updated.Username)
}
