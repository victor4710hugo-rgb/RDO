package services

import (
	"strings"

	"github.com/victor4710hugo-rgb/RDO/backend/internal/models"
	"github.com/victor4710hugo-rgb/RDO/backend/internal/repositories"
)

// FieldReportService concentra as regras de validação do RDO antes do acesso ao banco.
type FieldReportService struct {
	repo *repositories.Repository
}

// NewFieldReportService instancia o serviço de RDO com o repositório injetado.
func NewFieldReportService(repo *repositories.Repository) *FieldReportService {
	return &FieldReportService{repo: repo}
}

// List retorna todos os RDOs persistidos.
func (s *FieldReportService) List() ([]models.FieldReport, error) {
	return s.repo.ListFieldReports()
}

// Upsert normaliza e valida os campos obrigatórios (retornando ErrValidation -> 400) antes de gravar.
func (s *FieldReportService) Upsert(report models.FieldReport) (models.FieldReport, error) {
	report.ID = strings.TrimSpace(report.ID)
	report.ProjectID = strings.TrimSpace(report.ProjectID)
	report.ReportNumber = strings.TrimSpace(report.ReportNumber)
	report.ReportDate = strings.TrimSpace(report.ReportDate)
	report.Status = strings.TrimSpace(report.Status)

	if report.ID == "" {
		return report, validationError("id do RDO é obrigatório")
	}
	if report.ProjectID == "" {
		return report, validationError("id do projeto é obrigatório")
	}
	if report.ReportNumber == "" {
		return report, validationError("número do RDO é obrigatório")
	}
	if report.ReportDate == "" {
		return report, validationError("data do RDO é obrigatória")
	}
	if report.Payload == nil || strings.TrimSpace(*report.Payload) == "" {
		return report, validationError("payload do RDO é obrigatório")
	}

	return s.repo.UpsertFieldReport(report)
}
