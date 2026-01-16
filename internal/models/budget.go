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
	IsLocked     bool   `json:"is_locked"`
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

// BudgetMonitoringItem represents a row in the monitoring system
type BudgetMonitoringItem struct {
	BudgetID     int     `json:"budget_id"`
	CategoryID   int     `json:"category_id"`
	CategoryName string  `json:"category_name"`
	BudgetAmount float64 `json:"budget_amount"`
	SpentAmount  float64 `json:"spent_amount"`
	Percentage   float64 `json:"percentage"`
	IsLocked     bool    `json:"is_locked"`
}
