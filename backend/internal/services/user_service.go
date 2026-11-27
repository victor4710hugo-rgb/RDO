package services

import (
<<<<<<< HEAD
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"os"
	"strings"

	"golang.org/x/crypto/bcrypt"

	"github.com/victor4710hugo-rgb/RDO/backend/internal/models"
	"github.com/victor4710hugo-rgb/RDO/backend/internal/repositories"
)

const (
	RoleAdmin   = "admin"
	RoleManager = "gerente"
	RoleUser    = "user"
)

var (
	ErrInvalidInput     = errors.New("dados de usuário inválidos")
	ErrInvalidRole      = errors.New("perfil inválido: use admin, gerente ou user")
	ErrEmailAlreadyUsed = errors.New("e-mail já cadastrado")
	ErrInvalidLogin     = errors.New("credenciais inválidas")
)

type CreateUserInput struct {
	Name     string
	Email    string
	Password string
	Role     string
}

type UpdateUserInput struct {
	Name   string
	Email  string
	Role   string
	Active bool
}

type UserService struct {
	repo repositories.UserRepository
}

func NewUserService(repo repositories.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func normalizeAndValidateRole(role string) (string, error) {
	r := strings.ToLower(strings.TrimSpace(role))
	if r == "" {
		return RoleUser, nil
	}
	switch r {
	case RoleAdmin, RoleManager, RoleUser:
		return r, nil
	default:
		return "", ErrInvalidRole
	}
}

func GetAdminInitialPassword() string {
	if pwd := strings.TrimSpace(os.Getenv("ADMIN_INITIAL_PASSWORD")); pwd != "" {
		return pwd
	}
	// Fallback apenas para ambiente local/dev; defina ADMIN_INITIAL_PASSWORD em produção
	return "Admin@123456"
}

func (s *UserService) Create(ctx context.Context, input CreateUserInput) (*models.User, error) {
	name := strings.TrimSpace(input.Name)
	email := strings.ToLower(strings.TrimSpace(input.Email))
	if name == "" || email == "" || len(input.Password) < 6 {
		return nil, ErrInvalidInput
	}

	role, err := normalizeAndValidateRole(input.Role)
	if err != nil {
		return nil, err
	}

	if _, err := s.repo.GetByEmail(ctx, email); err == nil {
		return nil, ErrEmailAlreadyUsed
	} else if !errors.Is(err, repositories.ErrUserNotFound) {
		return nil, err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	id, err := generateUserID()
	if err != nil {
		return nil, err
	}

	user := &models.User{
		ID:           id,
		Name:         name,
		Email:        email,
		PasswordHash: string(hash),
		Role:         role,
		Active:       true,
	}

	if err := s.repo.Create(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}

// EnsureSeedUsers cria os usuários iniciais de admin e gerente caso ainda não existam.
func (s *UserService) EnsureSeedUsers(ctx context.Context) error {
	adminEmail := strings.ToLower(strings.TrimSpace(os.Getenv("ADMIN_EMAIL")))
	if adminEmail == "" {
		adminEmail = "admin@rdo.local"
	}

	if _, err := s.repo.GetByEmail(ctx, adminEmail); errors.Is(err, repositories.ErrUserNotFound) {
		_, err := s.Create(ctx, CreateUserInput{
			Name:     "Administrador",
			Email:    adminEmail,
			Password: GetAdminInitialPassword(),
			Role:     RoleAdmin,
		})
		if err != nil && !errors.Is(err, ErrEmailAlreadyUsed) {
			return err
		}
	} else if err != nil {
		return err
	}

	managerEmail := strings.ToLower(strings.TrimSpace(os.Getenv("MANAGER_EMAIL")))
	if managerEmail == "" {
		managerEmail = "gerente@rdo.local"
	}
	managerPassword := strings.TrimSpace(os.Getenv("MANAGER_INITIAL_PASSWORD"))
	if managerPassword == "" {
		managerPassword = "Gerente@123456"
	}

	if _, err := s.repo.GetByEmail(ctx, managerEmail); errors.Is(err, repositories.ErrUserNotFound) {
		_, err := s.Create(ctx, CreateUserInput{
			Name:     "Gerente Operacional",
			Email:    managerEmail,
			Password: managerPassword,
			Role:     RoleManager,
		})
		if err != nil && !errors.Is(err, ErrEmailAlreadyUsed) {
			return err
		}
	} else if err != nil {
		return err
	}

	return nil
}

func (s *UserService) GetByID(ctx context.Context, id string) (*models.User, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *UserService) List(ctx context.Context) ([]models.User, error) {
	return s.repo.List(ctx)
}

func (s *UserService) Update(ctx context.Context, id string, input UpdateUserInput) (*models.User, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	name := strings.TrimSpace(input.Name)
	email := strings.ToLower(strings.TrimSpace(input.Email))
	if name == "" || email == "" {
		return nil, ErrInvalidInput
	}

	role, err := normalizeAndValidateRole(input.Role)
	if err != nil {
		return nil, err
	}

	user.Name = name
	user.Email = email
	user.Role = role
	user.Active = input.Active

	if err := s.repo.Update(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s *UserService) Authenticate(ctx context.Context, email, password string) (*models.User, error) {
	normalizedEmail := strings.ToLower(strings.TrimSpace(email))
	user, err := s.repo.GetByEmail(ctx, normalizedEmail)
	if err != nil {
		if errors.Is(err, repositories.ErrUserNotFound) {
			return nil, ErrInvalidLogin
		}
		return nil, err
	}

	if !user.Active {
		return nil, ErrInvalidLogin
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, ErrInvalidLogin
	}

	return user, nil
}

func generateUserID() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
=======
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
>>>>>>> origin/main
}