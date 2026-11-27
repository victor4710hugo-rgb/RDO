package repositories

import (
	"fmt"

	"github.com/victor4710hugo-rgb/RDO/backend/internal/models"
)

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

	activities := make([]models.Activity, 0)
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
		if a.CreatedAt, err = parseTime(createdAt); err != nil {
			return nil, fmt.Errorf("converter created_at da atividade: %w", err)
		}
		if a.UpdatedAt, err = parseTime(updatedAt); err != nil {
			return nil, fmt.Errorf("converter updated_at da atividade: %w", err)
		}
		activities = append(activities, a)
	}
	return activities, rows.Err()
}