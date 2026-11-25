package repositories

import (
	"fmt"

	"github.com/victor4710hugo-rgb/RDO/backend/internal/models"
)

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

	orders := make([]models.WorkOrder, 0)
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
		if o.CreatedAt, err = parseTime(createdAt); err != nil {
			return nil, fmt.Errorf("converter created_at da OS: %w", err)
		}
		if o.UpdatedAt, err = parseTime(updatedAt); err != nil {
			return nil, fmt.Errorf("converter updated_at da OS: %w", err)
		}
		orders = append(orders, o)
	}
	return orders, rows.Err()
}