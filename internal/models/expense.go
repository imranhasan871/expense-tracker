package models

import (
	"fmt"
	"time"
)

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
	StartDate  string  `json:"start_date"`
	EndDate    string  `json:"end_date"`
	CategoryID int     `json:"category_id"`
	SearchText string  `json:"search_text"` // Search in remarks
	MinAmount  float64 `json:"min_amount"`
	MaxAmount  float64 `json:"max_amount"`
}

// Validate checks filter logic (e.g., start_date <= end_date, min <= max)
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

// CategorySpending represents spending per category for insights
type CategorySpending struct {
	CategoryID   int     `json:"category_id"`
	CategoryName string  `json:"category_name"`
	TotalAmount  float64 `json:"total_amount"`
	Count        int     `json:"count"`
}

// DaySpending represents spending per day-of-week
type DaySpending struct {
	DayOfWeek   int     `json:"day_of_week"` // 0=Sunday, 6=Saturday
	DayName     string  `json:"day_name"`
	TotalAmount float64 `json:"total_amount"`
	Count       int     `json:"count"`
}

// ExpenseInsights contains aggregated spending analysis
type ExpenseInsights struct {
	// Summary stats for current period
	TotalSpent       float64 `json:"total_spent"`
	TransactionCount int     `json:"transaction_count"`
	AverageExpense   float64 `json:"average_expense"`

	// Comparison with previous period
	PreviousPeriodTotal float64 `json:"previous_period_total"`
	SpendingChange      float64 `json:"spending_change"` // Percentage change

	// Top spending categories (sorted by total)
	TopCategories []CategorySpending `json:"top_categories"`

	// Most frequent spending days
	SpendingByDay []DaySpending `json:"spending_by_day"`
}
