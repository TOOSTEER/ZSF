package users

import "time"

// User represents a row in the users table.
type User struct {
	ID           int       `json:"id"`
	FirstName    string    `json:"first_name"`
	LastName     string    `json:"last_name"`
	GroupName    string    `json:"group_name"`
	Role         string    `json:"role"` // "host" | "member"
	TeamID       int       `json:"team_id"`
	Login        string    `json:"login"`
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"created_at"`
}
