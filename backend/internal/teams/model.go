package teams

// Team represents a row in the teams table.
type Team struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Color  string `json:"color"`
	HostID *int   `json:"host_id,omitempty"`
}

// ProgressItem describes a single task's completion state for a team.
type ProgressItem struct {
	TaskID           int    `json:"task_id"`
	Title            string `json:"title"`
	Type             string `json:"type"` // point | extra | crossword
	PointsPercent    int    `json:"points_percent"`
	CrosswordSubtype string `json:"crossword_subtype,omitempty"`
	IsCompleted      bool   `json:"is_completed"`
	PointsAwarded    int    `json:"points_awarded"`
}

// ProgressResponse is the payload for GET /api/teams/:id/progress.
type ProgressResponse struct {
	TeamID        int            `json:"team_id"`
	TeamName      string         `json:"team_name"`
	TeamColor     string         `json:"team_color"`
	TotalPercent  int            `json:"total_percent"`  // raw sum, can exceed 100
	DisplayPercent int           `json:"display_percent"` // min(total, 100)
	Items         []ProgressItem `json:"items"`
}
