
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

func Migrate(db *sql.DB) error {
	if db == nil {
		return fmt.Errorf("conexão com banco não inicializada")
	}

	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("iniciar migration: %w", err)
	}
	defer tx.Rollback()

	// Controle de versões das migrations.
	_, err = tx.Exec(`
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

	const version = 1

	hash := sha256.Sum256([]byte(schemaV1))
	checksum := hex.EncodeToString(hash[:])

	var existingChecksum string
	err = tx.QueryRow(
		"SELECT checksum FROM schema_migrations WHERE version = ?",
		version,
	).Scan(&existingChecksum)

	if err == nil {
		if existingChecksum != checksum {
			return fmt.Errorf("migration %d já aplicada com checksum diferente", version)
		}
		return tx.Commit()
	}

	if err != sql.ErrNoRows {
		return fmt.Errorf("consultar migrations: %w", err)
	}

	// Executa a estrutura completa do banco.
	if _, err := tx.Exec(schemaV1); err != nil {
		return fmt.Errorf("executar schema v1: %w", err)
	}

	_, err = tx.Exec(`
		INSERT INTO schema_migrations
			(version, name, checksum, applied_at)
		VALUES (?, ?, ?, ?)
	`,
		version,
		"initial_schema",
		checksum,
		time.Now().UTC().Format(time.RFC3339Nano),
	)
	if err != nil {
		return fmt.Errorf("registrar migration: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("confirmar migration: %w", err)
	}

	return nil
}
