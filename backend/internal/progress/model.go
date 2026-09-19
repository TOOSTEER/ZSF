package progress

import "time"

// TeamProgress represents a row in the team_progress table.
type TeamProgress struct {
	ID            int        `json:"id"`
	TeamID        int        `json:"team_id"`
	TaskID        int        `json:"task_id"`
	IsCompleted   bool       `json:"is_completed"`
	PointsAwarded int        `json:"points_awarded"`
	CompletedBy   *int       `json:"completed_by,omitempty"`
	CompletedAt   *time.Time `json:"completed_at,omitempty"`
}

// AwardRequest is the body of POST /api/progress/award.
type AwardRequest struct {
	TeamID int `json:"team_id"`
	TaskID int `json:"task_id"`
	Points int `json:"points"`
}
