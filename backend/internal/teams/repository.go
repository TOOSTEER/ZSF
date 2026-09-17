package teams

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrNotFound = errors.New("team not found")

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

// GetByID fetches a single team.
func (r *Repository) GetByID(ctx context.Context, id int) (*Team, error) {
	var t Team
	err := r.pool.QueryRow(ctx,
		`SELECT id, name, color, host_id FROM teams WHERE id = $1`, id,
	).Scan(&t.ID, &t.Name, &t.Color, &t.HostID)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &t, nil
}

// GetAll returns every team, ordered by id.
func (r *Repository) GetAll(ctx context.Context) ([]Team, error) {
	rows, err := r.pool.Query(ctx, `SELECT id, name, color, host_id FROM teams ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []Team
	for rows.Next() {
		var t Team
		if err := rows.Scan(&t.ID, &t.Name, &t.Color, &t.HostID); err != nil {
			return nil, err
		}
		out = append(out, t)
	}
	return out, rows.Err()
}

// LeastPopulatedTeamID returns the id of the team with the fewest
// members, used to distribute new participants evenly across the 5 teams.
// Ties are broken by team id (ascending), which keeps the distribution
// deterministic and simple. // TODO: уточнить, если нужна чистая случайность вместо "минимальной команды"
func (r *Repository) LeastPopulatedTeamID(ctx context.Context) (int, error) {
	query := `
		SELECT t.id
		FROM teams t
		LEFT JOIN users u ON u.team_id = t.id AND u.role = 'member'
		GROUP BY t.id
		ORDER BY COUNT(u.id) ASC, t.id ASC
		LIMIT 1`

	var id int
	err := r.pool.QueryRow(ctx, query).Scan(&id)
	return id, err
}

// SetHost assigns a host user to a team (used right after a host registers).
func (r *Repository) SetHost(ctx context.Context, teamID, hostUserID int) (*Team, error) {
	_, err := r.pool.Exec(ctx, `UPDATE teams SET host_id = $1 WHERE id = $2`, hostUserID, teamID)
	if err != nil {
		return nil, err
	}
	return r.GetByID(ctx, teamID)
}

// GetProgress returns the full task list joined with this team's
// completion state, plus the summed percentages.
func (r *Repository) GetProgress(ctx context.Context, teamID int) (*ProgressResponse, error) {
	team, err := r.GetByID(ctx, teamID)
	if err != nil {
		return nil, err
	}

	query := `
		SELECT
			t.id, t.title, t.type, t.points_percent, COALESCE(t.crossword_subtype, ''),
			COALESCE(tp.is_completed, FALSE), COALESCE(tp.points_awarded, 0)
		FROM tasks t
		LEFT JOIN team_progress tp ON tp.task_id = t.id AND tp.team_id = $1
		ORDER BY t.order_index`

	rows, err := r.pool.Query(ctx, query, teamID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	resp := &ProgressResponse{
		TeamID:    team.ID,
		TeamName:  team.Name,
		TeamColor: team.Color,
	}

	total := 0
	for rows.Next() {
		var item ProgressItem
		if err := rows.Scan(
			&item.TaskID, &item.Title, &item.Type, &item.PointsPercent, &item.CrosswordSubtype,
			&item.IsCompleted, &item.PointsAwarded,
		); err != nil {
			return nil, err
		}
		total += item.PointsAwarded
		resp.Items = append(resp.Items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	resp.TotalPercent = total
	if total > 100 {
		resp.DisplayPercent = 100
	} else {
		resp.DisplayPercent = total
	}

	return resp, nil
}
