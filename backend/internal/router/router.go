package router

import (
	"net/http"

	"github.com/victor4710hugo-rgb/RDO/backend/internal/handlers"
	"github.com/victor4710hugo-rgb/RDO/backend/internal/repositories"
	"github.com/victor4710hugo-rgb/RDO/backend/internal/services"
)

func New(repo *repositories.Repository) *http.ServeMux {
	mux := http.NewServeMux()

	operations := handlers.NewOperationsHandler(repo)
	fieldReports := handlers.NewFieldReportsHandler(services.NewFieldReportService(repo))
	sync := handlers.NewSyncHandler()

	mux.HandleFunc("GET /api/v1/health", handlers.Health)
	mux.HandleFunc("GET /api/v1/projects", operations.Projects)
	mux.HandleFunc("GET /api/v1/work-orders", operations.WorkOrders)
	mux.HandleFunc("GET /api/v1/activities", operations.Activities)
	mux.HandleFunc("GET /api/v1/field-reports", fieldReports.List)
	mux.HandleFunc("POST /api/v1/field-reports", fieldReports.Upsert)
	mux.HandleFunc("GET /api/v1/sync/pull", sync.Pull)
	mux.HandleFunc("POST /api/v1/sync/push", sync.Push)

	return mux
}