package handlers

import (
	"encoding/json"
	"net/http"
)

type SyncHandler struct{}

func NewSyncHandler() *SyncHandler {
	return &SyncHandler{}
}

func (h *SyncHandler) Pull(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func (h *SyncHandler) Push(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}