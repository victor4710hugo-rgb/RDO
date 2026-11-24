package database

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

func Open() (*sql.DB, error) {
	dir := "data"
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, fmt.Errorf("criar pasta de dados: %w", err)
	}

	db, err := sql.Open("sqlite", filepath.Join(dir, "rdo.db"))
	if err != nil {
		return nil, fmt.Errorf("abrir banco: %w", err)
	}

	if _, err := db.Exec("PRAGMA journal_mode=WAL; PRAGMA foreign_keys=ON;"); err != nil {
		return nil, fmt.Errorf("configurar pragmas: %w", err)
	}

	return db, nil
}