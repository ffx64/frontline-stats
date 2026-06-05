package repositories_test

import (
	"context"
	"log"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/ffx64/frontline-stats/internal/entities"
	"github.com/ffx64/frontline-stats/internal/repositories"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupMockStatsDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock, func()) {
	log.Printf("info: configurando mock do banco de dados para testes de stats")
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Printf("error: falha ao criar mock do banco de dados: %v", err)
		t.Fatalf("erro ao criar mock DB: %v", err)
	}

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	if err != nil {
		log.Printf("error: falha ao abrir gorm mock: %v", err)
		t.Fatalf("erro ao abrir gorm mock: %v", err)
	}

	closeFunc := func() { db.Close() }
	log.Printf("info: mock do banco de dados configurado com sucesso")
	return gormDB, mock, closeFunc
}

func TestStatsRepository_Save(t *testing.T) {
	log.Printf("info: iniciando teste - salvar stats do jogador")
	db, mock, closeFunc := setupMockStatsDB(t)
	defer closeFunc()

	repo := repositories.NewPlayersStatsRepository(db)
	stats := &entities.PlayersStats{}

	mock.ExpectBegin()
	mock.ExpectExec(`INSERT INTO "players_stats"`).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.Save(context.Background(), stats)
	if err != nil {
		log.Printf("error: falha ao salvar stats do jogador: %v", err)
		t.Fatalf("erro ao salvar stats: %v", err)
	} else {
		log.Printf("info: stats do jogador salvas com sucesso")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		log.Printf("error: expectativas do mock não foram atendidas: %v", err)
		t.Errorf("expectations não foram atendidas: %v", err)
	} else {
		log.Printf("info: todas as expectativas do mock foram atendidas")
	}
}

func TestStatsRepository_UpdatePlayerStats(t *testing.T) {
	log.Printf("info: iniciando teste - atualizar stats dos jogadores")
	db, mock, closeFunc := setupMockStatsDB(t)
	defer closeFunc()

	repo := repositories.NewPlayersStatsRepository(db)

	mock.ExpectExec("UPDATE players_stats").
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.UpdatePlayerStats(context.Background())
	if err != nil {
		log.Printf("error: falha ao atualizar stats dos jogadores: %v", err)
		t.Fatalf("erro ao atualizar stats: %v", err)
	} else {
		log.Printf("info: stats dos jogadores atualizadas com sucesso")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		log.Printf("error: expectativas do mock não foram atendidas: %v", err)
		t.Errorf("expectations não foram atendidas: %v", err)
	} else {
		log.Printf("info: todas as expectativas do mock foram atendidas")
	}
}
