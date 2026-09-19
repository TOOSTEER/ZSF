package users

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrNotFound = errors.New("user not found")

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

func (r *Repository) Create(ctx context.Context, u *User) (int, error) {
	query := `
		INSERT INTO users (first_name, last_name, group_name, role, team_id, login, password_hash)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id`

	var id int
	err := r.pool.QueryRow(ctx, query,
		u.FirstName, u.LastName, u.GroupName, u.Role, u.TeamID, u.Login, u.PasswordHash,
	).Scan(&id)
	return id, err
}

func (r *Repository) GetByID(ctx context.Context, id int) (*User, error) {
	query := `
		SELECT id, first_name, last_name, group_name, role, team_id, login, password_hash, created_at
		FROM users WHERE id = $1`

	var u User
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&u.ID, &u.FirstName, &u.LastName, &u.GroupName, &u.Role, &u.TeamID, &u.Login, &u.PasswordHash, &u.CreatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *Repository) GetByLogin(ctx context.Context, login string) (*User, error) {
	query := `
		SELECT id, first_name, last_name, group_name, role, team_id, login, password_hash, created_at
		FROM users WHERE login = $1`

	var u User
	err := r.pool.QueryRow(ctx, query, login).Scan(
		&u.ID, &u.FirstName, &u.LastName, &u.GroupName, &u.Role, &u.TeamID, &u.Login, &u.PasswordHash, &u.CreatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}

// GetHostByName ищет водящего по имени+фамилии (регистронезависимо).
// Используется в /api/auth/login.
func (r *Repository) GetHostByName(ctx context.Context, firstName, lastName string) (*User, error) {
	query := `
		SELECT id, first_name, last_name, group_name, role, team_id, login, password_hash, created_at
		FROM users
		WHERE LOWER(first_name) = LOWER($1)
		  AND LOWER(last_name)  = LOWER($2)
		  AND role = 'host'
		LIMIT 1`

	var u User
	err := r.pool.QueryRow(ctx, query, firstName, lastName).Scan(
		&u.ID, &u.FirstName, &u.LastName, &u.GroupName, &u.Role, &u.TeamID, &u.Login, &u.PasswordHash, &u.CreatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *Repository) LoginExists(ctx context.Context, login string) (bool, error) {
	var exists bool
	err := r.pool.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM users WHERE login = $1)`, login).Scan(&exists)
	return exists, err
}

func (r *Repository) CountMembers(ctx context.Context, teamID int) (int, error) {
	var count int
	err := r.pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM users WHERE team_id = $1 AND role = 'member'`, teamID,
	).Scan(&count)
	return count, err
}
