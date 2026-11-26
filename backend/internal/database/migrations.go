package database

import (
	"crypto/sha256"
	"database/sql"
	_ "embed"
	"encoding/hex"
	"fmt"
	"time"
)

//go:embed schema_v1.sql
var schemaV1 string

//go:embed 0002_auth_users.sql
var migration0002 string

type migration struct {
	version int
	name    string
	sql     string
}

var migrations = []migration{
	{version: 1, name: "initial_schema", sql: schemaV1},
	{version: 2, name: "auth_users", sql: migration0002},
}

func Migrate(db *sql.DB) error {
	if db == nil {
		return fmt.Errorf("conexão com banco não inicializada")
	}

	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version INTEGER PRIMARY KEY,
			name TEXT NOT NULL,
			checksum TEXT NOT NULL,
			applied_at TEXT NOT NULL
		)
	`)
	if err != nil {
		return fmt.Errorf("criar controle de migrations: %w", err)
	}

	for _, m := range migrations {
		if err := applyMigration(db, m); err != nil {
			return err
		}
	}
	return nil
}

func applyMigration(db *sql.DB, m migration) error {
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("iniciar migration %d: %w", m.version, err)
	}
	defer tx.Rollback()

	hash := sha256.Sum256([]byte(m.sql))
	checksum := hex.EncodeToString(hash[:])

	var existing string
	err = tx.QueryRow("SELECT checksum FROM schema_migrations WHERE version = ?", m.version).Scan(&existing)
	if err == nil {
		if existing != checksum {
			return fmt.Errorf("migration %d já aplicada com checksum diferente", m.version)
		}
		return tx.Commit()
	}
	if err != sql.ErrNoRows {
		return fmt.Errorf("consultar migration %d: %w", m.version, err)
	}

	if _, err := tx.Exec(m.sql); err != nil {
		return fmt.Errorf("executar migration %d (%s): %w", m.version, m.name, err)
	}

	_, err = tx.Exec(
		`INSERT INTO schema_migrations (version, name, checksum, applied_at) VALUES (?, ?, ?, ?)`,
		m.version, m.name, checksum, time.Now().UTC().Format(time.RFC3339Nano),
	)
	if err != nil {
		return fmt.Errorf("registrar migration %d: %w", m.version, err)
	}
	return tx.Commit()
}