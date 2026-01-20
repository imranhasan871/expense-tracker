package repository

import (
	"database/sql"
	"errors"
	"expense-tracker/internal/models"
)

type sqlBudgetRepository struct {
	db *sql.DB
}

func NewBudgetRepository(db *sql.DB) BudgetRepository {
	return &sqlBudgetRepository{db: db}
}

func (r *sqlBudgetRepository) GetAll(year int) ([]models.Budget, error) {
	query := `SELECT b.id, b.category_id, b.amount, b.year, b.created_at, b.updated_at, c.name as category_name, b.is_locked
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
		err := rows.Scan(&b.ID, &b.CategoryID, &b.Amount, &b.Year, &b.CreatedAt, &b.UpdatedAt, &b.CategoryName, &b.IsLocked)
		if err != nil {
			return nil, err
		}
		budgets = append(budgets, b)
	}

	return budgets, nil
}

func (r *sqlBudgetRepository) CreateOrUpdate(categoryID int, amount float64, year int) (*models.Budget, error) {
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

func (r *sqlBudgetRepository) GetDashboardSummary(year int) (*models.BudgetDashboardSummary, error) {
	summary := &models.BudgetDashboardSummary{}

	err := r.db.QueryRow("SELECT COALESCE(SUM(amount), 0) FROM budgets WHERE year = $1", year).Scan(&summary.TotalAnnualBudget)
	if err != nil {
		return nil, err
	}

	err = r.db.QueryRow("SELECT COALESCE(MAX(amount), 0) FROM budgets WHERE year = $1", year).Scan(&summary.HighestAllocation)
	if err != nil {
		return nil, err
	}

	summary.SavingsTarget = summary.TotalAnnualBudget * 0.2

	summary.RemainingBudget = summary.TotalAnnualBudget * 0.8

	return summary, nil
}

func (r *sqlBudgetRepository) GetByCategory(categoryID, year int) (*models.Budget, error) {
	query := `SELECT id, category_id, amount, year, is_locked FROM budgets WHERE category_id = $1 AND year = $2`
	var b models.Budget
	err := r.db.QueryRow(query, categoryID, year).Scan(&b.ID, &b.CategoryID, &b.Amount, &b.Year, &b.IsLocked)
	if err != nil {
		return nil, err
	}
	return &b, nil
}

func (r *sqlBudgetRepository) GetMonitoringData(year int) ([]models.BudgetMonitoringItem, error) {
	query := `
		SELECT 
			b.id, 
			b.category_id, 
			c.name, 
			b.amount as budget_amount, 
			b.is_locked,
			COALESCE(SUM(e.amount), 0) as spent_amount
		FROM budgets b
		JOIN categories c ON b.category_id = c.id
		LEFT JOIN expenses e ON b.category_id = e.category_id AND EXTRACT(YEAR FROM e.expense_date) = b.year
		WHERE b.year = $1
		GROUP BY b.id, b.category_id, c.name, b.amount, b.is_locked
		ORDER BY c.name ASC
	`

	rows, err := r.db.Query(query, year)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := []models.BudgetMonitoringItem{}
	for rows.Next() {
		var item models.BudgetMonitoringItem
		err := rows.Scan(&item.BudgetID, &item.CategoryID, &item.CategoryName, &item.BudgetAmount, &item.IsLocked, &item.SpentAmount)
		if err != nil {
			return nil, err
		}

		if item.BudgetAmount > 0 {
			item.Percentage = (item.SpentAmount / item.BudgetAmount) * 100
		}
		items = append(items, item)
	}

	return items, nil
}

func (r *sqlBudgetRepository) ToggleLock(budgetID int, isLocked bool) error {
	_, err := r.db.Exec("UPDATE budgets SET is_locked = $1 WHERE id = $2", isLocked, budgetID)
	return err
}

func (r *sqlBudgetRepository) IsLocked(categoryID, year int) (bool, error) {
	var isLocked bool
	err := r.db.QueryRow("SELECT is_locked FROM budgets WHERE category_id = $1 AND year = $2", categoryID, year).Scan(&isLocked)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return isLocked, nil
}
