package handlers

import "net/http"

// SyncHandler atende os endpoints iniciais de sincronização (pull/push).
type SyncHandler struct{}

// NewSyncHandler cria o handler de sincronização.
func NewSyncHandler() *SyncHandler {
	return &SyncHandler{}
}

// Pull atende GET /api/v1/sync/pull retornando 200 OK em JSON padronizado.
func (h *SyncHandler) Pull(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

// Push atende POST /api/v1/sync/push retornando 200 OK em JSON padronizado.
func (h *SyncHandler) Push(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}
