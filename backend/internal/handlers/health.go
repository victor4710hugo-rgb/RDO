package handlers

import "net/http"

// Health atende GET /api/v1/health retornando 200 OK em JSON padronizado.
func Health(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}
