package repository

import (
	"database/sql"
	"expense-tracker/internal/models"
	"fmt"
	"time"
)

// BudgetEntryRepository handles database operations for budget entries
type BudgetEntryRepository struct {
	db *sql.DB
}

// NewBudgetEntryRepository creates a new budget entry repository
func NewBudgetEntryRepository(db *sql.DB) *BudgetEntryRepository {
	return &BudgetEntryRepository{db: db}
}

// Create adds a new budget entry, enforcing the date to be Jan 1st of the budget's year
func (r *BudgetEntryRepository) Create(budgetID int, amount float64, description string) (*models.BudgetEntry, error) {
	// 1. Get the year of the budget to set the entry date
	var year int
	err := r.db.QueryRow("SELECT year FROM budgets WHERE id = $1", budgetID).Scan(&year)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("budget with id %d not found", budgetID)
		}
		return nil, err
	}

	// 2. Set date to Jan 1st of that year
	entryDate := time.Date(year, time.January, 1, 0, 0, 0, 0, time.Local)

	// 3. Insert the entry
	query := `INSERT INTO budget_entries (budget_id, amount, description, date, updated_at)
	          VALUES ($1, $2, $3, $4, CURRENT_TIMESTAMP)
	          RETURNING id, created_at, updated_at`

	var entry models.BudgetEntry
	entry.BudgetID = budgetID
	entry.Amount = amount
	entry.Description = description
	entry.Date = entryDate

	err = r.db.QueryRow(query, budgetID, amount, description, entryDate).Scan(
		&entry.ID, &entry.CreatedAt, &entry.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &entry, nil
}

// GetByBudgetID retrieves all entries for a specific budget
func (r *BudgetEntryRepository) GetByBudgetID(budgetID int) ([]models.BudgetEntry, error) {
	query := `SELECT id, budget_id, amount, description, date, created_at, updated_at 
	          FROM budget_entries 
	          WHERE budget_id = $1 
	          ORDER BY created_at DESC`

	rows, err := r.db.Query(query, budgetID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []models.BudgetEntry
	for rows.Next() {
		var e models.BudgetEntry
		err := rows.Scan(
			&e.ID, &e.BudgetID, &e.Amount, &e.Description, &e.Date, &e.CreatedAt, &e.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		entries = append(entries, e)
	}

	return entries, nil
}
