package services

import (
	"fmt"

	"github.com/victor4710hugo-rgb/RDO/backend/internal/models"
	"github.com/victor4710hugo-rgb/RDO/backend/internal/repositories"
)

type FieldReportService struct {
	repo *repositories.Repository
}

func NewFieldReportService(repo *repositories.Repository) *FieldReportService {
	return &FieldReportService{repo: repo}
}

func (s *FieldReportService) List() ([]models.FieldReport, error) {
	return s.repo.ListFieldReports()
}

func (s *FieldReportService) Upsert(report models.FieldReport) (models.FieldReport, error) {
	if report.ID == "" {
		return report, fmt.Errorf("id do RDO é obrigatório")
	}
	if report.ProjectID == "" {
		return report, fmt.Errorf("id do projeto é obrigatório")
	}
	if report.ReportNumber == "" {
		return report, fmt.Errorf("número do RDO é obrigatório")
	}
	if report.ReportDate == "" {
		return report, fmt.Errorf("data do RDO é obrigatória")
	}
	if report.Payload == nil || *report.Payload == "" {
		return report, fmt.Errorf("payload do RDO é obrigatório")
	}

	return s.repo.UpsertFieldReport(report)
}