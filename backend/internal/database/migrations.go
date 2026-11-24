package database

import (
	"database/sql"
	"fmt"
)

func Migrate(db *sql.DB) error {
	statements := []string{
		`CREATE TABLE IF NOT EXISTS projects (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			client_name TEXT NOT NULL,
			created_at TEXT NOT NULL,
			updated_at TEXT NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS work_orders (
			id TEXT PRIMARY KEY,
			project_id TEXT NOT NULL,
			number TEXT NOT NULL,
			title TEXT NOT NULL,
			created_at TEXT NOT NULL,
			updated_at TEXT NOT NULL,
			FOREIGN KEY (project_id) REFERENCES projects(id)
		)`,
		`CREATE TABLE IF NOT EXISTS activities (
			id TEXT PRIMARY KEY,
			work_order_id TEXT NOT NULL,
			title TEXT NOT NULL,
			status TEXT NOT NULL,
			created_at TEXT NOT NULL,
			updated_at TEXT NOT NULL,
			FOREIGN KEY (work_order_id) REFERENCES work_orders(id)
		)`,
		`CREATE TABLE IF NOT EXISTS field_reports (
			id TEXT PRIMARY KEY,
			work_order_id TEXT NOT NULL,
			report_date TEXT NOT NULL,
			version INTEGER NOT NULL DEFAULT 1,
			payload TEXT NOT NULL,
			created_at TEXT NOT NULL,
			updated_at TEXT NOT NULL,
			FOREIGN KEY (work_order_id) REFERENCES work_orders(id)
		)`,
		`CREATE TABLE IF NOT EXISTS sync_records (
			id TEXT PRIMARY KEY,
			entity_type TEXT NOT NULL,
			entity_id TEXT NOT NULL,
			status TEXT NOT NULL,
			created_at TEXT NOT NULL,
			updated_at TEXT NOT NULL
		)`,
	}

	for _, statement := range statements {
		if _, err := db.Exec(statement); err != nil {
			return fmt.Errorf("executar migration: %w", err)
		}
	}

	return nil
}