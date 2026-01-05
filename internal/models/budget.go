package models

import "time"

// Budget represents an annual budget for a category
type Budget struct {
	ID         int       `json:"id"`
	CategoryID int       `json:"category_id"`
	Amount     float64   `json:"amount"`
	Year       int       `json:"year"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`

	// Joined field for display
	CategoryName string `json:"category_name,omitempty"`
}

// BudgetRequest represents the payload to create or update a budget
type BudgetRequest struct {
	CategoryID int     `json:"category_id"`
	Amount     float64 `json:"amount"`
	Year       int     `json:"year"`
}

// BudgetDashboardSummary represents the calculated stats for the dashboard
type BudgetDashboardSummary struct {
	TotalAnnualBudget float64 `json:"total_annual_budget"`
	HighestAllocation float64 `json:"highest_allocation"`
	RemainingBudget   float64 `json:"remaining_budget"`
	SavingsTarget     float64 `json:"savings_target"`
}
