package services

import (
	"errors"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/victor4710hugo-rgb/RDO/backend/internal/models"
	"github.com/victor4710hugo-rgb/RDO/backend/internal/repositories"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("e-mail ou senha inválidos")
	ErrInactiveUser       = errors.New("usuário inativo ou bloqueado")
	ErrInvalidToken       = errors.New("token inválido ou expirado")
)

type AuthClaims struct {
	UserID string `json:"sub"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

type AuthService struct {
	repo   *repositories.Repository
	secret []byte
}

func NewAuthService(repo *repositories.Repository) (*AuthService, error) {
	secret := strings.TrimSpace(os.Getenv("JWT_SECRET"))
	if secret == "" {
		return nil, errors.New("JWT_SECRET não configurado no ambiente")
	}
	return &AuthService{repo: repo, secret: []byte(secret)}, nil
}

func (s *AuthService) Login(email, password string) (string, models.User, error) {
	email = strings.ToLower(strings.TrimSpace(email))
	if email == "" || password == "" {
		return "", models.User{}, ErrInvalidCredentials
	}

	user, err := s.repo.GetUserByEmail(email)
	if err != nil {
		return "", models.User{}, ErrInvalidCredentials
	}
	if user.Status != "active" {
		return "", models.User{}, ErrInactiveUser
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", models.User{}, ErrInvalidCredentials
	}

	now := time.Now().UTC()
	claims := AuthClaims{
		UserID: user.ID,
		Email:  user.Email,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   user.ID,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(8 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(s.secret)
	if err != nil {
		return "", models.User{}, err
	}

	return signed, user, nil
}

func (s *AuthService) ParseToken(tokenString string) (*AuthClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &AuthClaims{}, func(t *jwt.Token) (any, error) {
		if t.Method != jwt.SigningMethodHS256 {
			return nil, ErrInvalidToken
		}
		return s.secret, nil
	})
	if err != nil || !token.Valid {
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*AuthClaims)
	if !ok || claims.UserID == "" {
		return nil, ErrInvalidToken
	}
	return claims, nil
}

func (s *AuthService) GetMe(userID string) (models.User, error) {
	user, err := s.repo.GetUserByID(userID)
	if err != nil {
		return models.User{}, err
	}
	if user.Status != "active" {
		return models.User{}, ErrInactiveUser
	}
	return user, nil
}