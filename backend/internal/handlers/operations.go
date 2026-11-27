package handlers

import (
	"net/http"

	"github.com/victor4710hugo-rgb/RDO/backend/internal/repositories"
)

// OperationsHandler expõe as consultas operacionais (projetos, ordens de serviço e atividades).
type OperationsHandler struct {
	repo *repositories.Repository
}

// NewOperationsHandler cria o handler de consultas operacionais.
func NewOperationsHandler(repo *repositories.Repository) *OperationsHandler {
	return &OperationsHandler{repo: repo}
}

// Projects atende GET /api/v1/projects com resposta e erro sempre em JSON padronizado.
func (h *OperationsHandler) Projects(w http.ResponseWriter, r *http.Request) {
	projects, err := h.repo.ListProjects()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "erro interno ao listar projetos")
		return
	}
	writeJSON(w, http.StatusOK, projects)
}

// WorkOrders atende GET /api/v1/work-orders com resposta e erro sempre em JSON padronizado.
func (h *OperationsHandler) WorkOrders(w http.ResponseWriter, r *http.Request) {
	orders, err := h.repo.ListWorkOrders()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "erro interno ao listar ordens de serviço")
		return
	}
	writeJSON(w, http.StatusOK, orders)
}

// Activities atende GET /api/v1/activities com resposta e erro sempre em JSON padronizado.
func (h *OperationsHandler) Activities(w http.ResponseWriter, r *http.Request) {
	activities, err := h.repo.ListActivities()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "erro interno ao listar atividades")
		return
	}
	writeJSON(w, http.StatusOK, activities)
}
