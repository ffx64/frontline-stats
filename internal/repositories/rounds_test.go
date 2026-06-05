package repositories_test

import (
	"context"
	"log"
	"testing"

	"github.com/ffx64/frontline-stats/internal/database"
	"github.com/ffx64/frontline-stats/internal/entities"
	"github.com/ffx64/frontline-stats/internal/repositories"
	"github.com/google/uuid"
)

func setupRoundsDB(t *testing.T) *repositories.RoundsRepository {
	log.Printf("info: conectando ao banco de dados para testes de rounds")
	db := database.Connect()

	log.Printf("info: limpando tabela rounds antes dos testes")
	if err := db.Migrator().DropTable(&entities.Rounds{}); err != nil {
		log.Printf("error: falha ao dropar tabela rounds: %v", err)
		t.Fatalf("erro ao dropar tabela rounds: %v", err)
	}

	log.Printf("info: migrando entidade rounds")
	if err := db.AutoMigrate(&entities.Rounds{}); err != nil {
		log.Printf("error: falha ao migrar entidade rounds: %v", err)
		t.Fatalf("erro ao migrar entidade rounds: %v", err)
	}

	repo := repositories.NewRoundsRepository(db)
	log.Printf("info: repositório de rounds inicializado com sucesso")
	return &repo
}

func TestRoundsRepository_Save(t *testing.T) {
	log.Printf("info: iniciando teste - criação de rodada")
	repo := *setupRoundsDB(t)

	round := &entities.Rounds{
		ServerID: uuid.New(),
		Status:   "IN_PROGRESS",
	}

	saved, err := repo.Save(context.Background(), round)
	if err != nil {
		log.Printf("error: falha ao criar rodada: %v", err)
		t.Fatalf("erro ao criar rodada: %v", err)
	}

	log.Printf("info: rodada criada com sucesso: id=%v", saved.ID)

	if saved.ID == uuid.Nil {
		log.Printf("error: id da rodada é nulo após salvar")
		t.Errorf("id não deveria ser nulo após salvar")
	}

	if saved.Status != "IN_PROGRESS" {
		log.Printf("error: status incorreto ao salvar rodada, esperado 'IN_PROGRESS', obtido '%s'", saved.Status)
		t.Errorf("esperava status 'IN_PROGRESS', obteve '%s'", saved.Status)
	}
}

func TestRoundsRepository_FindByID(t *testing.T) {
	log.Printf("info: iniciando teste - busca de rodada por id")
	repo := *setupRoundsDB(t)

	round := &entities.Rounds{
		ServerID: uuid.New(),
		Status:   "WAITING",
	}
	saved, _ := repo.Save(context.Background(), round)
	log.Printf("info: rodada salva para teste de busca: id=%v", saved.ID)

	found, err := repo.FindByID(context.Background(), saved.ID)
	if err != nil {
		log.Printf("error: falha ao buscar rodada por id: %v", err)
		t.Fatalf("erro ao buscar rodada: %v", err)
	}

	if found == nil {
		log.Printf("error: rodada não encontrada no banco de dados")
		t.Fatalf("esperava encontrar rodada, obteve nil")
	}

	if found.ID != saved.ID {
		log.Printf("error: ids divergentes, esperado %v, obtido %v", saved.ID, found.ID)
		t.Errorf("ids diferentes: esperado %v, obtido %v", saved.ID, found.ID)
	}
}

func TestRoundsRepository_FindAll(t *testing.T) {
	log.Printf("info: iniciando teste - listagem paginada de rodadas")

	repo := *setupRoundsDB(t)
	ctx := context.Background()

	rounds := []entities.Rounds{
		{ServerID: uuid.New(), Status: "ROUND_1"},
		{ServerID: uuid.New(), Status: "ROUND_2"},
		{ServerID: uuid.New(), Status: "ROUND_3"},
	}

	for _, r := range rounds {
		if _, err := repo.Save(ctx, &r); err != nil {
			log.Printf("error: falha ao salvar rodada '%s': %v", r.Status, err)
			t.Fatalf("erro ao salvar rodada: %v", err)
		}
		log.Printf("info: rodada '%s' salva com sucesso", r.Status)
	}

	// primeira página (limit 2)
	pageLimit := 2
	pageOffset := 0

	all, total, err := repo.FindAll(ctx, pageLimit, pageOffset)
	if err != nil {
		log.Printf("error: falha ao listar rodadas: %v", err)
		t.Fatalf("erro ao listar rodadas: %v", err)
	}

	log.Printf("info: total de rodadas encontradas: %d", total)
	log.Printf("info: rodadas retornadas nesta página: %d", len(all))

	if total != int64(len(rounds)) {
		t.Errorf("esperava total %d, obteve %d", len(rounds), total)
	}

	if len(all) != pageLimit {
		t.Errorf("esperava %d rodadas na página, obteve %d", pageLimit, len(all))
	}

	// segunda página
	pageOffset = pageLimit
	secondPage, total2, err := repo.FindAll(ctx, pageLimit, pageOffset)
	if err != nil {
		t.Fatalf("erro ao listar segunda página: %v", err)
	}

	if total2 != total {
		t.Errorf("o total de registros deveria ser o mesmo (%d), obteve %d", total, total2)
	}

	if len(secondPage) != 1 {
		t.Errorf("esperava 1 rodada na segunda página, obteve %d", len(secondPage))
	}

	log.Printf("info: teste de paginação concluído com sucesso")
}
