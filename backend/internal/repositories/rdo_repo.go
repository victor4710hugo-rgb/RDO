package repositories

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/victor4710hugo-rgb/RDO/backend/internal/models"
)

func (r *Repository) ListFieldReports() ([]models.FieldReport, error) {
	const query = `
		SELECT
			id,
			project_id,
			report_number,
			report_date,
			status,
			version,
			payload,
			created_at,
			updated_at
		FROM field_reports
		ORDER BY report_date DESC, created_at DESC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("listar RDOs: %w", err)
	}
	defer rows.Close()

	reports := make([]models.FieldReport, 0)
	for rows.Next() {
		var fr models.FieldReport
		var createdAt, updatedAt string
		if err := rows.Scan(
			&fr.ID,
			&fr.ProjectID,
			&fr.ReportNumber,
			&fr.ReportDate,
			&fr.Status,
			&fr.Version,
			&fr.Payload,
			&createdAt,
			&updatedAt,
		); err != nil {
			return nil, fmt.Errorf("ler RDO: %w", err)
		}
		if fr.CreatedAt, err = parseTime(createdAt); err != nil {
			return nil, fmt.Errorf("converter created_at do RDO: %w", err)
		}
		if fr.UpdatedAt, err = parseTime(updatedAt); err != nil {
			return nil, fmt.Errorf("converter updated_at do RDO: %w", err)
		}
		reports = append(reports, fr)
	}
	return reports, rows.Err()
}

func (r *Repository) UpsertFieldReport(report models.FieldReport) (models.FieldReport, error) {
	now := time.Now().UTC()
	if report.Status == "" {
		report.Status = "draft"
	}

	var currentVersion int
	err := r.db.QueryRow("SELECT version FROM field_reports WHERE id = ?", report.ID).Scan(&currentVersion)
	switch {
	case err == sql.ErrNoRows:
		report.Version = 1
		report.CreatedAt = now
		report.UpdatedAt = now
		_, err = r.db.Exec(
			`INSERT INTO field_reports (
				id, project_id, report_number, report_date, status, version, payload, created_at, updated_at
			) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			report.ID,
			report.ProjectID,
			report.ReportNumber,
			report.ReportDate,
			report.Status,
			report.Version,
			report.Payload,
			report.CreatedAt.Format(time.RFC3339),
			report.UpdatedAt.Format(time.RFC3339),
		)
		if err != nil {
			return report, fmt.Errorf("inserir RDO: %w", err)
		}
	case err != nil:
		return report, fmt.Errorf("consultar RDO: %w", err)
	default:
		if report.Version != currentVersion {
			return report, fmt.Errorf("conflito de versão: esperado %d, recebido %d", currentVersion, report.Version)
		}
		report.Version = currentVersion + 1
		report.UpdatedAt = now
		_, err = r.db.Exec(
			`UPDATE field_reports
			SET report_number = ?, report_date = ?, status = ?, version = ?, payload = ?, updated_at = ?
			WHERE id = ?`,
			report.ReportNumber,
			report.ReportDate,
			report.Status,
			report.Version,
			report.Payload,
			report.UpdatedAt.Format(time.RFC3339),
			report.ID,
		)
		if err != nil {
			return report, fmt.Errorf("atualizar RDO: %w", err)
		}
	}

	return report, nil
}