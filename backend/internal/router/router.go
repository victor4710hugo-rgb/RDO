package router

import (
	"net/http"

	"github.com/victor4710hugo-rgb/RDO/backend/internal/handlers"
	"github.com/victor4710hugo-rgb/RDO/backend/internal/repositories"
	"github.com/victor4710hugo-rgb/RDO/backend/internal/services"
)

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
	sync := handlers.NewSyncHandler()

	onlyAdmin := auth.RequireRole("admin")
	writeOps := auth.RequireRole("admin", "gerente")

	// Rotas públicas
	mux.HandleFunc("GET /api/v1/health", handlers.Health)
	mux.HandleFunc("POST /api/v1/auth/login", auth.Login)

	// Sessão (qualquer usuário autenticado)
	mux.HandleFunc("GET /api/v1/auth/me", auth.RequireAuth(auth.Me))

	// Gestão de usuários (somente admin)
	mux.HandleFunc("GET /api/v1/users", onlyAdmin(users.List))
	mux.HandleFunc("GET /api/v1/users/{id}", onlyAdmin(users.GetByID))
	mux.HandleFunc("POST /api/v1/users", onlyAdmin(users.Create))
	mux.HandleFunc("PUT /api/v1/users/{id}", onlyAdmin(users.Update))
	mux.HandleFunc("DELETE /api/v1/users/{id}", onlyAdmin(users.Delete))

	// Leitura operacional (qualquer usuário autenticado)
	mux.HandleFunc("GET /api/v1/projects", auth.RequireAuth(operations.Projects))
	mux.HandleFunc("GET /api/v1/work-orders", auth.RequireAuth(operations.WorkOrders))
	mux.HandleFunc("GET /api/v1/activities", auth.RequireAuth(operations.Activities))
	mux.HandleFunc("GET /api/v1/field-reports", auth.RequireAuth(fieldReports.List))
	mux.HandleFunc("GET /api/v1/sync/pull", auth.RequireAuth(sync.Pull))

	// Escrita operacional (admin e gerente)
	mux.HandleFunc("POST /api/v1/field-reports", writeOps(fieldReports.Upsert))
	mux.HandleFunc("POST /api/v1/sync/push", writeOps(sync.Push))

	return mux, nil
}
