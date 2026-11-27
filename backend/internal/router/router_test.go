package router_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/victor4710hugo-rgb/RDO/backend/internal/database"
	"github.com/victor4710hugo-rgb/RDO/backend/internal/models"
	"github.com/victor4710hugo-rgb/RDO/backend/internal/repositories"
	"github.com/victor4710hugo-rgb/RDO/backend/internal/router"
	"github.com/victor4710hugo-rgb/RDO/backend/internal/services"
)

func setupTestRouter(t *testing.T) (*http.ServeMux, string, string) {
	t.Helper()
	t.Setenv("JWT_SECRET", "segredo-de-teste-para-validacao-rbac")

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

	repo := repositories.New(db)
	userSvc := services.NewUserService(repo)

	if _, err := userSvc.Create(models.CreateUserInput{
		Name:     "Admin Teste",
		Email:    "admin@test.local",
		Password: "password123",
		Role:     "admin",
	}); err != nil {
		t.Fatalf("criar admin de teste: %v", err)
	}

	if _, err := userSvc.Create(models.CreateUserInput{
		Name:     "Gerente Teste",
		Email:    "gerente@test.local",
		Password: "password123",
		Role:     "gerente",
	}); err != nil {
		t.Fatalf("criar gerente de teste: %v", err)
	}

	mux, err := router.New(repo)
	if err != nil {
		t.Fatalf("criar router: %v", err)
	}

	adminToken := loginForToken(t, mux, "admin@test.local", "password123")
	gerenteToken := loginForToken(t, mux, "gerente@test.local", "password123")

	return mux, adminToken, gerenteToken
}

func loginForToken(t *testing.T, mux *http.ServeMux, email, password string) string {
	t.Helper()
	body, _ := json.Marshal(map[string]string{
		"email":    email,
		"password": password,
	})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("login falhou para %s: status %d (%s)", email, rec.Code, rec.Body.String())
	}

	var resp struct {
		Token string `json:"token"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil || resp.Token == "" {
		t.Fatalf("token ausente na resposta de login para %s", email)
	}
	return resp.Token
}

func TestRouterAuthAndRBAC(t *testing.T) {
	mux, adminToken, gerenteToken := setupTestRouter(t)

	tests := []struct {
		name       string
		method     string
		path       string
		token      string
		wantStatus int
	}{
		{
			name:       "health publica sem token retorna 200",
			method:     http.MethodGet,
			path:       "/api/v1/health",
			wantStatus: http.StatusOK,
		},
		{
			name:       "users sem token retorna 401",
			method:     http.MethodGet,
			path:       "/api/v1/users",
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "projects sem token retorna 401",
			method:     http.MethodGet,
			path:       "/api/v1/projects",
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "sync push sem token retorna 401",
			method:     http.MethodPost,
			path:       "/api/v1/sync/push",
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "users com token invalido retorna 401",
			method:     http.MethodGet,
			path:       "/api/v1/users",
			token:      "token-invalido",
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "gerente em users retorna 403",
			method:     http.MethodGet,
			path:       "/api/v1/users",
			token:      gerenteToken,
			wantStatus: http.StatusForbidden,
		},
		{
			name:       "gerente em projects retorna 200",
			method:     http.MethodGet,
			path:       "/api/v1/projects",
			token:      gerenteToken,
			wantStatus: http.StatusOK,
		},
		{
			name:       "admin em users retorna 200",
			method:     http.MethodGet,
			path:       "/api/v1/users",
			token:      adminToken,
			wantStatus: http.StatusOK,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(tc.method, tc.path, nil)
			if tc.token != "" {
				req.Header.Set("Authorization", "Bearer "+tc.token)
			}
			rec := httptest.NewRecorder()
			mux.ServeHTTP(rec, req)

			if rec.Code != tc.wantStatus {
				t.Fatalf("esperado status %d, obtido %d (body: %s)", tc.wantStatus, rec.Code, rec.Body.String())
			}
		})
	}

	t.Run("gerente em sync push passa pela autorizacao RBAC", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/v1/sync/push", bytes.NewReader([]byte(`{}`)))
		req.Header.Set("Authorization", "Bearer "+gerenteToken)
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		if rec.Code == http.StatusUnauthorized || rec.Code == http.StatusForbidden {
			t.Fatalf("gerente deveria ter acesso a /api/v1/sync/push, mas recebeu %d", rec.Code)
		}
	})
}
