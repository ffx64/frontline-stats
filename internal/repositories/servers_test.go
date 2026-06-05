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

func setupServersDB(t *testing.T) *repositories.ServersRepository {
	log.Printf("info: conectando ao banco de dados para testes de servers")
	db := database.Connect()

	log.Printf("info: limpando tabela servers antes dos testes")
	if err := db.Migrator().DropTable(&entities.Servers{}); err != nil {
		log.Printf("error: falha ao dropar tabela servers: %v", err)
		t.Fatalf("erro ao dropar tabela rounds: %v", err)
	}

	log.Printf("info: migrando entidade servers")
	if err := db.AutoMigrate(&entities.Servers{}); err != nil {
		log.Printf("error: falha ao migrar entidade servers: %v", err)
		t.Fatalf("erro ao migrar entidade servers: %v", err)
	}

	repo := repositories.NewServersRepository(db)
	log.Printf("info: repositório de servers inicializado com sucesso")
	return &repo
}

func TestServersRepository_Save(t *testing.T) {
	log.Printf("info: iniciando teste - criação de servidor")
	repo := *setupServersDB(t)

	server := &entities.Servers{Name: "SAS - Test Server"}

	saved, err := repo.Save(context.Background(), server)
	if err != nil {
		log.Printf("error: falha ao salvar servidor: %v", err)
		t.Fatalf("erro ao salvar servidor: %v", err)
	}
	log.Printf("info: servidor salvo com sucesso: id=%v", saved.ID)

	if saved.ID == uuid.Nil {
		log.Printf("error: id nulo após salvar servidor")
		t.Errorf("id não deveria ser nulo após salvar")
	}

	if saved.Name != "SAS - Test Server" {
		log.Printf("error: nome incorreto ao salvar servidor, esperado 'SAS - Test Server', obtido '%s'", saved.Name)
		t.Errorf("esperado nome 'SAS - Test Server', obtido '%s'", saved.Name)
	}
}

func TestServersRepository_FindByID(t *testing.T) {
	log.Printf("info: iniciando teste - busca de servidor por id")
	repo := *setupServersDB(t)

	server := &entities.Servers{Name: "FindByID Server"}
	saved, _ := repo.Save(context.Background(), server)
	log.Printf("info: servidor salvo para teste de busca: id=%v", saved.ID)

	found, err := repo.FindByID(context.Background(), saved.ID)
	if err != nil {
		log.Printf("error: falha ao buscar servidor: %v", err)
		t.Fatalf("erro ao buscar servidor: %v", err)
	}

	if found == nil {
		log.Printf("error: servidor não encontrado no banco de dados")
		t.Fatalf("esperava encontrar servidor, obteve nil")
	}

	if found.ID != saved.ID {
		log.Printf("error: ids divergentes, esperado %v, obtido %v", saved.ID, found.ID)
		t.Errorf("ids diferentes: esperado %v, obtido %v", saved.ID, found.ID)
	}
}

func TestServersRepository_FindAll(t *testing.T) {
	log.Printf("info: iniciando teste - listagem de servidores")
	repo := *setupServersDB(t)

	servers := []entities.Servers{
		{Name: "Server 1"},
		{Name: "Server 2"},
	}

	for _, s := range servers {
		if _, err := repo.Save(context.Background(), &s); err != nil {
			log.Printf("error: falha ao salvar servidor '%s': %v", s.Name, err)
			t.Fatalf("erro ao salvar servidor: %v", err)
		}
		log.Printf("info: servidor '%s' salvo com sucesso", s.Name)
	}

	all, err := repo.FindAll(context.Background())
	if err != nil {
		log.Printf("error: falha ao listar servidores: %v", err)
		t.Fatalf("erro ao listar servidores: %v", err)
	}

	log.Printf("info: total de servidores retornados: %d", len(all))
	if len(all) != 2 {
		log.Printf("error: quantidade incorreta de servidores, esperado 2, obtido %d", len(all))
		t.Errorf("esperava 2 servidores, obteve %d", len(all))
	}
}

func TestServersRepository_Update(t *testing.T) {
	log.Printf("info: iniciando teste - atualização de servidor")
	repo := *setupServersDB(t)

	server := &entities.Servers{Name: "Server Antigo"}
	saved, _ := repo.Save(context.Background(), server)
	log.Printf("info: servidor salvo para atualização: id=%v", saved.ID)

	saved.Name = "Server Atualizado"
	updated, err := repo.Update(context.Background(), saved)
	if err != nil {
		log.Printf("error: falha ao atualizar servidor: %v", err)
		t.Fatalf("erro ao atualizar servidor: %v", err)
	}

	if updated.Name != "Server Atualizado" {
		log.Printf("error: nome não foi atualizado corretamente, obtido '%s'", updated.Name)
		t.Errorf("esperava nome atualizado, obteve %s", updated.Name)
	}
}

func TestServersRepository_Delete(t *testing.T) {
	log.Printf("info: iniciando teste - exclusão de servidor")
	repo := *setupServersDB(t)

	server := &entities.Servers{Name: "Server Pra Deletar"}
	saved, _ := repo.Save(context.Background(), server)
	log.Printf("info: servidor salvo para exclusão: id=%v", saved.ID)

	err := repo.Delete(context.Background(), saved.ID)
	if err != nil {
		log.Printf("error: falha ao deletar servidor: %v", err)
		t.Fatalf("erro ao deletar servidor: %v", err)
	}

	found, _ := repo.FindByID(context.Background(), saved.ID)
	if found != nil {
		log.Printf("error: servidor ainda existe após tentativa de exclusão: id=%v", saved.ID)
		t.Errorf("esperava servidor deletado, mas ele ainda existe")
	} else {
		log.Printf("info: servidor deletado com sucesso: id=%v", saved.ID)
	}
}
