package database

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
)

func Seed(db *sql.DB) error {
	var count int
	if err := db.QueryRow("SELECT COUNT(*) FROM projects").Scan(&count); err != nil {
		return fmt.Errorf("contar projetos: %w", err)
	}
	if count > 0 {
		return nil
	}

	now := time.Now().UTC().Format(time.RFC3339)

	projectID := uuid.NewString()
	if _, err := db.Exec(
		"INSERT INTO projects (id, name, client_name, created_at, updated_at) VALUES (?, ?, ?, ?, ?)",
		projectID, "Projeto Demonstracao", "Cliente Demonstracao", now, now,
	); err != nil {
		return fmt.Errorf("inserir projeto: %w", err)
	}

	workOrderID := uuid.NewString()
	if _, err := db.Exec(
		"INSERT INTO work_orders (id, project_id, number, title, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)",
		workOrderID, projectID, "OS-001", "Manutencao preventiva", now, now,
	); err != nil {
		return fmt.Errorf("inserir OS: %w", err)
	}

	activities := []string{"Inspecao inicial", "Substituicao de componente", "Teste final"}
	for _, title := range activities {
		if _, err := db.Exec(
			"INSERT INTO activities (id, work_order_id, title, status, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)",
			uuid.NewString(), workOrderID, title, "pending", now, now,
		); err != nil {
			return fmt.Errorf("inserir atividade: %w", err)
		}
	}

	return nil
}