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

type RegisterRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	GroupName string `json:"group_name"`
	Role      string `json:"role"`
	Password  string `json:"password,omitempty"`
}

type AuthResponse struct {
	Token string      `json:"token"`
	User  users.User  `json:"user"`
	Team  *teams.Team `json:"team,omitempty"`
}

// LoginRequest — вход водящего по логину+паролю (латиница).
// Логины: host_blue, host_red, host_orange, host_green, host_yellow.
type LoginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

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
	// Регистрация доступна только участникам. Водящие создаются админом.
	if req.Role != "member" {
		writeJSONError(w, http.StatusBadRequest, "registration is only for members")
		return
	}

	ctx := r.Context()

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
}

// Login — вход водящего по логину+паролю.
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid body")
		return
	}

	req.Login = strings.TrimSpace(req.Login)
	if req.Login == "" || req.Password == "" {
		writeJSONError(w, http.StatusBadRequest, "login and password are required")
		return
	}

	ctx := r.Context()
	u, err := h.usersRepo.GetByLogin(ctx, req.Login)
	if err != nil {
		writeJSONError(w, http.StatusUnauthorized, "invalid credentials")
		return
	}
	if u.Role != "host" {
		writeJSONError(w, http.StatusUnauthorized, "invalid credentials")
		return
	}
	if u.PasswordHash == "" {
		writeJSONError(w, http.StatusUnauthorized, "invalid credentials")
		return
	}
	if bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(req.Password)) != nil {
		writeJSONError(w, http.StatusUnauthorized, "invalid credentials")
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

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeJSONError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}
