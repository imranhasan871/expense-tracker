package repository

import (
	"database/sql"
	"expense-tracker/internal/models"
	"fmt"
	"strings"
	"time"
)

// ExpenseRepository handles database operations for expenses
type ExpenseRepository struct {
	db *sql.DB
}

// NewExpenseRepository creates a new expense repository
func NewExpenseRepository(db *sql.DB) *ExpenseRepository {
	return &ExpenseRepository{db: db}
}

// Create records a new expense
func (r *ExpenseRepository) Create(req models.ExpenseRequest) (*models.Expense, error) {
	expenseDate, err := time.Parse("2006-01-02", req.ExpenseDate)
	if err != nil {
		return nil, fmt.Errorf("invalid date format: %v", err)
	}

	var e models.Expense
	query := `INSERT INTO expenses (category_id, amount, expense_date, remarks) 
	          VALUES ($1, $2, $3, $4) 
	          RETURNING id, category_id, amount, expense_date, remarks, created_at, updated_at`

	err = r.db.QueryRow(query, req.CategoryID, req.Amount, expenseDate, req.Remarks).Scan(
		&e.ID, &e.CategoryID, &e.Amount, &e.ExpenseDate, &e.Remarks, &e.CreatedAt, &e.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &e, nil
}

// GetAll retrieves expenses with optional filters
func (r *ExpenseRepository) GetAll(filter models.ExpenseFilter) ([]models.Expense, error) {
	query := `SELECT e.id, e.category_id, e.amount, e.expense_date, e.remarks, e.created_at, e.updated_at, c.name as category_name 
	          FROM expenses e 
	          JOIN categories c ON e.category_id = c.id`

	var conditions []string
	var args []interface{}
	argCount := 1

	if filter.StartDate != "" {
		conditions = append(conditions, fmt.Sprintf("e.expense_date >= $%d", argCount))
		args = append(args, filter.StartDate)
		argCount++
	}

	if filter.EndDate != "" {
		conditions = append(conditions, fmt.Sprintf("e.expense_date <= $%d", argCount))
		args = append(args, filter.EndDate)
		argCount++
	}

	if filter.CategoryID > 0 {
		conditions = append(conditions, fmt.Sprintf("e.category_id = $%d", argCount))
		args = append(args, filter.CategoryID)
		argCount++
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	query += " ORDER BY e.expense_date DESC, e.created_at DESC"

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	expenses := []models.Expense{}
	for rows.Next() {
		var e models.Expense
		err := rows.Scan(&e.ID, &e.CategoryID, &e.Amount, &e.ExpenseDate, &e.Remarks, &e.CreatedAt, &e.UpdatedAt, &e.CategoryName)
		if err != nil {
			return nil, err
		}
		expenses = append(expenses, e)
	}

	return expenses, nil
}

// Delete removes an expense record
func (r *ExpenseRepository) Delete(id int) error {
	_, err := r.db.Exec("DELETE FROM expenses WHERE id = $1", id)
	return err
}
