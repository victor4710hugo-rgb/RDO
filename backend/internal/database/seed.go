package database

import (
	"database/sql"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// Seed popula o banco com dados de demonstração (cliente, projeto, OS e
// atividades). É idempotente: se já existir algum projeto, não faz nada.
func Seed(db *sql.DB) error {
	if err := seedOperational(db); err != nil {
		return err
	}
	return seedAdmin(db)
}

func seedOperational(db *sql.DB) error {
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

func seedAdmin(db *sql.DB) error {
	email := strings.ToLower(strings.TrimSpace(os.Getenv("ADMIN_EMAIL")))
	password := os.Getenv("ADMIN_PASSWORD")
	if email == "" || password == "" {
		return nil
	}

	var exists int
	if err := db.QueryRow("SELECT COUNT(*) FROM users WHERE LOWER(email) = ?", email).Scan(&exists); err != nil {
		return fmt.Errorf("verificar admin inicial: %w", err)
	}
	if exists > 0 {
		return nil
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("gerar hash do admin: %w", err)
	}

	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("iniciar transacao de seed admin: %w", err)
	}
	defer tx.Rollback()

	now := time.Now().UTC().Format(time.RFC3339)
	adminID := uuid.NewString()

	if _, err := tx.Exec(
		"INSERT INTO users (id, name, email, status, password_hash, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?)",
		adminID, "Administrador RDO", email, "active", string(hash), now, now,
	); err != nil {
		return fmt.Errorf("inserir admin inicial: %w", err)
	}

	var roleID string
	if err := tx.QueryRow("SELECT id FROM roles WHERE name = 'admin'").Scan(&roleID); err != nil {
		return fmt.Errorf("buscar papel admin: %w", err)
	}

	if _, err := tx.Exec(
		"INSERT INTO user_roles (user_id, role_id, created_at) VALUES (?, ?, ?)",
		adminID, roleID, now,
	); err != nil {
		return fmt.Errorf("vincular papel admin: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("confirmar seed admin: %w", err)
	}

	return nil
}
