package models

import "time"

// Category represents an expense category
type Category struct {
	ID        int       `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	IsActive  bool      `json:"is_active" db:"is_active"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// CategoryRequest represents the request body for creating/updating a category
type CategoryRequest struct {
	Name     string `json:"name"`
	IsActive *bool  `json:"is_active,omitempty"`
}
