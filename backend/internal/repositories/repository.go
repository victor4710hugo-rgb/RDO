package repositories

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/victor4710hugo-rgb/RDO/backend/internal/models"
)

type Repository struct {
	db *sql.DB
}

func New(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func parseTime(value string) (time.Time, error) {
	return time.Parse(time.RFC3339, value)
}

func (r *Repository) ListProjects() ([]models.Project, error) {
	const query = `
		SELECT
			p.id,
			p.client_id,
			c.name,
			p.name,
			p.workforce_control_mode,
			p.status,
			p.created_at,
			p.updated_at
		FROM projects p
		INNER JOIN clients c ON c.id = p.client_id
		WHERE p.deleted_at IS NULL
		ORDER BY p.name
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("listar projetos: %w", err)
	}
	defer rows.Close()

	var projects []models.Project
	for rows.Next() {
		var p models.Project
		var createdAt, updatedAt string
		if err := rows.Scan(
			&p.ID,
			&p.ClientID,
			&p.ClientName,
			&p.Name,
			&p.WorkforceControlMode,
			&p.Status,
			&createdAt,
			&updatedAt,
		); err != nil {
			return nil, fmt.Errorf("ler projeto: %w", err)
		}
		p.CreatedAt, err = parseTime(createdAt)
		if err != nil {
			return nil, fmt.Errorf("converter created_at do projeto: %w", err)
		}
		p.UpdatedAt, err = parseTime(updatedAt)
		if err != nil {
			return nil, fmt.Errorf("converter updated_at do projeto: %w", err)
		}
		projects = append(projects, p)
	}
	return projects, rows.Err()
}

func (r *Repository) ListWorkOrders() ([]models.WorkOrder, error) {
	const query = `
		SELECT
			id,
			project_id,
			unit_id,
			number,
			title,
			priority,
			status,
			created_at,
			updated_at
		FROM work_orders
		WHERE deleted_at IS NULL
		ORDER BY number
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("listar OS: %w", err)
	}
	defer rows.Close()

	var orders []models.WorkOrder
	for rows.Next() {
		var o models.WorkOrder
		var createdAt, updatedAt string
		if err := rows.Scan(
			&o.ID,
			&o.ProjectID,
			&o.UnitID,
			&o.Number,
			&o.Title,
			&o.Priority,
			&o.Status,
			&createdAt,
			&updatedAt,
		); err != nil {
			return nil, fmt.Errorf("ler OS: %w", err)
		}
		o.CreatedAt, err = parseTime(createdAt)
		if err != nil {
			return nil, fmt.Errorf("converter created_at da OS: %w", err)
		}
		o.UpdatedAt, err = parseTime(updatedAt)
		if err != nil {
			return nil, fmt.Errorf("converter updated_at da OS: %w", err)
		}
		orders = append(orders, o)
	}
	return orders, rows.Err()
}

func (r *Repository) ListActivities() ([]models.Activity, error) {
	const query = `
		SELECT
			id,
			work_order_id,
			title,
			priority,
			status,
			planned_progress,
			actual_progress,
			created_at,
			updated_at
		FROM activities
		WHERE deleted_at IS NULL
		ORDER BY title
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("listar atividades: %w", err)
	}
	defer rows.Close()

	var activities []models.Activity
	for rows.Next() {
		var a models.Activity
		var createdAt, updatedAt string
		if err := rows.Scan(
			&a.ID,
			&a.WorkOrderID,
			&a.Title,
			&a.Priority,
			&a.Status,
			&a.PlannedProgress,
			&a.ActualProgress,
			&createdAt,
			&updatedAt,
		); err != nil {
			return nil, fmt.Errorf("ler atividade: %w", err)
		}
		a.CreatedAt, err = parseTime(createdAt)
		if err != nil {
			return nil, fmt.Errorf("converter created_at da atividade: %w", err)
		}
		a.UpdatedAt, err = parseTime(updatedAt)
		if err != nil {
			return nil, fmt.Errorf("converter updated_at da atividade: %w", err)
		}
		activities = append(activities, a)
	}
	return activities, rows.Err()
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