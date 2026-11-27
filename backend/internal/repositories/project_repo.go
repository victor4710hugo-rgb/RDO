package repositories

import (
	"fmt"

	"github.com/victor4710hugo-rgb/RDO/backend/internal/models"
)

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

	projects := make([]models.Project, 0)
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
		if p.CreatedAt, err = parseTime(createdAt); err != nil {
			return nil, fmt.Errorf("converter created_at do projeto: %w", err)
		}
		if p.UpdatedAt, err = parseTime(updatedAt); err != nil {
			return nil, fmt.Errorf("converter updated_at do projeto: %w", err)
		}
		projects = append(projects, p)
	}
	return projects, rows.Err()
}