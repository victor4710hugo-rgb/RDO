package repositories

import (
	"database/sql"
	"time"
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
	rows, err := r.db.Query("SELECT id, name, client_name, created_at, updated_at FROM projects ORDER BY name")
	if err != nil {
		return nil, fmt.Errorf("listar projetos: %w", err)
	}
	defer rows.Close()

	var projects []models.Project
	for rows.Next() {
		var p models.Project
		var createdAt, updatedAt string
		if err := rows.Scan(&p.ID, &p.Name, &p.ClientName, &createdAt, &updatedAt); err != nil {
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
	rows, err := r.db.Query("SELECT id, project_id, number, title, created_at, updated_at FROM work_orders ORDER BY number")
	if err != nil {
		return nil, fmt.Errorf("listar OS: %w", err)
	}
	defer rows.Close()

	var orders []models.WorkOrder
	for rows.Next() {
		var o models.WorkOrder
		var createdAt, updatedAt string
		if err := rows.Scan(&o.ID, &o.ProjectID, &o.Number, &o.Title, &createdAt, &updatedAt); err != nil {
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
	rows, err := r.db.Query("SELECT id, work_order_id, title, status, created_at, updated_at FROM activities ORDER BY title")
	if err != nil {
		return nil, fmt.Errorf("listar atividades: %w", err)
	}
	defer rows.Close()

	var activities []models.Activity
	for rows.Next() {
		var a models.Activity
		var createdAt, updatedAt string
		if err := rows.Scan(&a.ID, &a.WorkOrderID, &a.Title, &a.Status, &createdAt, &updatedAt); err != nil {
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
func (r *Repository) ListFieldReports() ([]models.FieldReport, error) {
	const query = `
		SELECT id, project_id, report_number, report_date, status, version, payload, created_at, updated_at
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
			&fr.ID, &fr.ProjectID, &fr.ReportNumber, &fr.ReportDate,
			&fr.Status, &fr.Version, &fr.Payload, &createdAt, &updatedAt,
		); err != nil {
			return nil, fmt.Errorf("ler RDO: %w", err)
		}
		if fr.CreatedAt, err = parseTime(createdAt); err != nil {
			return nil, err
		}
		if fr.UpdatedAt, err = parseTime(updatedAt); err != nil {
			return nil, err
		}
		reports = append(reports, fr)
	}
	return reports, rows.Err()
}