package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/victor4710hugo-rgb/RDO/backend/internal/models"
	"github.com/victor4710hugo-rgb/RDO/backend/internal/services"
)

type FieldReportsHandler struct {
	service *services.FieldReportService
}

func NewFieldReportsHandler(service *services.FieldReportService) *FieldReportsHandler {
	return &FieldReportsHandler{service: service}
}

func (h *FieldReportsHandler) List(w http.ResponseWriter, r *http.Request) {
	reports, err := h.service.List()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, reports)
}

func (h *FieldReportsHandler) Upsert(w http.ResponseWriter, r *http.Request) {
	var report models.FieldReport
	if err := json.NewDecoder(r.Body).Decode(&report); err != nil {
		writeError(w, http.StatusBadRequest, "JSON inválido")
		return
	}

	saved, err := h.service.Upsert(report)
	if err != nil {
		writeError(w, http.StatusConflict, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, saved)
}
