package controllers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ffx64/gamestats-backend/internal/controllers"
	"github.com/ffx64/gamestats-backend/internal/dtos"
	appErrors "github.com/ffx64/gamestats-backend/internal/errors"
	"github.com/gin-gonic/gin"
)

type fakeKillsService struct {
	saveFunc func(ctx context.Context, dto []dtos.KillsSaveDTO) error
}

func (f *fakeKillsService) SaveKills(ctx context.Context, dto []dtos.KillsSaveDTO) error {
	if f.saveFunc != nil {
		return f.saveFunc(ctx, dto)
	}
	return nil
}

func TestKillsController_SaveKills(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("deve retornar 400 quando o json for inválido", func(t *testing.T) {
		log.Println("info: iniciando teste - json inválido")

		controller := controllers.NewKillsController(&fakeKillsService{})
		router := gin.Default()
		router.POST("/kills", controller.SaveKills)

		req, _ := http.NewRequest(http.MethodPost, "/kills", bytes.NewBuffer([]byte(`invalid-json`)))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		if resp.Code != http.StatusBadRequest {
			log.Printf("error: esperado 400, recebido %d", resp.Code)
			t.Fatalf("esperava 400, recebeu %d", resp.Code)
		}

		log.Println("info: teste finalizado com sucesso - json inválido")
	})

	t.Run("deve retornar erro customizado vindo do service", func(t *testing.T) {
		log.Println("info: iniciando teste - erro customizado do service")

		fakeService := &fakeKillsService{
			saveFunc: func(ctx context.Context, dto []dtos.KillsSaveDTO) error {
				return appErrors.ErrServerNotFound
			},
		}
		controller := controllers.NewKillsController(fakeService)
		router := gin.Default()
		router.POST("/kills", controller.SaveKills)

		body, _ := json.Marshal([]dtos.KillsSaveDTO{{ServerID: "id"}})
		req, _ := http.NewRequest(http.MethodPost, "/kills", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		if resp.Code != appErrors.ErrServerNotFound.Status {
			log.Printf("error: esperado %d, recebido %d", appErrors.ErrServerNotFound.Status, resp.Code)
			t.Fatalf("esperava %d, recebeu %d", appErrors.ErrServerNotFound.Status, resp.Code)
		}

		log.Println("info: teste finalizado com sucesso - erro customizado do service")
	})

	t.Run("deve retornar 500 em erro genérico do service", func(t *testing.T) {
		log.Println("info: iniciando teste - erro genérico do service")

		fakeService := &fakeKillsService{
			saveFunc: func(ctx context.Context, dto []dtos.KillsSaveDTO) error {
				return errors.New("erro inesperado no service")
			},
		}
		controller := controllers.NewKillsController(fakeService)
		router := gin.Default()
		router.POST("/kills", controller.SaveKills)

		body, _ := json.Marshal([]dtos.KillsSaveDTO{{ServerID: "id"}})
		req, _ := http.NewRequest(http.MethodPost, "/kills", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		if resp.Code != http.StatusInternalServerError {
			log.Printf("error: esperado 500, recebido %d", resp.Code)
			t.Fatalf("esperava 500, recebeu %d", resp.Code)
		}

		log.Println("info: teste finalizado com sucesso - erro genérico do service")
	})

	t.Run("deve retornar 201 quando as kills forem salvas com sucesso", func(t *testing.T) {
		log.Println("info: iniciando teste - salvamento de kills bem-sucedido")

		fakeService := &fakeKillsService{
			saveFunc: func(ctx context.Context, dto []dtos.KillsSaveDTO) error {
				return nil
			},
		}
		controller := controllers.NewKillsController(fakeService)
		router := gin.Default()
		router.POST("/kills", controller.SaveKills)

		body, _ := json.Marshal([]dtos.KillsSaveDTO{{ServerID: "id"}})
		req, _ := http.NewRequest(http.MethodPost, "/kills", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		if resp.Code != http.StatusCreated {
			log.Printf("error: esperado 201, recebido %d", resp.Code)
			t.Fatalf("esperava 201, recebeu %d", resp.Code)
		}

		log.Println("info: teste finalizado com sucesso - salvamento de kills bem-sucedido")
	})
}
