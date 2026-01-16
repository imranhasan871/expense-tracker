package repository

import (
	"database/sql"
	"errors"
	"expense-tracker/internal/models"
)

// BudgetRepository handles database operations for budgets
type BudgetRepository struct {
	db *sql.DB
}

// NewBudgetRepository creates a new budget repository
func NewBudgetRepository(db *sql.DB) *BudgetRepository {
	return &BudgetRepository{db: db}
}

// GetAll retrieves all budgets for a specific year
func (r *BudgetRepository) GetAll(year int) ([]models.Budget, error) {
	query := `SELECT b.id, b.category_id, b.amount, b.year, b.created_at, b.updated_at, c.name as category_name 
	          FROM budgets b 
	          JOIN categories c ON b.category_id = c.id 
	          WHERE b.year = $1 
	          ORDER BY c.name ASC`

	rows, err := r.db.Query(query, year)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	budgets := []models.Budget{}
	for rows.Next() {
		var b models.Budget
		err := rows.Scan(&b.ID, &b.CategoryID, &b.Amount, &b.Year, &b.CreatedAt, &b.UpdatedAt, &b.CategoryName)
		if err != nil {
			return nil, err
		}
		budgets = append(budgets, b)
	}

	return budgets, nil
}

// CreateOrUpdate sets a budget for a category and year
func (r *BudgetRepository) CreateOrUpdate(categoryID int, amount float64, year int) (*models.Budget, error) {
	var b models.Budget
	const minBudget = 10000

	if amount <= minBudget {
		return nil, errors.New("budget amount must be greater than 10000")
	}
	query := `INSERT INTO budgets (category_id, amount, year, updated_at) 
	          VALUES ($1, $2, $3, CURRENT_TIMESTAMP) 
	          ON CONFLICT (category_id, year) 
	          DO UPDATE SET amount = EXCLUDED.amount, updated_at = CURRENT_TIMESTAMP 
	          RETURNING id, category_id, amount, year, created_at, updated_at`

	err := r.db.QueryRow(query, categoryID, amount, year).Scan(
		&b.ID, &b.CategoryID, &b.Amount, &b.Year, &b.CreatedAt, &b.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &b, nil
}

// GetDashboardSummary calculates stats for the dashboard
func (r *BudgetRepository) GetDashboardSummary(year int) (*models.BudgetDashboardSummary, error) {
	summary := &models.BudgetDashboardSummary{}

	// Total Budget
	err := r.db.QueryRow("SELECT COALESCE(SUM(amount), 0) FROM budgets WHERE year = $1", year).Scan(&summary.TotalAnnualBudget)
	if err != nil {
		return nil, err
	}

	// Highest Allocation
	err = r.db.QueryRow("SELECT COALESCE(MAX(amount), 0) FROM budgets WHERE year = $1", year).Scan(&summary.HighestAllocation)
	if err != nil {
		return nil, err
	}

	// For simple design, let's assume savings target is 20% of total
	summary.SavingsTarget = summary.TotalAnnualBudget * 0.2

	// Remaining budget (simple calculation for now, we'll refine with actual spend later)
	summary.RemainingBudget = summary.TotalAnnualBudget * 0.8

	return summary, nil
}
