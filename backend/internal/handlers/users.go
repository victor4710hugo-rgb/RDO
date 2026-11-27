package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/victor4710hugo-rgb/RDO/backend/internal/models"
	"github.com/victor4710hugo-rgb/RDO/backend/internal/repositories"
	"github.com/victor4710hugo-rgb/RDO/backend/internal/services"
)

// UsersHandler traduz requisições HTTP de /api/v1/users para chamadas ao UserService.
type UsersHandler struct {
	service *services.UserService
}

// NewUsersHandler cria o handler de usuários recebendo o serviço injetado.
func NewUsersHandler(service *services.UserService) *UsersHandler {
	return &UsersHandler{service: service}
}

// writeJSON padroniza o Content-Type, o status HTTP e a serialização JSON.
func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

// writeError retorna erros da API no formato padrão {"erro": "..."}.
func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"erro": msg})
}

// mapUserError converte erros de domínio/repositório no status HTTP correspondente.
func mapUserError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, services.ErrValidation):
		writeError(w, http.StatusBadRequest, err.Error())
	case errors.Is(err, services.ErrForbiddenAction):
		// Novo (Issue #3): devolve 403 para tentativa de autoexclusão ou autorebaixamento.
		writeError(w, http.StatusForbidden, err.Error())
	case errors.Is(err, repositories.ErrEmailExists):
		writeError(w, http.StatusConflict, err.Error())
	case errors.Is(err, repositories.ErrUserNotFound):
		writeError(w, http.StatusNotFound, err.Error())
	default:
		writeError(w, http.StatusInternalServerError, "erro interno ao processar usuário")
	}
}

// List atende GET /api/v1/users e devolve todos os usuários ativos.
func (h *UsersHandler) List(w http.ResponseWriter, r *http.Request) {
	users, err := h.service.List()
	if err != nil {
		mapUserError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, users)
}

// GetByID atende GET /api/v1/users/{id} e devolve um usuário pelo ID da rota.
func (h *UsersHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	user, err := h.service.GetByID(r.PathValue("id"))
	if err != nil {
		mapUserError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, user)
}

// Create atende POST /api/v1/users e retorna 201 Created.
func (h *UsersHandler) Create(w http.ResponseWriter, r *http.Request) {
	var in models.CreateUserInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeError(w, http.StatusBadRequest, "JSON inválido")
		return
	}
	user, err := h.service.Create(in)
	if err != nil {
		mapUserError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, user)
}

// Update atende PUT /api/v1/users/{id}, repassando o ID do usuário autenticado (actorID).
func (h *UsersHandler) Update(w http.ResponseWriter, r *http.Request) {
	var in models.UpdateUserInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeError(w, http.StatusBadRequest, "JSON inválido")
		return
	}

	// Novo (Issue #3): extrai o ID do usuário logado das claims do JWT para bloquear autorebaixamento.
	var actorID string
	if claims, ok := ClaimsFromContext(r.Context()); ok {
		actorID = claims.UserID
	}

	user, err := h.service.Update(actorID, r.PathValue("id"), in)
	if err != nil {
		mapUserError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, user)
}

// Delete atende DELETE /api/v1/users/{id}, repassando o actorID para impedir autoexclusão.
func (h *UsersHandler) Delete(w http.ResponseWriter, r *http.Request) {
	// Novo (Issue #3): extrai o ID do usuário logado das claims do JWT para bloquear autoexclusão.
	var actorID string
	if claims, ok := ClaimsFromContext(r.Context()); ok {
		actorID = claims.UserID
	}

	if err := h.service.Delete(actorID, r.PathValue("id")); err != nil {
		mapUserError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
