package models

import "time"

// Expense represents a single transaction record
type Expense struct {
	ID          int       `json:"id"`
	CategoryID  int       `json:"category_id"`
	Amount      float64   `json:"amount"`
	ExpenseDate time.Time `json:"expense_date"`
	Remarks     string    `json:"remarks"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// Joined field for display
	CategoryName string `json:"category_name,omitempty"`
}

// ExpenseRequest represents the payload to record a new expense
type ExpenseRequest struct {
	CategoryID  int     `json:"category_id"`
	Amount      float64 `json:"amount"`
	ExpenseDate string  `json:"expense_date"` // String from form: YYYY-MM-DD
	Remarks     string  `json:"remarks"`
}

// ExpenseFilter represents criteria for searching expenses
type ExpenseFilter struct {
	StartDate  string `json:"start_date"`
	EndDate    string `json:"end_date"`
	CategoryID int    `json:"category_id"`
}
