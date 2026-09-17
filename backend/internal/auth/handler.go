package auth

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"strings"

	"golang.org/x/crypto/bcrypt"

	"quest-backend/internal/teams"
	"quest-backend/internal/users"
)

type Handler struct {
	usersRepo *users.Repository
	teamsRepo *teams.Repository
	cfg       Config
}

func NewHandler(usersRepo *users.Repository, teamsRepo *teams.Repository, cfg Config) *Handler {
	return &Handler{usersRepo: usersRepo, teamsRepo: teamsRepo, cfg: cfg}
}

// RegisterRequest is the body of POST /api/auth/register.
type RegisterRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	GroupName string `json:"group_name"`
	Role      string `json:"role"` // "host" | "member"
	Password  string `json:"password,omitempty"` // required for role="host"
}

// AuthResponse is returned by both register and login.
type AuthResponse struct {
	Token string      `json:"token"`
	User  users.User  `json:"user"`
	Team  *teams.Team `json:"team,omitempty"`
}

// LoginRequest is the body of POST /api/auth/login (host only).
type LoginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

// Register handles POST /api/auth/register.
func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid body")
		return
	}

	req.FirstName = strings.TrimSpace(req.FirstName)
	req.LastName = strings.TrimSpace(req.LastName)
	req.GroupName = strings.TrimSpace(req.GroupName)

	if req.FirstName == "" || req.LastName == "" || req.GroupName == "" {
		writeJSONError(w, http.StatusBadRequest, "first_name, last_name and group_name are required")
		return
	}
	if req.Role != "host" && req.Role != "member" {
		writeJSONError(w, http.StatusBadRequest, "role must be 'host' or 'member'")
		return
	}

	ctx := r.Context()

	if req.Role == "member" {
		teamID, err := h.teamsRepo.LeastPopulatedTeamID(ctx)
		if err != nil {
			writeJSONError(w, http.StatusInternalServerError, "internal error")
			return
		}

		login, err := h.generateUniqueLogin(ctx)
		if err != nil {
			writeJSONError(w, http.StatusInternalServerError, "internal error")
			return
		}

		u := &users.User{
			FirstName: req.FirstName,
			LastName:  req.LastName,
			GroupName: req.GroupName,
			Role:      "member",
			TeamID:    teamID,
			Login:     login,
		}

		id, err := h.usersRepo.Create(ctx, u)
		if err != nil {
			writeJSONError(w, http.StatusInternalServerError, "internal error")
			return
		}
		u.ID = id

		token, err := GenerateToken(h.cfg, u.ID, u.Role, u.TeamID)
		if err != nil {
			writeJSONError(w, http.StatusInternalServerError, "internal error")
			return
		}

		team, _ := h.teamsRepo.GetByID(ctx, teamID)
		writeJSON(w, http.StatusCreated, AuthResponse{Token: token, User: *u, Team: team})
		return
	}

	// role == "host"
	if req.Password == "" {
		writeJSONError(w, http.StatusBadRequest, "password is required for role='host'")
		return
	}

	// Найти команду без назначенного водящего.
	// TODO: уточнить, если нужен явный выбор команды для водящего вместо автоназначения
	teamID, err := h.firstTeamWithoutHost(ctx)
	if err != nil {
		writeJSONError(w, http.StatusConflict, "all teams already have a host assigned")
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "internal error")
		return
	}

	login, err := h.generateUniqueLogin(ctx)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "internal error")
		return
	}

	u := &users.User{
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		GroupName:    req.GroupName,
		Role:         "host",
		TeamID:       teamID,
		Login:        login,
		PasswordHash: string(hash),
	}

	id, err := h.usersRepo.Create(ctx, u)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "internal error")
		return
	}
	u.ID = id

	if _, err := h.teamsRepo.SetHost(ctx, teamID, id); err != nil {
		writeJSONError(w, http.StatusInternalServerError, "internal error")
		return
	}

	token, err := GenerateToken(h.cfg, u.ID, u.Role, u.TeamID)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "internal error")
		return
	}

	team, _ := h.teamsRepo.GetByID(ctx, teamID)
	writeJSON(w, http.StatusCreated, AuthResponse{Token: token, User: *u, Team: team})
}

// Login handles POST /api/auth/login (host only).
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid body")
		return
	}

	ctx := r.Context()
	u, err := h.usersRepo.GetByLogin(ctx, req.Login)
	if err != nil {
		writeJSONError(w, http.StatusUnauthorized, "invalid login or password")
		return
	}
	if u.Role != "host" {
		writeJSONError(w, http.StatusUnauthorized, "invalid login or password")
		return
	}
	if bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(req.Password)) != nil {
		writeJSONError(w, http.StatusUnauthorized, "invalid login or password")
		return
	}

	token, err := GenerateToken(h.cfg, u.ID, u.Role, u.TeamID)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "internal error")
		return
	}

	team, _ := h.teamsRepo.GetByID(ctx, u.TeamID)
	writeJSON(w, http.StatusOK, AuthResponse{Token: token, User: *u, Team: team})
}

// generateUniqueLogin creates a short random login like "u48213" and
// retries on the rare collision against existing logins.
func (h *Handler) generateUniqueLogin(ctx context.Context) (string, error) {
	for i := 0; i < 10; i++ {
		n, err := rand.Int(rand.Reader, big.NewInt(90000))
		if err != nil {
			return "", err
		}
		login := fmt.Sprintf("u%05d", n.Int64()+10000)

		exists, err := h.usersRepo.LoginExists(ctx, login)
		if err != nil {
			return "", err
		}
		if !exists {
			return login, nil
		}
	}
	return "", fmt.Errorf("could not generate unique login")
}

func (h *Handler) firstTeamWithoutHost(ctx context.Context) (int, error) {
	allTeams, err := h.teamsRepo.GetAll(ctx)
	if err != nil {
		return 0, err
	}
	for _, t := range allTeams {
		if t.HostID == nil {
			return t.ID, nil
		}
	}
	return 0, fmt.Errorf("no team without host")
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeJSONError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}
