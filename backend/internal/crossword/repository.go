package crossword

import (
	"context"
	"errors"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

// ErrAlreadyCompleted is returned when a team already got points for a task.
var ErrAlreadyCompleted = errors.New("task already completed for this team")

// CrosswordPoints is the fixed reward for solving any crossword.
// Проценты фиксированы (см. раздел 5 ТЗ).
const CrosswordPoints = 3

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

// GetQuestions returns the clue list for a crossword task, ordered by number.
func (r *Repository) GetQuestions(ctx context.Context, taskID int) ([]Question, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, task_id, question_number, question_text, COALESCE(direction, ''), COALESCE(answer_length, 0)
		FROM crossword_questions
		WHERE task_id = $1
		ORDER BY question_number`, taskID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []Question
	for rows.Next() {
		var q Question
		if err := rows.Scan(&q.ID, &q.TaskID, &q.QuestionNumber, &q.QuestionText, &q.Direction, &q.AnswerLength); err != nil {
			return nil, err
		}
		out = append(out, q)
	}
	return out, rows.Err()
}

// getAnswers loads the correct-answer map (question_number -> answer) for a task.
func (r *Repository) getAnswers(ctx context.Context, taskID int) (map[int]string, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT question_number, correct_answer FROM crossword_answers WHERE task_id = $1`, taskID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make(map[int]string)
	for rows.Next() {
		var num int
		var answer string
		if err := rows.Scan(&num, &answer); err != nil {
			return nil, err
		}
		out[num] = answer
	}
	return out, rows.Err()
}

func normalize(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}

// Check compares submitted answers against crossword_answers. If every
// answer is correct it awards CrosswordPoints to the team inside a
// transaction guarded by the team_progress unique constraint, so a
// concurrent double-submit safely results in ErrAlreadyCompleted instead
// of double crediting.
func (r *Repository) Check(ctx context.Context, taskID, teamID, userID int, submitted []CheckAnswerInput) (*CheckResponse, error) {
	correct, err := r.getAnswers(ctx, taskID)
	if err != nil {
		return nil, err
	}

	var wrong []int
	for _, a := range submitted {
		expected, ok := correct[a.QuestionNumber]
		// If there is no reference answer yet (content not filled in),
		// treat it as unresolved/wrong rather than silently passing.
		// TODO: уточнить поведение, когда правильные ответы ещё не внесены
		if !ok || normalize(expected) != normalize(a.Answer) {
			wrong = append(wrong, a.QuestionNumber)
		}
	}

	if len(wrong) > 0 {
		return &CheckResponse{Success: true, Awarded: false, WrongNumbers: wrong}, nil
	}

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx) //nolint:errcheck

	var alreadyDone bool
	err = tx.QueryRow(ctx,
		`SELECT is_completed FROM team_progress WHERE team_id = $1 AND task_id = $2`,
		teamID, taskID,
	).Scan(&alreadyDone)
	if err == nil && alreadyDone {
		return nil, ErrAlreadyCompleted
	}

	_, err = tx.Exec(ctx, `
		INSERT INTO team_progress (team_id, task_id, is_completed, points_awarded, completed_by, completed_at)
		VALUES ($1, $2, TRUE, $3, $4, NOW())
		ON CONFLICT (team_id, task_id) DO UPDATE
			SET is_completed = TRUE, points_awarded = EXCLUDED.points_awarded,
			    completed_by = EXCLUDED.completed_by, completed_at = NOW()
			WHERE team_progress.is_completed = FALSE`,
		teamID, taskID, CrosswordPoints, userID,
	)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return &CheckResponse{Success: true, Awarded: true}, nil
}
