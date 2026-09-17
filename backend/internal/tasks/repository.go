package tasks

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrNotFound = errors.New("task not found")

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

// GetByID fetches a single task.
func (r *Repository) GetByID(ctx context.Context, id int) (*Task, error) {
	var t Task
	err := r.pool.QueryRow(ctx, `
		SELECT id, title, type, points_percent, order_index, crossword_subtype, created_at
		FROM tasks WHERE id = $1`, id,
	).Scan(&t.ID, &t.Title, &t.Type, &t.PointsPercent, &t.OrderIndex, &t.CrosswordSubtype, &t.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &t, nil
}

// GetAll returns every task ordered for display.
func (r *Repository) GetAll(ctx context.Context) ([]Task, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, title, type, points_percent, order_index, crossword_subtype, created_at
		FROM tasks ORDER BY order_index`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []Task
	for rows.Next() {
		var t Task
		if err := rows.Scan(&t.ID, &t.Title, &t.Type, &t.PointsPercent, &t.OrderIndex, &t.CrosswordSubtype, &t.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, t)
	}
	return out, rows.Err()
}
