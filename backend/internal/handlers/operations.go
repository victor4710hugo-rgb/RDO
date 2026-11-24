package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/victor4710hugo-rgb/RDO/backend/internal/repositories"
)

type OperationsHandler struct {
	repo *repositories.Repository
}

func NewOperationsHandler(repo *repositories.Repository) *OperationsHandler {
	return &OperationsHandler{repo: repo}
}

func (h *OperationsHandler) Projects(w http.ResponseWriter, r *http.Request) {
	projects, err := h.repo.ListProjects()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(projects)
}

func (h *OperationsHandler) WorkOrders(w http.ResponseWriter, r *http.Request) {
	orders, err := h.repo.ListWorkOrders()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orders)
}

func (h *OperationsHandler) Activities(w http.ResponseWriter, r *http.Request) {
	activities, err := h.repo.ListActivities()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(activities)
}