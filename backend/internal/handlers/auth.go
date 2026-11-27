package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/victor4710hugo-rgb/RDO/backend/internal/services"
)

type contextKey string

const authClaimsKey contextKey = "auth_claims"

type AuthHandler struct {
	service *services.AuthService
}

func NewAuthHandler(service *services.AuthService) *AuthHandler {
	return &AuthHandler{service: service}
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "JSON inválido")
		return
	}

	token, user, err := h.service.Login(req.Email, req.Password)
	if err != nil {
		if errors.Is(err, services.ErrInactiveUser) {
			writeError(w, http.StatusForbidden, err.Error())
			return
		}
		writeError(w, http.StatusUnauthorized, services.ErrInvalidCredentials.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"token": token,
		"user":  user,
	})
}

func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	claims, ok := ClaimsFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "não autenticado")
		return
	}

	user, err := h.service.GetMe(claims.UserID)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "sessão inválida")
		return
	}

	writeJSON(w, http.StatusOK, user)
}

func (h *AuthHandler) RequireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := strings.TrimSpace(r.Header.Get("Authorization"))
		if !strings.HasPrefix(authHeader, "Bearer ") {
			writeError(w, http.StatusUnauthorized, "token ausente ou malformado")
			return
		}

		tokenStr := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))
		claims, err := h.service.ParseToken(tokenStr)
		if err != nil {
			writeError(w, http.StatusUnauthorized, services.ErrInvalidToken.Error())
			return
		}

		ctx := context.WithValue(r.Context(), authClaimsKey, claims)
		next(w, r.WithContext(ctx))
	}
}

func (h *AuthHandler) RequireRole(allowedRoles ...string) func(http.HandlerFunc) http.HandlerFunc {
	allowed := make(map[string]struct{}, len(allowedRoles))
	for _, role := range allowedRoles {
		allowed[strings.ToLower(strings.TrimSpace(role))] = struct{}{}
	}

	return func(next http.HandlerFunc) http.HandlerFunc {
		return h.RequireAuth(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := ClaimsFromContext(r.Context())
			if !ok {
				writeError(w, http.StatusUnauthorized, "não autenticado")
				return
			}

			if _, permitted := allowed[strings.ToLower(claims.Role)]; !permitted {
				writeError(w, http.StatusForbidden, "permissão insuficiente para este recurso")
				return
			}

			next(w, r)
		})
	}
}

func ClaimsFromContext(ctx context.Context) (*services.AuthClaims, bool) {
	claims, ok := ctx.Value(authClaimsKey).(*services.AuthClaims)
	return claims, ok
}
