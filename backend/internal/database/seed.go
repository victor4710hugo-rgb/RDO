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

func Seed(db *sql.DB) error {
	if err := seedOperational(db); err != nil {
		return err
	}
	return seedAdmin(db)
}

func seedOperational(db *sql.DB) error {
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