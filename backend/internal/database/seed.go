package database

import (
	"database/sql"
	"fmt"
	"time"
)

// Seed popula o banco com dados de demonstração (cliente, projeto, OS e
// atividades). É idempotente: se já existir algum projeto, não faz nada.
func Seed(db *sql.DB) error {
	var count int
	if err := db.QueryRow("SELECT COUNT(*) FROM projects").Scan(&count); err != nil {
		return fmt.Errorf("consultar seed de projetos: %w", err)
	}
	if count > 0 {
		return nil
	}

	// Uma única transação: ou todo o seed é gravado, ou nada.
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("iniciar seed: %w", err)
	}
	defer tx.Rollback() // no-op se o Commit já tiver sido feito

	now := time.Now().UTC().Format(time.RFC3339)

	// IDs fixos e legíveis: facilitam testes e chamadas manuais da API.
	const (
		clientID    = "cli-001"
		projectID   = "proj-001"
		workOrderID = "os-001"
	)

	// Cliente
	if _, err := tx.Exec(
		`INSERT INTO clients (id, name, status, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?)`,
		clientID, "Cliente Demonstração S.A.", "active", now, now,
	); err != nil {
		return fmt.Errorf("inserir cliente: %w", err)
	}

	// Projeto
	if _, err := tx.Exec(
		`INSERT INTO projects (id, client_id, name, workforce_control_mode, status, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?)`,
		projectID, clientID, "Obra Industrial — Unidade Norte", "both", "in_progress", now, now,
	); err != nil {
		return fmt.Errorf("inserir projeto: %w", err)
	}

	// Ordem de serviço
	if _, err := tx.Exec(
		`INSERT INTO work_orders (id, project_id, number, title, priority, status, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		workOrderID, projectID, "OS-2026-001", "Montagem Eletromecânica Linha A", "high", "in_progress", now, now,
	); err != nil {
		return fmt.Errorf("inserir OS: %w", err)
	}

	// Atividades (em lista para facilitar a inclusão de novas)
	activities := []struct {
		id              string
		title           string
		priority        string
		status          string
		plannedProgress float64
		actualProgress  float64
	}{
		{"act-001", "Instalação de eletrocalhas e painéis", "high", "in_progress", 60.0, 45.0},
	}

	for _, a := range activities {
		if _, err := tx.Exec(
			`INSERT INTO activities (id, work_order_id, title, priority, status, planned_progress, actual_progress, created_at, updated_at)
			 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			a.id, workOrderID, a.title, a.priority, a.status, a.plannedProgress, a.actualProgress, now, now,
		); err != nil {
			return fmt.Errorf("inserir atividade %s: %w", a.id, err)
		}
	}

	// Confirma tudo de uma vez.
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("confirmar seed: %w", err)
	}
	return nil
}