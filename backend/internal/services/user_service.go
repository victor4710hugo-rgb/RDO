package services

import (
	"errors"
	"net/mail"
	"strings"

	"github.com/google/uuid"
	"github.com/victor4710hugo-rgb/RDO/backend/internal/models"
	"github.com/victor4710hugo-rgb/RDO/backend/internal/repositories"
	"golang.org/x/crypto/bcrypt"
)

// ErrValidation marca qualquer erro de entrada inválida (handler responde 400).
// ErrForbiddenAction marca ações bloqueadas por regra de segurança/RBAC (handler responde 403).
var (
	ErrValidation      = errors.New("dados inválidos")
	ErrForbiddenAction = errors.New("ação não permitida")
)

// minPasswordLen é o tamanho mínimo da senha (criação e troca).
const minPasswordLen = 8

// UserService concentra as regras de negócio de usuários (sem SQL).
type UserService struct {
	repo *repositories.Repository
}

// NewUserService cria o serviço usando o repositório compartilhado.
func NewUserService(repo *repositories.Repository) *UserService {
	return &UserService{repo: repo}
}

// validationError padroniza erros de validação com errors.Join.
func validationError(msg string) error {
	return errors.Join(ErrValidation, errors.New(msg))
}

// forbiddenError padroniza erros de ação proibida (403) com errors.Join.
func forbiddenError(msg string) error {
	return errors.Join(ErrForbiddenAction, errors.New(msg))
}

// isValidRole aceita apenas os papéis existentes no sistema.
func isValidRole(role string) bool {
	return role == "admin" || role == "gerente"
}

// isValidStatus aceita apenas os status permitidos pelo banco.
func isValidStatus(status string) bool {
	return status == "active" || status == "inactive" || status == "blocked"
}

// validateEmail valida o formato com net/mail e rejeita formas como
// "Nome <a@b.com>". Usada na criação e na atualização.
func validateEmail(email string) error {
	addr, err := mail.ParseAddress(email)
	if err != nil || addr.Address != email {
		return validationError("e-mail inválido")
	}
	return nil
}

// validateIdentity valida nome, e-mail e papel (regras comuns a criar/atualizar).
func validateIdentity(name, email, role string) error {
	if name == "" {
		return validationError("nome é obrigatório")
	}
	if err := validateEmail(email); err != nil {
		return err
	}
	if !isValidRole(role) {
		return validationError("papel permitido: admin ou gerente")
	}
	return nil
}

// hashPassword valida o tamanho mínimo e gera o hash bcrypt da senha.
func hashPassword(password string) (string, error) {
	if len(password) < minPasswordLen {
		return "", validationError("a senha deve ter no mínimo 8 caracteres")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// List devolve todos os usuários ativos.
func (s *UserService) List() ([]models.User, error) {
	return s.repo.ListUsers()
}

// GetByID busca um usuário pelo identificador.
func (s *UserService) GetByID(id string) (models.User, error) {
	return s.repo.GetUserByID(id)
}

// Create valida a entrada, gera o hash da senha e cria o usuário como "active".
func (s *UserService) Create(in models.CreateUserInput) (models.User, error) {
	name := strings.TrimSpace(in.Name)
	email := strings.ToLower(strings.TrimSpace(in.Email))
	role := strings.ToLower(strings.TrimSpace(in.Role))

	if err := validateIdentity(name, email, role); err != nil {
		return models.User{}, err
	}
	hash, err := hashPassword(in.Password)
	if err != nil {
		return models.User{}, err
	}

	return s.repo.CreateUser(models.User{
		ID:           uuid.NewString(),
		Name:         name,
		Email:        email,
		Role:         role,
		Status:       "active",
		PasswordHash: hash,
	})
}

// Update atualiza nome, e-mail, papel e status. Regras:
//   - papel é OBRIGATÓRIO (vazio é rejeitado, nunca preservado em silêncio);
//   - o usuário autenticado (actorID) não pode alterar o próprio papel (403);
//   - status vazio preserva o atual; se enviado, deve ser active|inactive|blocked;
//   - senha vazia preserva o hash atual; se enviada, é validada e re-hasheada.
func (s *UserService) Update(actorID, id string, in models.UpdateUserInput) (models.User, error) {
	existing, err := s.repo.GetUserByID(id)
	if err != nil {
		return models.User{}, err
	}

	name := strings.TrimSpace(in.Name)
	email := strings.ToLower(strings.TrimSpace(in.Email))
	role := strings.ToLower(strings.TrimSpace(in.Role))
	status := strings.ToLower(strings.TrimSpace(in.Status))
	if status == "" {
		status = existing.Status
	}

	if err := validateIdentity(name, email, role); err != nil {
		return models.User{}, err
	}
	if !isValidStatus(status) {
		return models.User{}, validationError("status permitido: active, inactive ou blocked")
	}
	// Regra de segurança RBAC: impede que o usuário autenticado rebaixe ou troque o próprio papel.
	if actorID != "" && actorID == id && role != existing.Role {
		return models.User{}, forbiddenError("não é permitido alterar o próprio papel")
	}

	passwordHash := existing.PasswordHash
	if in.Password != "" {
		if passwordHash, err = hashPassword(in.Password); err != nil {
			return models.User{}, err
		}
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

// Delete remove o usuário pelo identificador, impedindo que o usuário autenticado exclua a si mesmo.
func (s *UserService) Delete(actorID, id string) error {
	if actorID != "" && actorID == id {
		return forbiddenError("não é permitido excluir o próprio usuário")
	}
	return s.repo.DeleteUser(id)
}
