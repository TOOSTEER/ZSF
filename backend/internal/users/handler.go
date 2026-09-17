package users

import (
	"encoding/json"
	"net/http"

	"quest-backend/internal/auth/reqctx"
	"quest-backend/internal/teams"
)

type Handler struct {
	repo      *Repository
	teamsRepo *teams.Repository
}

func NewHandler(repo *Repository, teamsRepo *teams.Repository) *Handler {
	return &Handler{repo: repo, teamsRepo: teamsRepo}
}

// MeResponse is the payload for GET /api/me.
type MeResponse struct {
	User User       `json:"user"`
	Team *teams.Team `json:"team,omitempty"`
}

// Me handles GET /api/me.
func (h *Handler) Me(w http.ResponseWriter, r *http.Request) {
	claims, ok := reqctx.FromContext(r.Context())
	if !ok {
		writeJSONError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	u, err := h.repo.GetByID(r.Context(), claims.UserID)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "internal error")
		return
	}

	resp := MeResponse{User: *u}

	team, err := h.teamsRepo.GetByID(r.Context(), u.TeamID)
	if err == nil {
		resp.Team = team
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
