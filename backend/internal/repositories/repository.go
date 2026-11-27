package repositories

import (
	"database/sql"
	"time"
)

// Repository é o único ponto de acesso ao SQLite. Os métodos ficam
// divididos por domínio em project_repo.go, os_repo.go, activity_repo.go,
// rdo_repo.go e user_repo.go.
type Repository struct {
	db *sql.DB
}

// New cria o repositório sobre a conexão já aberta.
func New(db *sql.DB) *Repository {
	return &Repository{db: db}
}

// parseTime converte timestamps RFC3339 gravados no banco; é compartilhado
// pelos repositórios de cada domínio.
func parseTime(value string) (time.Time, error) {
	return time.Parse(time.RFC3339, value)
}
