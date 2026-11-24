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

func (s *FieldReportService) Upsert(report models.FieldReport) (models.FieldReport, error) {
	if report.ID == "" {
		return report, fmt.Errorf("id do RDO é obrigatório")
	}
	if report.WorkOrderID == "" {
		return report, fmt.Errorf("id da OS é obrigatório")
	}
	if report.ReportDate == "" {
		return report, fmt.Errorf("data do RDO é obrigatória")
	}
	if report.Payload == "" {
		return report, fmt.Errorf("payload do RDO é obrigatório")
	}

	return s.repo.UpsertFieldReport(report)
}	