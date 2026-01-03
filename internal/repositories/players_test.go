package repositories_test

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/ffx64/gamestats-backend/internal/database"
	"github.com/ffx64/gamestats-backend/internal/entities"
	"github.com/ffx64/gamestats-backend/internal/repositories"
	"github.com/google/uuid"
)

func TestPlayersRepository_Create(t *testing.T) {
	log.Println("info: iniciando teste - criar jogador")

	db := database.Connect()

	if err := db.Migrator().DropTable(&entities.Players{}); err != nil {
		log.Printf("error: falha ao dropar tabela: %v", err)
		t.Fatalf("erro ao dropar tabela: %v", err)
	}

	if err := db.AutoMigrate(&entities.Players{}); err != nil {
		log.Printf("error: falha ao migrar entidade: %v", err)
		t.Fatalf("erro ao migrar entidade: %v", err)
	}

	repo := repositories.NewPlayersRepository(db)
	ctx := context.Background()

	now := time.Now()
	player := &entities.Players{
		ID:        uuid.New(),
		GUID:      uuid.NewString(),
		Username:  "JogadorTest",
		Platform:  "pc",
		IsActive:  true,
		IsBanned:  false,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := repo.Save(ctx, player); err != nil {
		log.Printf("error: falha ao criar jogador: %v", err)
		t.Fatalf("erro ao criar jogador: %v", err)
	}

	if player.ID == uuid.Nil {
		log.Println("error: id do jogador não foi gerado")
		t.Fatal("esperava ID gerado automaticamente, mas veio vazio")
	}

	log.Println("info: teste finalizado com sucesso - criar jogador")
}

func TestPlayersRepository_FindByID(t *testing.T) {
	log.Println("info: iniciando teste - buscar jogador por ID")

	db := database.Connect()

	if err := db.Migrator().DropTable(&entities.Players{}); err != nil {
		log.Printf("error: falha ao dropar tabela: %v", err)
		t.Fatalf("erro ao dropar tabela: %v", err)
	}

	if err := db.AutoMigrate(&entities.Players{}); err != nil {
		log.Printf("error: falha ao migrar entidade: %v", err)
		t.Fatalf("erro ao migrar entidade: %v", err)
	}

	repo := repositories.NewPlayersRepository(db)
	ctx := context.Background()

	now := time.Now()
	player := &entities.Players{
		ID:        uuid.New(),
		GUID:      uuid.NewString(),
		Username:  "Jogador123",
		Platform:  "pc",
		IsActive:  true,
		IsBanned:  false,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := repo.Save(ctx, player); err != nil {
		log.Printf("error: falha ao criar jogador: %v", err)
		t.Fatalf("erro ao criar jogador: %v", err)
	}

	found, err := repo.FindByID(ctx, player.ID)
	if err != nil {
		log.Printf("error: falha ao buscar jogador: %v", err)
		t.Fatalf("erro ao buscar jogador: %v", err)
	}

	if found.Username != player.Username {
		log.Printf("error: nome incorreto - esperado %s, recebido %s", player.Username, found.Username)
		t.Errorf("esperava nome %s, recebeu %s", player.Username, found.Username)
	}

	log.Println("info: teste finalizado com sucesso - buscar jogador por ID")
}

func TestPlayersRepository_FindByGUID(t *testing.T) {
	log.Println("info: iniciando teste - buscar jogador por GUID")

	db := database.Connect()

	if err := db.Migrator().DropTable(&entities.Players{}); err != nil {
		log.Printf("error: falha ao dropar tabela: %v", err)
		t.Fatalf("erro ao dropar tabela: %v", err)
	}

	if err := db.AutoMigrate(&entities.Players{}); err != nil {
		log.Printf("error: falha ao migrar entidade: %v", err)
		t.Fatalf("erro ao migrar entidade: %v", err)
	}

	repo := repositories.NewPlayersRepository(db)
	ctx := context.Background()

	now := time.Now()
	player := &entities.Players{
		ID:        uuid.New(),
		GUID:      uuid.NewString(),
		Username:  "Jogador123",
		Platform:  "pc",
		IsActive:  true,
		IsBanned:  false,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := repo.Save(ctx, player); err != nil {
		log.Printf("error: falha ao criar jogador: %v", err)
		t.Fatalf("erro ao criar jogador: %v", err)
	}

	found, err := repo.FindByGUID(ctx, player.GUID)
	if err != nil {
		log.Printf("error: falha ao buscar jogador: %v", err)
		t.Fatalf("erro ao buscar jogador: %v", err)
	}

	if found.Username != player.Username {
		log.Printf("error: nome incorreto - esperado %s, recebido %s", player.Username, found.Username)
		t.Errorf("esperava nome %s, recebeu %s", player.Username, found.Username)
	}

	log.Println("info: teste finalizado com sucesso - buscar jogador por GUID")
}

func TestPlayersRepository_FindAll(t *testing.T) {
	log.Println("info: iniciando teste - listar jogadores")

	db := database.Connect()

	if err := db.Migrator().DropTable(&entities.Players{}); err != nil {
		log.Printf("error: falha ao dropar tabela: %v", err)
		t.Fatalf("erro ao dropar tabela: %v", err)
	}

	if err := db.AutoMigrate(&entities.Players{}); err != nil {
		log.Printf("error: falha ao migrar entidade: %v", err)
		t.Fatalf("erro ao migrar entidade: %v", err)
	}

	repo := repositories.NewPlayersRepository(db)
	ctx := context.Background()

	now := time.Now()
	if err := repo.Save(ctx, &entities.Players{ID: uuid.New(), GUID: uuid.NewString(), Username: "A", Platform: "pc", IsActive: true, CreatedAt: now, UpdatedAt: now}); err != nil {
		log.Printf("error: falha ao criar jogador A: %v", err)
		t.Fatalf("erro ao criar jogador A: %v", err)
	}
	if err := repo.Save(ctx, &entities.Players{ID: uuid.New(), GUID: uuid.NewString(), Username: "B", Platform: "pc", IsActive: true, CreatedAt: now, UpdatedAt: now}); err != nil {
		log.Printf("error: falha ao criar jogador B: %v", err)
		t.Fatalf("erro ao criar jogador B: %v", err)
	}

	list, err := repo.FindAll(ctx)
	if err != nil {
		log.Printf("error: falha ao listar jogadores: %v", err)
		t.Fatalf("erro ao listar jogadores: %v", err)
	}

	if len(list) != 2 {
		log.Printf("error: esperado 2 jogadores, recebido %d", len(list))
		t.Errorf("esperava 2 jogadores, recebeu %d", len(list))
	}

	log.Println("info: teste finalizado com sucesso - listar jogadores")
}

func TestPlayersRepository_Update(t *testing.T) {
	log.Println("info: iniciando teste - atualizar jogador")

	db := database.Connect()

	if err := db.Migrator().DropTable(&entities.Players{}); err != nil {
		log.Printf("error: falha ao dropar tabela: %v", err)
		t.Fatalf("erro ao dropar tabela: %v", err)
	}

	if err := db.AutoMigrate(&entities.Players{}); err != nil {
		log.Printf("error: falha ao migrar entidade: %v", err)
		t.Fatalf("erro ao migrar entidade: %v", err)
	}

	repo := repositories.NewPlayersRepository(db)
	ctx := context.Background()

	now := time.Now()
	player := &entities.Players{
		ID:        uuid.New(),
		GUID:      uuid.NewString(),
		Username:  "AntigoNome",
		Platform:  "pc",
		IsActive:  true,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := repo.Save(ctx, player); err != nil {
		log.Printf("error: falha ao criar jogador: %v", err)
		t.Fatalf("erro ao criar jogador: %v", err)
	}

	player.Username = "NovoNome"
	player.UpdatedAt = time.Now()
	if err := repo.Update(ctx, player); err != nil {
		log.Printf("error: falha ao atualizar jogador: %v", err)
		t.Fatalf("erro ao atualizar jogador: %v", err)
	}

	updated, _ := repo.FindByID(ctx, player.ID)
	if updated.Username != "NovoNome" {
		log.Printf("error: nome não atualizado corretamente - esperado 'NovoNome', recebido '%s'", updated.Username)
		t.Errorf("esperava nome atualizado 'NovoNome', recebeu '%s'", updated.Username)
	}

	log.Println("info: teste finalizado com sucesso - atualizar jogador")
}

func TestPlayersRepository_Delete(t *testing.T) {
	log.Println("info: iniciando teste - deletar jogador")

	db := database.Connect()

	if err := db.Migrator().DropTable(&entities.Players{}); err != nil {
		log.Printf("error: falha ao dropar tabela: %v", err)
		t.Fatalf("erro ao dropar tabela: %v", err)
	}

	if err := db.AutoMigrate(&entities.Players{}); err != nil {
		log.Printf("error: falha ao migrar entidade: %v", err)
		t.Fatalf("erro ao migrar entidade: %v", err)
	}

	repo := repositories.NewPlayersRepository(db)
	ctx := context.Background()

	now := time.Now()
	player := &entities.Players{
		ID:        uuid.New(),
		GUID:      uuid.NewString(),
		Username:  "DeletarMe",
		Platform:  "pc",
		IsActive:  true,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := repo.Save(ctx, player); err != nil {
		log.Printf("error: falha ao criar jogador: %v", err)
		t.Fatalf("erro ao criar jogador: %v", err)
	}

	if err := repo.Delete(ctx, player.ID); err != nil {
		log.Printf("error: falha ao deletar jogador: %v", err)
		t.Fatalf("erro ao deletar jogador: %v", err)
	}

	_, err := repo.FindByID(ctx, player.ID)
	if err == nil {
		log.Println("error: jogador ainda encontrado após exclusão")
		t.Error("esperava erro ao buscar jogador deletado, mas não ocorreu")
	}

	log.Println("info: teste finalizado com sucesso - deletar jogador")
}
