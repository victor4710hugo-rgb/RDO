package repositories

import (
	"database/sql"
	"time"
)

type Repository struct {
	db *sql.DB
}

func New(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func parseTime(value string) (time.Time, error) {
	return time.Parse(time.RFC3339, value)
}