package database

import (
	"database/sql"
	"fmt"
	"time"
)

func Seed(db *sql.DB) error {
	var count int
	if err := db.QueryRow("SELECT COUNT(*) FROM projects").Scan(&count); err != nil {
		return fmt.Errorf("consultar seed de projetos: %w", err)
	}
	if count > 0 {
		return nil
	}

	now := time.Now().UTC().Format(time.RFC3339)

	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("iniciar seed: %w", err)
	}
	defer tx.Rollback()

	if _, err := tx.Exec(
		`INSERT INTO clients (id, name, status, created_at, updated_at) VALUES (?, ?, ?, ?, ?)`,
		"cli-001", "Cliente Demonstração S.A.", "active", now, now,
	); err != nil {
		return fmt.Errorf("inserir cliente: %w", err)
	}

	if _, err := tx.Exec(
		`INSERT INTO projects (id, client_id, name, workforce_control_mode, status, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?)`,
		"proj-001", "cli-001", "Obra Industrial — Unidade Norte", "both", "in_progress", now, now,
	); err != nil {
		return fmt.Errorf("inserir projeto: %w", err)
	}

	if _, err := tx.Exec(
		`INSERT INTO work_orders (id, project_id, number, title, priority, status, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		"os-001", "proj-001", "OS-2026-001", "Montagem Eletromecânica Linha A", "high", "in_progress", now, now,
	); err != nil {
		return fmt.Errorf("inserir OS: %w", err)
	}

	if _, err := tx.Exec(
		`INSERT INTO activities (id, work_order_id, title, priority, status, planned_progress, actual_progress, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		"act-001", "os-001", "Instalação de eletrocalhas e painéis", "high", "in_progress", 60.0, 45.0, now, now,
	); err != nil {
		return fmt.Errorf("inserir atividade: %w", err)
	}

	return tx.Commit()
}