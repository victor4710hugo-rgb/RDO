package repositories

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/victor4710hugo-rgb/RDO/backend/internal/models"
)

var (
	ErrUserNotFound = errors.New("usuário não encontrado")
	ErrEmailExists  = errors.New("e-mail já cadastrado")
)

const userSelect = `
	SELECT u.id, u.name, u.email, COALESCE(r.name, ''), u.status, u.password_hash, u.created_at, u.updated_at
	FROM users u
	LEFT JOIN user_roles ur ON ur.user_id = u.id
	LEFT JOIN roles r ON r.id = ur.role_id
`

func scanUser(scanner interface{ Scan(...any) error }) (models.User, error) {
	var u models.User
	var createdAt, updatedAt string
	if err := scanner.Scan(&u.ID, &u.Name, &u.Email, &u.Role, &u.Status, &u.PasswordHash, &createdAt, &updatedAt); err != nil {
		return u, err
	}
	var err error
	if u.CreatedAt, err = parseTime(createdAt); err != nil {
		return u, err
	}
	if u.UpdatedAt, err = parseTime(updatedAt); err != nil {
		return u, err
	}
	return u, nil
}

func (r *Repository) ListUsers() ([]models.User, error) {
	rows, err := r.db.Query(userSelect + ` WHERE u.deleted_at IS NULL ORDER BY u.created_at DESC`)
	if err != nil {
		return nil, fmt.Errorf("listar usuários: %w", err)
	}
	defer rows.Close()

	users := make([]models.User, 0)
	for rows.Next() {
		u, err := scanUser(rows)
		if err != nil {
			return nil, fmt.Errorf("ler usuário: %w", err)
		}
		users = append(users, u)
	}
	return users, rows.Err()
}

func (r *Repository) GetUserByID(id string) (models.User, error) {
	row := r.db.QueryRow(userSelect+` WHERE u.id = ? AND u.deleted_at IS NULL`, id)
	u, err := scanUser(row)
	if errors.Is(err, sql.ErrNoRows) {
		return u, ErrUserNotFound
	}
	return u, err
}

func (r *Repository) GetUserByEmail(email string) (models.User, error) {
	row := r.db.QueryRow(userSelect+` WHERE LOWER(u.email) = LOWER(?) AND u.deleted_at IS NULL`, email)
	u, err := scanUser(row)
	if errors.Is(err, sql.ErrNoRows) {
		return u, ErrUserNotFound
	}
	return u, err
}

func (r *Repository) CreateUser(u models.User) (models.User, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return u, err
	}
	defer tx.Rollback()

	now := time.Now().UTC()
	u.CreatedAt, u.UpdatedAt = now, now
	if u.Status == "" {
		u.Status = "active"
	}

	_, err = tx.Exec(
		`INSERT INTO users (id, name, email, status, password_hash, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?)`,
		u.ID, u.Name, u.Email, u.Status, u.PasswordHash, now.Format(time.RFC3339), now.Format(time.RFC3339),
	)
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "unique") {
			return u, ErrEmailExists
		}
		return u, fmt.Errorf("inserir usuário: %w", err)
	}

	var roleID string
	if err := tx.QueryRow(`SELECT id FROM roles WHERE name = ?`, u.Role).Scan(&roleID); err != nil {
		return u, fmt.Errorf("buscar papel %s: %w", u.Role, err)
	}
	if _, err := tx.Exec(`INSERT INTO user_roles (user_id, role_id, created_at) VALUES (?, ?, ?)`, u.ID, roleID, now.Format(time.RFC3339)); err != nil {
		return u, fmt.Errorf("vincular papel: %w", err)
	}

	return u, tx.Commit()
}

func (r *Repository) UpdateUser(u models.User) (models.User, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return u, err
	}
	defer tx.Rollback()

	now := time.Now().UTC()
	u.UpdatedAt = now

	var res sql.Result
	if u.PasswordHash != "" {
		res, err = tx.Exec(
			`UPDATE users SET name = ?, email = ?, status = ?, password_hash = ?, updated_at = ? WHERE id = ? AND deleted_at IS NULL`,
			u.Name, u.Email, u.Status, u.PasswordHash, now.Format(time.RFC3339), u.ID,
		)
	} else {
		res, err = tx.Exec(
			`UPDATE users SET name = ?, email = ?, status = ?, updated_at = ? WHERE id = ? AND deleted_at IS NULL`,
			u.Name, u.Email, u.Status, now.Format(time.RFC3339), u.ID,
		)
	}
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "unique") {
			return u, ErrEmailExists
		}
		return u, fmt.Errorf("atualizar usuário: %w", err)
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return u, ErrUserNotFound
	}

	var roleID string
	if err := tx.QueryRow(`SELECT id FROM roles WHERE name = ?`, u.Role).Scan(&roleID); err != nil {
		return u, fmt.Errorf("buscar papel %s: %w", u.Role, err)
	}
	if _, err := tx.Exec(`DELETE FROM user_roles WHERE user_id = ?`, u.ID); err != nil {
		return u, err
	}
	if _, err := tx.Exec(`INSERT INTO user_roles (user_id, role_id, created_at) VALUES (?, ?, ?)`, u.ID, roleID, now.Format(time.RFC3339)); err != nil {
		return u, err
	}

	return u, tx.Commit()
}

func (r *Repository) DeleteUser(id string) error {
	now := time.Now().UTC().Format(time.RFC3339)
	res, err := r.db.Exec(`UPDATE users SET deleted_at = ?, status = 'inactive', updated_at = ? WHERE id = ? AND deleted_at IS NULL`, now, now, id)
	if err != nil {
		return fmt.Errorf("excluir usuário: %w", err)
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return ErrUserNotFound
	}
	return nil
}