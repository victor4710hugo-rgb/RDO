package services

import (
	"errors"
	"strings"

	"github.com/google/uuid"
	"github.com/victor4710hugo-rgb/RDO/backend/internal/models"
	"github.com/victor4710hugo-rgb/RDO/backend/internal/repositories"
	"golang.org/x/crypto/bcrypt"
)

var ErrValidation = errors.New("dados inválidos")

type UserService struct {
	repo *repositories.Repository
}

func NewUserService(repo *repositories.Repository) *UserService {
	return &UserService{repo: repo}
}

func isValidRole(role string) bool {
	return role == "admin" || role == "gerente"
}

func (s *UserService) List() ([]models.User, error) {
	return s.repo.ListUsers()
}

func (s *UserService) GetByID(id string) (models.User, error) {
	return s.repo.GetUserByID(id)
}

func (s *UserService) Create(in models.CreateUserInput) (models.User, error) {
	name := strings.TrimSpace(in.Name)
	email := strings.ToLower(strings.TrimSpace(in.Email))
	role := strings.ToLower(strings.TrimSpace(in.Role))

	if name == "" || email == "" || !strings.Contains(email, "@") {
		return models.User{}, errors.Join(ErrValidation, errors.New("nome e e-mail válido são obrigatórios"))
	}
	if len(in.Password) < 8 {
		return models.User{}, errors.Join(ErrValidation, errors.New("a senha deve ter no mínimo 8 caracteres"))
	}
	if !isValidRole(role) {
		return models.User{}, errors.Join(ErrValidation, errors.New("papel permitido: admin ou gerente"))
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
	if err != nil {
		return models.User{}, err
	}

	return s.repo.CreateUser(models.User{
		ID:           uuid.NewString(),
		Name:         name,
		Email:        email,
		Role:         role,
		Status:       "active",
		PasswordHash: string(hash),
	})
}

func (s *UserService) Update(id string, in models.UpdateUserInput) (models.User, error) {
	existing, err := s.repo.GetUserByID(id)
	if err != nil {
		return models.User{}, err
	}

	name := strings.TrimSpace(in.Name)
	email := strings.ToLower(strings.TrimSpace(in.Email))
	role := strings.ToLower(strings.TrimSpace(in.Role))
	status := strings.TrimSpace(in.Status)
	if status == "" {
		status = existing.Status
	}

	if name == "" || email == "" || !strings.Contains(email, "@") {
		return models.User{}, errors.Join(ErrValidation, errors.New("nome e e-mail válido são obrigatórios"))
	}
	if !isValidRole(role) {
		return models.User{}, errors.Join(ErrValidation, errors.New("papel permitido: admin ou gerente"))
	}

	var passwordHash string
	if in.Password != "" {
		if len(in.Password) < 8 {
			return models.User{}, errors.Join(ErrValidation, errors.New("a senha deve ter no mínimo 8 caracteres"))
		}
		hash, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
		if err != nil {
			return models.User{}, err
		}
		passwordHash = string(hash)
	}

	return s.repo.UpdateUser(models.User{
		ID:           id,
		Name:         name,
		Email:        email,
		Role:         role,
		Status:       status,
		PasswordHash: passwordHash,
		CreatedAt:    existing.CreatedAt,
	})
}

func (s *UserService) Delete(id string) error {
	return s.repo.DeleteUser(id)
}