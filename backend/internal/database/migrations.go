package database

import (
	"crypto/sha256"
	"database/sql"
	_ "embed"
	"encoding/hex"
	"fmt"
	"time"
)

// SQL das migrations embutido no binário.
// Os arquivos .sql devem ficar na mesma pasta deste arquivo .go.

//go:embed schema_v1.sql
var schemaV1 string

//go:embed 0002_auth_users.sql
var migration0002 string

// migration descreve uma migration versionada.
type migration struct {
	version int
	name    string
	sql     string
}

// migrations lista, em ordem crescente de versão, todas as migrations.
// Para adicionar uma nova: crie o .sql, embuta com //go:embed e inclua aqui.
// Nunca altere o conteúdo de uma migration já aplicada: o checksum mudaria
// e o Migrate retornaria erro.
var migrations = []migration{
	{version: 1, name: "initial_schema", sql: schemaV1},
	{version: 2, name: "auth_users", sql: migration0002},
}

// Migrate garante que o banco esteja com todas as migrations aplicadas.
// É idempotente: migrations já aplicadas (mesmo checksum) são ignoradas.
func Migrate(db *sql.DB) error {
	if db == nil {
		return fmt.Errorf("conexão com banco não inicializada")
	}

	// Tabela de controle das versões aplicadas.
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

// applyMigration aplica uma única migration dentro de uma transação:
// ou o SQL e o registro em schema_migrations são confirmados juntos,
// ou nada é alterado.
func applyMigration(db *sql.DB, m migration) error {
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("iniciar migration %d: %w", m.version, err)
	}
	defer tx.Rollback() // no-op se o Commit já tiver sido feito

	hash := sha256.Sum256([]byte(m.sql))
	checksum := hex.EncodeToString(hash[:])

	// Verifica se a versão já foi aplicada.
	var existing string
	err = tx.QueryRow(
		"SELECT checksum FROM schema_migrations WHERE version = ?",
		m.version,
	).Scan(&existing)

	if err == nil {
		// Já aplicada: o conteúdo não pode ter mudado.
		if existing != checksum {
			return fmt.Errorf("migration %d já aplicada com checksum diferente", m.version)
		}
		return tx.Commit()
	}
	if err != sql.ErrNoRows {
		return fmt.Errorf("consultar migration %d: %w", m.version, err)
	}

	// Ainda não aplicada: executa o SQL.
	if _, err := tx.Exec(m.sql); err != nil {
		return fmt.Errorf("executar migration %d (%s): %w", m.version, m.name, err)
	}

	// Registra a aplicação.
	_, err = tx.Exec(`
		INSERT INTO schema_migrations
			(version, name, checksum, applied_at)
		VALUES (?, ?, ?, ?)
	`,
		m.version,
		m.name,
		checksum,
		time.Now().UTC().Format(time.RFC3339Nano),
	)
	if err != nil {
		return fmt.Errorf("registrar migration %d: %w", m.version, err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("confirmar migration %d: %w", m.version, err)
	}
	return nil
}