package repository

import (
	"database/sql"
	"expense-tracker/internal/models"
	"fmt"
	"time"
)

type BudgetEntryRepository struct {
	db *sql.DB
}

func NewBudgetEntryRepository(db *sql.DB) *BudgetEntryRepository {
	return &BudgetEntryRepository{db: db}
}

func (r *BudgetEntryRepository) Create(budgetID int, amount float64, description string) (*models.BudgetEntry, error) {
	var year int
	err := r.db.QueryRow("SELECT year FROM budgets WHERE id = $1", budgetID).Scan(&year)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("budget with id %d not found", budgetID)
		}
		return nil, err
	}

	entryDate := time.Date(year, time.January, 1, 0, 0, 0, 0, time.Local)

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
