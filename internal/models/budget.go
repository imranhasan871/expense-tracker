package models

import "time"

type Budget struct {
	ID         int       `json:"id"`
	CategoryID int       `json:"category_id"`
	Amount     float64   `json:"amount"`
	Year       int       `json:"year"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`

	CategoryName string `json:"category_name,omitempty"`
	IsLocked     bool   `json:"is_locked"`
}

type BudgetRequest struct {
	CategoryID int     `json:"category_id"`
	Amount     float64 `json:"amount"`
	Year       int     `json:"year"`
}

type BudgetDashboardSummary struct {
	TotalAnnualBudget float64 `json:"total_annual_budget"`
	HighestAllocation float64 `json:"highest_allocation"`
	RemainingBudget   float64 `json:"remaining_budget"`
	SavingsTarget     float64 `json:"savings_target"`
}

type BudgetMonitoringItem struct {
	BudgetID     int     `json:"budget_id"`
	CategoryID   int     `json:"category_id"`
	CategoryName string  `json:"category_name"`
	BudgetAmount float64 `json:"budget_amount"`
	SpentAmount  float64 `json:"spent_amount"`
	Percentage   float64 `json:"percentage"`
	IsLocked     bool    `json:"is_locked"`
}
