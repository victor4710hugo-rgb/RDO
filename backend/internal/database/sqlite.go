
package database

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

func Open() (*sql.DB, error) {
	dir := "data"

	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("criar pasta de dados: %w", err)
	}

	dbPath, err := filepath.Abs(filepath.Join(dir, "rdo.db"))
	if err != nil {
		return nil, fmt.Errorf("obter caminho do banco: %w", err)
	}

	// Configura os PRAGMAs para cada conexão aberta pelo driver.
	dsn := "file:" + filepath.ToSlash(dbPath) +
		"?_pragma=foreign_keys(1)" +
		"&_pragma=journal_mode(WAL)" +
		"&_pragma=busy_timeout(5000)"

	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, fmt.Errorf("abrir banco: %w", err)
	}

	ctx := context.Background()

	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, fmt.Errorf("conectar ao banco: %w", err)
	}

	if err := Migrate(db); err != nil {
		db.Close()
		return nil, fmt.Errorf("inicializar banco: %w", err)
	}

	return db, nil
}
