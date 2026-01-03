package services_test

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/ffx64/gamestats-backend/internal/database"
	"github.com/ffx64/gamestats-backend/internal/dtos"
	"github.com/ffx64/gamestats-backend/internal/entities"
	"github.com/ffx64/gamestats-backend/internal/repositories"
	"github.com/ffx64/gamestats-backend/internal/services"
	"github.com/google/uuid"
)

func setupServersDB(t *testing.T) *repositories.ServersRepository {
	log.Printf("info: conectando ao banco de dados em memória para testes de servers service")
	db := database.Connect()

	log.Printf("info: limpando tabela servers antes dos testes")
	if err := db.Migrator().DropTable(&entities.Servers{}); err != nil {
		log.Printf("error: falha ao dropar tabela servers: %v", err)
		t.Fatalf("erro ao dropar tabela servers: %v", err)
	}

	log.Printf("info: migrando entidades Servers")
	if err := db.AutoMigrate(&entities.Servers{}); err != nil {
		log.Printf("error: falha ao migrar entidades: %v", err)
		t.Fatalf("erro ao migrar entidades: %v", err)
	}

	repo := repositories.NewServersRepository(db)

	log.Printf("info: banco de dados preparado com sucesso")
	return &repo
}

func TestServersService_SaveServer(t *testing.T) {
	log.Printf("info: iniciando teste - salvar servidor")
	repo := *setupServersDB(t)

	service := services.NewServersService(repo)
	ctx := context.Background()

	dto := &dtos.ServersSaveDTO{Name: "ServidorTeste"}
	log.Printf("info: salvando servidor com nome=%s", dto.Name)
	created, err := service.SaveServer(ctx, dto)
	if err != nil {
		log.Printf("error: falha ao salvar servidor: %v", err)
		t.Fatalf("falha ao salvar servidor: %v", err)
	}
	log.Printf("info: servidor salvo com sucesso, ID=%s", created.ID)

	if created.Name != dto.Name {
		log.Printf("error: nome do servidor divergente, esperado %s, obtido %s", dto.Name, created.Name)
		t.Errorf("esperava nome %s, recebeu %s", dto.Name, created.Name)
	}

	// tenta criar o mesmo novamente (deve falhar)
	log.Printf("info: tentando criar servidor duplicado")
	_, err = service.SaveServer(ctx, dto)
	if err == nil {
		log.Printf("error: servidor duplicado criado sem erro")
		t.Fatal("esperava erro de servidor existente, mas não ocorreu")
	}
	log.Printf("info: servidor duplicado corretamente retornou erro")
}

func TestServersService_GetServerByID(t *testing.T) {
	log.Printf("info: iniciando teste - buscar servidor por ID")
	repo := *setupServersDB(t)

	service := services.NewServersService(repo)
	ctx := context.Background()

	server := &entities.Servers{ID: uuid.New(), Name: "BuscarServidor", CreatedAt: time.Now()}
	log.Printf("info: salvando servidor para teste, ID=%s", server.ID)
	repo.Save(ctx, server)

	log.Printf("info: buscando servidor pelo ID=%s", server.ID)
	found, err := service.GetServerByID(ctx, server.ID)
	if err != nil {
		log.Printf("error: falha ao buscar servidor: %v", err)
		t.Fatalf("falha ao buscar servidor: %v", err)
	}
	log.Printf("info: servidor encontrado com sucesso, nome=%s", found.Name)

	if found.Name != server.Name {
		log.Printf("error: nome divergente, esperado %s, obtido %s", server.Name, found.Name)
		t.Errorf("esperava nome %s, recebeu %s", server.Name, found.Name)
	}
}

func TestServersService_GetAllServers(t *testing.T) {
	log.Printf("info: iniciando teste - listar todos os servidores")
	repo := *setupServersDB(t)

	service := services.NewServersService(repo)
	ctx := context.Background()

	log.Printf("info: salvando servidores para teste de listagem")
	repo.Save(ctx, &entities.Servers{ID: uuid.New(), Name: "ServidorA", CreatedAt: time.Now()})
	repo.Save(ctx, &entities.Servers{ID: uuid.New(), Name: "ServidorB", CreatedAt: time.Now()})

	log.Printf("info: listando todos os servidores")
	list, err := service.GetAllServers(ctx)
	if err != nil {
		log.Printf("error: falha ao listar servidores: %v", err)
		t.Fatalf("falha ao listar servidores: %v", err)
	}
	log.Printf("info: total de servidores retornados: %d", len(list))

	if len(list) != 2 {
		log.Printf("error: quantidade de servidores incorreta, esperado 2, obtido %d", len(list))
		t.Errorf("esperava 2 servidores, recebeu %d", len(list))
	}
}

func TestServersService_UpdateServer(t *testing.T) {
	log.Printf("info: iniciando teste - atualizar servidor")
	repo := *setupServersDB(t)

	service := services.NewServersService(repo)
	ctx := context.Background()

	server := &entities.Servers{ID: uuid.New(), Name: "ServidorAntigo", CreatedAt: time.Now()}
	log.Printf("info: salvando servidor para teste de atualização, ID=%s", server.ID)
	repo.Save(ctx, server)

	updateDTO := &dtos.ServersSaveDTO{Name: "ServidorNovo"}
	log.Printf("info: atualizando servidor ID=%s para nome=%s", server.ID, updateDTO.Name)
	updated, err := service.UpdateServer(ctx, server.ID, updateDTO)
	if err != nil {
		log.Printf("error: falha ao atualizar servidor: %v", err)
		t.Fatalf("falha ao atualizar servidor: %v", err)
	}
	log.Printf("info: servidor atualizado com sucesso, novo nome=%s", updated.Name)

	if updated.Name != "ServidorNovo" {
		log.Printf("error: nome não atualizado, esperado 'ServidorNovo', obtido '%s'", updated.Name)
		t.Errorf("esperava nome atualizado 'ServidorNovo', recebeu '%s'", updated.Name)
	}
}

func TestServersService_DeleteServer(t *testing.T) {
	db := database.Connect()
	_ = db.Migrator().DropTable(&entities.Servers{})
	_ = db.AutoMigrate(&entities.Servers{})

	repo := repositories.NewServersRepository(db)
	service := services.NewServersService(repo)
	ctx := context.Background()

	// cria servidor pra deletar
	server := &entities.Servers{
		ID:        uuid.New(),
		Name:      "ServerPraDeletar",
		CreatedAt: time.Now(),
	}
	if _, err := repo.Save(ctx, server); err != nil {
		t.Fatalf("falha ao criar servidor de teste: %v", err)
	}

	// deleta o servidor
	ok, err := service.DeleteServer(ctx, server.ID)
	if err != nil {
		t.Fatalf("erro inesperado ao deletar servidor: %v", err)
	}
	if !ok {
		t.Errorf("esperava retorno true, recebeu false")
	}

	// tenta buscar novamente (deve não existir)
	found, _ := repo.FindByID(ctx, server.ID)
	if found != nil {
		t.Errorf("esperava servidor deletado, mas ainda existe")
	}

	// tenta deletar um id inexistente
	randomID := uuid.New()
	ok, err = service.DeleteServer(ctx, randomID)
	if err == nil {
		t.Errorf("esperava erro ao deletar servidor inexistente, mas não ocorreu")
	}
	if ok {
		t.Errorf("esperava false ao deletar servidor inexistente, recebeu true")
	}
}
