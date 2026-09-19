package progress

import (
	"encoding/json"
	"net/http"

	"quest-backend/internal/auth/reqctx"
)

type Handler struct {
	repo *Repository
}

func NewHandler(repo *Repository) *Handler {
	return &Handler{repo: repo}
}

// Award handles POST /api/progress/award (host only).
func (h *Handler) Award(w http.ResponseWriter, r *http.Request) {
	claims, ok := reqctx.FromContext(r.Context())
	if !ok {
		writeJSONError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req AwardRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid body")
		return
	}

	// Host может начислять баллы только своей команде.
	if req.TeamID != claims.TeamID {
		writeJSONError(w, http.StatusForbidden, "you can only award points to your own team")
		return
	}
	if req.Points <= 0 {
		writeJSONError(w, http.StatusBadRequest, "points must be positive")
		return
	}

	tp, err := h.repo.Award(r.Context(), req.TeamID, req.TaskID, req.Points, claims.UserID)
	if err == ErrAlreadyCompleted {
		writeJSONError(w, http.StatusConflict, "task already completed for this team")
		return
	}
	if err == ErrTaskNotFound {
		writeJSONError(w, http.StatusNotFound, "task not found")
		return
	}
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "internal error")
		return
	}

	writeJSON(w, http.StatusOK, tp)
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeJSONError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}
