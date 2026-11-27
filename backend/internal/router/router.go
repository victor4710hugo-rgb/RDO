package router

import (
	"net/http"

	"github.com/victor4710hugo-rgb/RDO/backend/internal/handlers"
	"github.com/victor4710hugo-rgb/RDO/backend/internal/repositories"
	"github.com/victor4710hugo-rgb/RDO/backend/internal/services"
)

// New monta o roteador da API.
// Rotas públicas: GET /api/v1/health e POST /api/v1/auth/login.
// Todas as demais exigem JWT válido (auth.RequireAuth).
func New(repo *repositories.Repository) (*http.ServeMux, error) {
	mux := http.NewServeMux()

	authSvc, err := services.NewAuthService(repo)
	if err != nil {
		return nil, err
	}

	auth := handlers.NewAuthHandler(authSvc)
	operations := handlers.NewOperationsHandler(repo)
	fieldReports := handlers.NewFieldReportsHandler(services.NewFieldReportService(repo))
	users := handlers.NewUsersHandler(services.NewUserService(repo))
	syncH := handlers.NewSyncHandler() // evita sombrear o pacote "sync" da stdlib

	// ---------------------------------------------------------------
	// Rotas públicas
	// ---------------------------------------------------------------
	mux.HandleFunc("GET /api/v1/health", handlers.Health)
	mux.HandleFunc("POST /api/v1/auth/login", auth.Login)

	// ---------------------------------------------------------------
	// Rotas protegidas por JWT
	// ---------------------------------------------------------------

	// Autenticação
	mux.HandleFunc("GET /api/v1/auth/me", auth.RequireAuth(auth.Me))

	// Operações
	mux.HandleFunc("GET /api/v1/projects", auth.RequireAuth(operations.Projects))
	mux.HandleFunc("GET /api/v1/work-orders", auth.RequireAuth(operations.WorkOrders))
	mux.HandleFunc("GET /api/v1/activities", auth.RequireAuth(operations.Activities))

	// Relatórios de campo
	mux.HandleFunc("GET /api/v1/field-reports", auth.RequireAuth(fieldReports.List))
	mux.HandleFunc("POST /api/v1/field-reports", auth.RequireAuth(fieldReports.Upsert))

	// Usuários
	mux.HandleFunc("GET /api/v1/users", auth.RequireAuth(users.List))
	mux.HandleFunc("GET /api/v1/users/{id}", auth.RequireAuth(users.GetByID))
	mux.HandleFunc("POST /api/v1/users", auth.RequireAuth(users.Create))
	mux.HandleFunc("PUT /api/v1/users/{id}", auth.RequireAuth(users.Update))
	mux.HandleFunc("DELETE /api/v1/users/{id}", auth.RequireAuth(users.Delete))

	// Sincronização
	mux.HandleFunc("GET /api/v1/sync/pull", auth.RequireAuth(syncH.Pull))
	mux.HandleFunc("POST /api/v1/sync/push", auth.RequireAuth(syncH.Push))

	return mux, nil
}