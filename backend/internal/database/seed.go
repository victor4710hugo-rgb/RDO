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

	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("iniciar transacao de seed: %w", err)
	}
	defer tx.Rollback()

	now := time.Now().UTC().Format(time.RFC3339)

	clientID := uuid.NewString()
	if _, err := tx.Exec(
		"INSERT INTO clients (id, name, created_at, updated_at) VALUES (?, ?, ?, ?)",
		clientID, "Cliente Demonstracao", now, now,
	); err != nil {
		return fmt.Errorf("inserir cliente: %w", err)
	}

	projectID := uuid.NewString()
	if _, err := tx.Exec(
		"INSERT INTO projects (id, client_id, name, created_at, updated_at) VALUES (?, ?, ?, ?, ?)",
		projectID, clientID, "Projeto Demonstracao", now, now,
	); err != nil {
		return fmt.Errorf("inserir projeto: %w", err)
	}

	workOrderID := uuid.NewString()
	if _, err := tx.Exec(
		"INSERT INTO work_orders (id, project_id, number, title, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)",
		workOrderID, projectID, "OS-001", "Manutencao preventiva", now, now,
	); err != nil {
		return fmt.Errorf("inserir OS: %w", err)
	}

	activities := []string{"Inspecao inicial", "Substituicao de componente", "Teste final"}
	for _, title := range activities {
		if _, err := tx.Exec(
			"INSERT INTO activities (id, work_order_id, title, status, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)",
			uuid.NewString(), workOrderID, title, "pending", now, now,
		); err != nil {
			return fmt.Errorf("inserir atividade: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("confirmar seed: %w", err)
	}

	return nil
}