package models

import (
	"fmt"
	"time"
)

type Expense struct {
	ID          int       `json:"id"`
	CategoryID  int       `json:"category_id"`
	UserID      *int      `json:"user_id,omitempty"`
	Amount      float64   `json:"amount"`
	ExpenseDate time.Time `json:"expense_date"`
	Remarks     string    `json:"remarks"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	CategoryName string `json:"category_name,omitempty"`
	UserName     string `json:"user_name,omitempty"`
}

type ExpenseRequest struct {
	CategoryID  int     `json:"category_id"`
	UserID      int     `json:"user_id"`
	Amount      float64 `json:"amount"`
	ExpenseDate string  `json:"expense_date"`
	Remarks     string  `json:"remarks"`
}

type ExpenseFilter struct {
	StartDate  string  `json:"start_date"`
	EndDate    string  `json:"end_date"`
	CategoryID int     `json:"category_id"`
	UserID     int     `json:"user_id"`
	SearchText string  `json:"search_text"`
	MinAmount  float64 `json:"min_amount"`
	MaxAmount  float64 `json:"max_amount"`
}

func (f *ExpenseFilter) Validate() error {
	if f.StartDate != "" && f.EndDate != "" {
		start, err1 := time.Parse("2006-01-02", f.StartDate)
		end, err2 := time.Parse("2006-01-02", f.EndDate)
		if err1 == nil && err2 == nil && start.After(end) {
			return fmt.Errorf("start date must be before or equal to end date")
		}
	}
	if f.MinAmount > 0 && f.MaxAmount > 0 && f.MinAmount > f.MaxAmount {
		return fmt.Errorf("minimum amount must be less than or equal to maximum amount")
	}
	return nil
}

type CategorySpending struct {
	CategoryID   int     `json:"category_id"`
	CategoryName string  `json:"category_name"`
	TotalAmount  float64 `json:"total_amount"`
	Count        int     `json:"count"`
}

type DaySpending struct {
	DayOfWeek   int     `json:"day_of_week"`
	DayName     string  `json:"day_name"`
	TotalAmount float64 `json:"total_amount"`
	Count       int     `json:"count"`
}

type ExpenseInsights struct {
	TotalSpent       float64 `json:"total_spent"`
	TransactionCount int     `json:"transaction_count"`
	AverageExpense   float64 `json:"average_expense"`

	PreviousPeriodTotal float64 `json:"previous_period_total"`
	SpendingChange      float64 `json:"spending_change"`

	TopCategories []CategorySpending `json:"top_categories"`

	SpendingByDay []DaySpending `json:"spending_by_day"`
}
