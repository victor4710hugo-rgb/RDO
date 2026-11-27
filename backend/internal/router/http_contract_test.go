package router_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/victor4710hugo-rgb/RDO/backend/internal/database"
	"github.com/victor4710hugo-rgb/RDO/backend/internal/models"
	"github.com/victor4710hugo-rgb/RDO/backend/internal/repositories"
	"github.com/victor4710hugo-rgb/RDO/backend/internal/router"
	"github.com/victor4710hugo-rgb/RDO/backend/internal/services"
)

// setupContractRouter inicializa banco temporário com seed (proj-001) e admin autenticado.
func setupContractRouter(t *testing.T) (*http.ServeMux, string) {
	t.Helper()
	t.Setenv("JWT_SECRET", "segredo-de-teste-contrato-http")

	origDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("obter diretório atual: %v", err)
	}
	if err := os.Chdir(t.TempDir()); err != nil {
		t.Fatalf("trocar para diretório temporário: %v", err)
	}
	t.Cleanup(func() { _ = os.Chdir(origDir) })

	db, err := database.Open()
	if err != nil {
		t.Fatalf("abrir banco de teste: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })

	// Executa o seed operacional para garantir o projeto proj-001 nos testes de RDO.
	if err := database.Seed(db); err != nil {
		t.Fatalf("executar seed: %v", err)
	}

	repo := repositories.New(db)
	userSvc := services.NewUserService(repo)
	if _, err := userSvc.Create(models.CreateUserInput{
		Name:     "Admin Contrato",
		Email:    "admin.http@test.local",
		Password: "password123",
		Role:     "admin",
	}); err != nil {
		t.Fatalf("criar admin: %v", err)
	}

	mux, err := router.New(repo)
	if err != nil {
		t.Fatalf("criar router: %v", err)
	}

	return mux, loginForToken(t, mux, "admin.http@test.local", "password123")
}

// assertJSONError garante status esperado, Content-Type application/json e campo "erro" preenchido.
func assertJSONError(t *testing.T, rec *httptest.ResponseRecorder, wantStatus int) {
	t.Helper()
	if rec.Code != wantStatus {
		t.Fatalf("esperado status %d, obtido %d (body: %s)", wantStatus, rec.Code, rec.Body.String())
	}
	if ct := rec.Header().Get("Content-Type"); !strings.Contains(ct, "application/json") {
		t.Fatalf("esperado Content-Type application/json, obtido %q", ct)
	}
	var payload map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &payload); err != nil {
		t.Fatalf("resposta de erro não é JSON válido: %v (%s)", err, rec.Body.String())
	}
	if strings.TrimSpace(payload["erro"]) == "" {
		t.Fatalf("esperado campo 'erro' não vazio na resposta JSON: %s", rec.Body.String())
	}
}

// TestHTTPStandardization valida status HTTP e contrato {"erro": "..."} em usuários e RDO.
func TestHTTPStandardization(t *testing.T) {
	mux, adminToken := setupContractRouter(t)

	// 1. Usuário inexistente deve retornar 404 com {"erro": "..."}.
	t.Run("GET user inexistente retorna 404 JSON", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/users/nao-existe", nil)
		req.Header.Set("Authorization", "Bearer "+adminToken)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)
		assertJSONError(t, rec, http.StatusNotFound)
	})

	// 2. RDO sem campos obrigatórios deve retornar 400 (e não 409).
	t.Run("POST field-report invalido retorna 400 JSON", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/v1/field-reports", bytes.NewReader([]byte(`{"id":"rdo-1"}`)))
		req.Header.Set("Authorization", "Bearer "+adminToken)
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)
		assertJSONError(t, rec, http.StatusBadRequest)
	})

	// 3. Criação de RDO novo deve retornar 201 Created com version = 1.
	t.Run("POST field-report novo retorna 201 Created", func(t *testing.T) {
		body := []byte(`{"id":"rdo-http-1","project_id":"proj-001","report_number":"RDO-001","report_date":"2026-11-27","payload":"{\"clima\":\"sol\"}"}`)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/field-reports", bytes.NewReader(body))
		req.Header.Set("Authorization", "Bearer "+adminToken)
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)
		if rec.Code != http.StatusCreated {
			t.Fatalf("esperado 201 Created, obtido %d (body: %s)", rec.Code, rec.Body.String())
		}
	})

	// 4. Atualização de RDO com versão correta (1) deve retornar 200 OK (passa para v2).
	t.Run("POST field-report atualizacao valida retorna 200 OK", func(t *testing.T) {
		body := []byte(`{"id":"rdo-http-1","project_id":"proj-001","report_number":"RDO-001","report_date":"2026-11-27","version":1,"payload":"{\"clima\":\"nublado\"}"}`)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/field-reports", bytes.NewReader(body))
		req.Header.Set("Authorization", "Bearer "+adminToken)
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)
		if rec.Code != http.StatusOK {
			t.Fatalf("esperado 200 OK, obtido %d (body: %s)", rec.Code, rec.Body.String())
		}
	})

	// 5. Atualização de RDO com versão defasada (1 quando atual é 2) deve retornar 409 Conflict JSON.
	t.Run("POST field-report com conflito de versao retorna 409 JSON", func(t *testing.T) {
		body := []byte(`{"id":"rdo-http-1","project_id":"proj-001","report_number":"RDO-001","report_date":"2026-11-27","version":1,"payload":"{\"clima\":\"chuva\"}"}`)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/field-reports", bytes.NewReader(body))
		req.Header.Set("Authorization", "Bearer "+adminToken)
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)
		assertJSONError(t, rec, http.StatusConflict)
	})
}
