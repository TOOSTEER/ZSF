package crossword

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"quest-backend/internal/auth/reqctx"
)

type Handler struct {
	repo *Repository
}

func NewHandler(repo *Repository) *Handler {
	return &Handler{repo: repo}
}

// Check handles POST /api/crossword/:task_id/check (member only).
func (h *Handler) Check(w http.ResponseWriter, r *http.Request) {
	claims, ok := reqctx.FromContext(r.Context())
	if !ok {
		writeJSONError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	taskID, err := strconv.Atoi(chi.URLParam(r, "task_id"))
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid task_id")
		return
	}

	var req CheckRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid body")
		return
	}

	resp, err := h.repo.Check(r.Context(), taskID, claims.TeamID, claims.UserID, req.Answers)
	if err == ErrAlreadyCompleted {
		writeJSONError(w, http.StatusConflict, "task already completed for this team")
		return
	}
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "internal error")
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeJSONError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}
