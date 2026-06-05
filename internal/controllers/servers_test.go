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

	"github.com/ffx64/frontline-stats/internal/controllers"
	"github.com/ffx64/frontline-stats/internal/database"
	"github.com/ffx64/frontline-stats/internal/dtos"
	"github.com/ffx64/frontline-stats/internal/entities"
	"github.com/ffx64/frontline-stats/internal/repositories"
	"github.com/ffx64/frontline-stats/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func setupServersController() *controllers.ServersController {
	db := database.Connect()
	_ = db.Migrator().DropTable(&entities.Servers{})
	_ = db.AutoMigrate(&entities.Servers{})
	repo := repositories.NewServersRepository(db)
	service := services.NewServersService(repo)
	return controllers.NewServersController(service)
}

func TestServersController_SaveServer(t *testing.T) {
	gin.SetMode(gin.TestMode)
	controller := setupServersController()
	router := gin.Default()
	router.POST("/api/v1/servers", controller.SaveServer)

	log.Println("info: iniciando teste - salvar servidor")

	body := dtos.ServersSaveDTO{Name: "ServidorDeTeste"}
	jsonBody, _ := json.Marshal(body)
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/servers", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		log.Printf("error: esperado status 201, recebido %d", w.Code)
		t.Fatalf("esperava 201, recebeu %d", w.Code)
	}

	var resp dtos.ServersDTO
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		log.Println("error: falha ao fazer unmarshal da resposta:", err)
		t.Fatal(err)
	}

	if resp.Name != "ServidorDeTeste" {
		log.Printf("error: esperado nome 'ServidorDeTeste', recebido '%s'", resp.Name)
		t.Fatalf("nome incorreto")
	}

	log.Println("info: teste finalizado com sucesso - salvar servidor")
}

func TestServersController_GetServerByID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	controller := setupServersController()
	router := gin.Default()
	router.GET("/api/v1/servers/:id", controller.GetServerByID)

	log.Println("info: iniciando teste - buscar servidor por id")

	db := database.Connect()
	repo := repositories.NewServersRepository(db)
	ctx := context.Background()
	server := &entities.Servers{ID: uuid.New(), Name: "BuscarServer", CreatedAt: time.Now()}
	repo.Save(ctx, server)

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/servers/"+server.ID.String(), nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		log.Printf("error: esperado status 200, recebido %d", w.Code)
		t.Fatalf("esperava 200, recebeu %d", w.Code)
	}

	var resp dtos.ServersDTO
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		log.Println("error: falha ao fazer unmarshal da resposta:", err)
		t.Fatal(err)
	}

	if resp.Name != "BuscarServer" {
		log.Printf("error: esperado nome 'BuscarServer', recebido '%s'", resp.Name)
		t.Fatalf("nome incorreto")
	}

	log.Println("info: teste finalizado com sucesso - buscar servidor por id")
}

func TestServersController_GetAllServers(t *testing.T) {
	gin.SetMode(gin.TestMode)
	controller := setupServersController()
	router := gin.Default()
	router.GET("/api/v1/servers", controller.GetAllServers)

	log.Println("info: iniciando teste - obter todos os servidores")

	db := database.Connect()
	repo := repositories.NewServersRepository(db)
	ctx := context.Background()
	repo.Save(ctx, &entities.Servers{ID: uuid.New(), Name: "SrvA", CreatedAt: time.Now()})
	repo.Save(ctx, &entities.Servers{ID: uuid.New(), Name: "SrvB", CreatedAt: time.Now()})

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/servers", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		log.Printf("error: esperado status 200, recebido %d", w.Code)
		t.Fatalf("esperava 200, recebeu %d", w.Code)
	}

	var resp []dtos.ServersDTO
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		log.Println("error: falha ao fazer unmarshal da resposta:", err)
		t.Fatal(err)
	}

	if len(resp) != 2 {
		log.Printf("error: esperado 2 servidores, recebidos %d", len(resp))
		t.Fatalf("quantidade incorreta")
	}

	log.Println("info: teste finalizado com sucesso - obter todos os servidores")
}

func TestServersController_UpdateServer(t *testing.T) {
	gin.SetMode(gin.TestMode)
	controller := setupServersController()
	router := gin.Default()
	router.PUT("/api/v1/servers/:id", controller.UpdateServer)

	log.Println("info: iniciando teste - atualizar servidor")

	db := database.Connect()
	repo := repositories.NewServersRepository(db)
	ctx := context.Background()
	server := &entities.Servers{ID: uuid.New(), Name: "VelhoServer", CreatedAt: time.Now()}
	repo.Save(ctx, server)

	updateBody := dtos.ServersSaveDTO{Name: "NovoServer"}
	jsonBody, _ := json.Marshal(updateBody)
	req, _ := http.NewRequest(http.MethodPut, "/api/v1/servers/"+server.ID.String(), bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		log.Printf("error: esperado status 200, recebido %d", w.Code)
		t.Fatalf("esperava 200, recebeu %d", w.Code)
	}

	var resp dtos.ServersDTO
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		log.Println("error: falha ao fazer unmarshal da resposta:", err)
		t.Fatal(err)
	}

	if resp.Name != "NovoServer" {
		log.Printf("error: esperado nome 'NovoServer', recebido '%s'", resp.Name)
		t.Fatalf("nome incorreto")
	}

	log.Println("info: teste finalizado com sucesso - atualizar servidor")
}

func TestServersController_DeleteServer(t *testing.T) {
	gin.SetMode(gin.TestMode)
	controller := setupServersController()
	router := gin.Default()
	router.DELETE("/api/v1/servers/:id", controller.DeleteServer)

	log.Println("info: iniciando teste - deletar servidor")

	db := database.Connect()
	repo := repositories.NewServersRepository(db)
	ctx := context.Background()
	server := &entities.Servers{ID: uuid.New(), Name: "DeleteMe", CreatedAt: time.Now()}
	repo.Save(ctx, server)

	req, _ := http.NewRequest(http.MethodDelete, "/api/v1/servers/"+server.ID.String(), nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		log.Printf("error: esperado status 200, recebido %d", w.Code)
		t.Fatalf("esperava 200, recebeu %d", w.Code)
	}
	log.Println("info: servidor deletado com sucesso")

	req2, _ := http.NewRequest(http.MethodDelete, "/api/v1/servers/"+server.ID.String(), nil)
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req2)

	if w2.Code != http.StatusNotFound {
		log.Printf("error: esperado 404, recebido %d", w2.Code)
		t.Fatalf("esperava 404, recebeu %d", w2.Code)
	}
	log.Println("info: tentativa de deletar servidor inexistente retornou 404 conforme esperado")

	req3, _ := http.NewRequest(http.MethodDelete, "/api/v1/servers/uuid-invalido", nil)
	w3 := httptest.NewRecorder()
	router.ServeHTTP(w3, req3)

	if w3.Code != http.StatusBadRequest {
		log.Printf("error: esperado 400, recebido %d", w3.Code)
		t.Fatalf("esperava 400, recebeu %d", w3.Code)
	}

	log.Println("info: teste finalizado com sucesso - deletar servidor")
}
