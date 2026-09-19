package tasks

import "time"

// Task represents a row in the tasks table.
type Task struct {
	ID               int       `json:"id"`
	Title            string    `json:"title"`
	Type             string    `json:"type"` // point | extra | crossword
	PointsPercent    int       `json:"points_percent"`
	OrderIndex       int       `json:"order_index"`
	CrosswordSubtype *string   `json:"crossword_subtype,omitempty"`
	CreatedAt        time.Time `json:"created_at"`
}
