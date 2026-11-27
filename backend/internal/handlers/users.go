package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/victor4710hugo-rgb/RDO/backend/internal/models"
	"github.com/victor4710hugo-rgb/RDO/backend/internal/repositories"
	"github.com/victor4710hugo-rgb/RDO/backend/internal/services"
)

type UsersHandler struct {
	service *services.UserService
}

func NewUsersHandler(service *services.UserService) *UsersHandler {
	return &UsersHandler{service: service}
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"erro": msg})
}

func mapUserError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, services.ErrValidation):
		writeError(w, http.StatusBadRequest, err.Error())
	case errors.Is(err, repositories.ErrEmailExists):
		writeError(w, http.StatusConflict, err.Error())
	case errors.Is(err, repositories.ErrUserNotFound):
		writeError(w, http.StatusNotFound, err.Error())
	default:
		writeError(w, http.StatusInternalServerError, "erro interno ao processar usuário")
	}
}

func (h *UsersHandler) List(w http.ResponseWriter, r *http.Request) {
	users, err := h.service.List()
	if err != nil {
		mapUserError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, users)
}

func (h *UsersHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	user, err := h.service.GetByID(r.PathValue("id"))
	if err != nil {
		mapUserError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, user)
}

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

func (h *UsersHandler) Update(w http.ResponseWriter, r *http.Request) {
	var in models.UpdateUserInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeError(w, http.StatusBadRequest, "JSON inválido")
		return
	}
	user, err := h.service.Update(r.PathValue("id"), in)
	if err != nil {
		mapUserError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, user)
}

func (h *UsersHandler) Delete(w http.ResponseWriter, r *http.Request) {
	if err := h.service.Delete(r.PathValue("id")); err != nil {
		mapUserError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}