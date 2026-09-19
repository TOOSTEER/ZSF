package tasks

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"quest-backend/internal/crossword"
)

type Handler struct {
	repo          *Repository
	crosswordRepo *crossword.Repository
}

func NewHandler(repo *Repository, crosswordRepo *crossword.Repository) *Handler {
	return &Handler{repo: repo, crosswordRepo: crosswordRepo}
}

// GetByID handles GET /api/tasks/:id.
func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid id")
		return
	}

	task, err := h.repo.GetByID(r.Context(), id)
	if err == ErrNotFound {
		writeJSONError(w, http.StatusNotFound, "task not found")
		return
	}
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "internal error")
		return
	}

	writeJSON(w, http.StatusOK, task)
}

// GetCrossword handles GET /api/tasks/:id/crossword.
//
// Задание считается «кроссвордным», если у него есть crossword_subtype
// (classic / circles / wordsearch / text). Поле type при этом может быть
// 'point' (точки РК3..ССФРК) или 'crossword' — оба варианта допустимы.
func (h *Handler) GetCrossword(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid id")
		return
	}

	task, err := h.repo.GetByID(r.Context(), id)
	if err == ErrNotFound {
		writeJSONError(w, http.StatusNotFound, "task not found")
		return
	}
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "internal error")
		return
	}

	// ИСПРАВЛЕНО: пропускаем и точки (type='point'), у которых задан
	// crossword_subtype, и полноценные 'crossword'.
	if task.CrosswordSubtype == nil || *task.CrosswordSubtype == "" {
		writeJSONError(w, http.StatusBadRequest, "task has no crossword content")
		return
	}

	questions, err := h.crosswordRepo.GetQuestions(r.Context(), id)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "internal error")
		return
	}
	if questions == nil {
		questions = []crossword.Question{}
	}

	resp := crossword.StructureResponse{
		TaskID:    task.ID,
		Title:     task.Title,
		Subtype:   *task.CrosswordSubtype,
		Questions: questions,
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
