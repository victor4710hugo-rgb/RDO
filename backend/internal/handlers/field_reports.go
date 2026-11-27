package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/victor4710hugo-rgb/RDO/backend/internal/models"
	"github.com/victor4710hugo-rgb/RDO/backend/internal/repositories"
	"github.com/victor4710hugo-rgb/RDO/backend/internal/services"
)

// FieldReportsHandler traduz requisições HTTP de /api/v1/field-reports para o FieldReportService.
type FieldReportsHandler struct {
	service *services.FieldReportService
}

// NewFieldReportsHandler cria o handler de RDO com o serviço injetado.
func NewFieldReportsHandler(service *services.FieldReportService) *FieldReportsHandler {
	return &FieldReportsHandler{service: service}
}

// mapFieldReportError converte erros de validação (400), conflito de versão (409) e falha interna (500).
func mapFieldReportError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, services.ErrValidation):
		writeError(w, http.StatusBadRequest, err.Error())
	case errors.Is(err, repositories.ErrVersionConflict):
		writeError(w, http.StatusConflict, err.Error())
	default:
		writeError(w, http.StatusInternalServerError, "erro interno ao processar RDO")
	}
}

// List atende GET /api/v1/field-reports retornando 200 OK ou 500 JSON padronizado.
func (h *FieldReportsHandler) List(w http.ResponseWriter, r *http.Request) {
	reports, err := h.service.List()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "erro interno ao listar RDOs")
		return
	}
	writeJSON(w, http.StatusOK, reports)
}

// Upsert atende POST /api/v1/field-reports retornando 201 na criação (v1) e 200 na atualização (v2+).
func (h *FieldReportsHandler) Upsert(w http.ResponseWriter, r *http.Request) {
	var report models.FieldReport
	if err := json.NewDecoder(r.Body).Decode(&report); err != nil {
		writeError(w, http.StatusBadRequest, "JSON inválido")
		return
	}

	saved, err := h.service.Upsert(report)
	if err != nil {
		mapFieldReportError(w, err)
		return
	}

	status := http.StatusOK
	if saved.Version == 1 {
		status = http.StatusCreated
	}
	writeJSON(w, status, saved)
}
