package progress

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrAlreadyCompleted = errors.New("task already completed for this team")
	ErrTaskNotFound     = errors.New("task not found")
)

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

// Award credits points to a team for a task. It fails with
// ErrAlreadyCompleted if the task was already awarded to that team
// (повторное начисление запрещено — раздел 5 ТЗ).
func (r *Repository) Award(ctx context.Context, teamID, taskID, points, hostUserID int) (*TeamProgress, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx) //nolint:errcheck

	var exists bool
	if err := tx.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM tasks WHERE id = $1)`, taskID).Scan(&exists); err != nil {
		return nil, err
	}
	if !exists {
		return nil, ErrTaskNotFound
	}

	var alreadyDone bool
	err = tx.QueryRow(ctx,
		`SELECT is_completed FROM team_progress WHERE team_id = $1 AND task_id = $2`,
		teamID, taskID,
	).Scan(&alreadyDone)
	if err == nil && alreadyDone {
		return nil, ErrAlreadyCompleted
	}

	var tp TeamProgress
	err = tx.QueryRow(ctx, `
		INSERT INTO team_progress (team_id, task_id, is_completed, points_awarded, completed_by, completed_at)
		VALUES ($1, $2, TRUE, $3, $4, NOW())
		ON CONFLICT (team_id, task_id) DO UPDATE
			SET is_completed = TRUE, points_awarded = EXCLUDED.points_awarded,
			    completed_by = EXCLUDED.completed_by, completed_at = NOW()
			WHERE team_progress.is_completed = FALSE
		RETURNING id, team_id, task_id, is_completed, points_awarded, completed_by, completed_at`,
		teamID, taskID, points, hostUserID,
	).Scan(&tp.ID, &tp.TeamID, &tp.TaskID, &tp.IsCompleted, &tp.PointsAwarded, &tp.CompletedBy, &tp.CompletedAt)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return &tp, nil
}
