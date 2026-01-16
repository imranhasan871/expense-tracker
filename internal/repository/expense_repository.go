package repository

import (
	"database/sql"
	"expense-tracker/internal/models"
	"fmt"
	"strings"
	"time"
)

type PostgresExpenseRepository struct {
	db *sql.DB
}

func NewExpenseRepository(db *sql.DB) *PostgresExpenseRepository {
	return &PostgresExpenseRepository{db: db}
}

func (r *PostgresExpenseRepository) Create(req models.ExpenseRequest) (*models.Expense, error) {
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

func (r *PostgresExpenseRepository) GetAll(filter models.ExpenseFilter) ([]models.Expense, error) {
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

	if filter.SearchText != "" {
		conditions = append(conditions, fmt.Sprintf("e.remarks ILIKE $%d", argCount))
		args = append(args, "%"+filter.SearchText+"%")
		argCount++
	}

	if filter.MinAmount > 0 {
		conditions = append(conditions, fmt.Sprintf("e.amount >= $%d", argCount))
		args = append(args, filter.MinAmount)
		argCount++
	}

	if filter.MaxAmount > 0 {
		conditions = append(conditions, fmt.Sprintf("e.amount <= $%d", argCount))
		args = append(args, filter.MaxAmount)
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

func (r *PostgresExpenseRepository) Delete(id int) error {
	_, err := r.db.Exec("DELETE FROM expenses WHERE id = $1", id)
	return err
}

func (r *PostgresExpenseRepository) GetInsights(filter models.ExpenseFilter) (*models.ExpenseInsights, error) {
	insights := &models.ExpenseInsights{}

	currentStats, err := r.getPeriodStats(filter)
	if err != nil {
		return nil, err
	}
	insights.TotalSpent = currentStats.total
	insights.TransactionCount = currentStats.count
	if currentStats.count > 0 {
		insights.AverageExpense = currentStats.total / float64(currentStats.count)
	}

	if filter.StartDate != "" && filter.EndDate != "" {
		start, _ := time.Parse("2006-01-02", filter.StartDate)
		end, _ := time.Parse("2006-01-02", filter.EndDate)
		duration := end.Sub(start)

		prevFilter := models.ExpenseFilter{
			StartDate:  start.Add(-duration - 24*time.Hour).Format("2006-01-02"),
			EndDate:    start.Add(-24 * time.Hour).Format("2006-01-02"),
			CategoryID: filter.CategoryID,
		}
		prevStats, err := r.getPeriodStats(prevFilter)
		if err == nil {
			insights.PreviousPeriodTotal = prevStats.total
			if prevStats.total > 0 {
				insights.SpendingChange = ((currentStats.total - prevStats.total) / prevStats.total) * 100
			} else if currentStats.total > 0 {
				insights.SpendingChange = 100
			}
		}
	}

	topCategories, err := r.getTopCategories(filter)
	if err != nil {
		return nil, err
	}
	insights.TopCategories = topCategories

	spendingByDay, err := r.getSpendingByDay(filter)
	if err != nil {
		return nil, err
	}
	insights.SpendingByDay = spendingByDay

	return insights, nil
}

type periodStats struct {
	total float64
	count int
}

func (r *PostgresExpenseRepository) getPeriodStats(filter models.ExpenseFilter) (*periodStats, error) {
	query := `SELECT COALESCE(SUM(amount), 0), COUNT(*) FROM expenses WHERE 1=1`
	var args []interface{}
	argCount := 1

	if filter.StartDate != "" {
		query += fmt.Sprintf(" AND expense_date >= $%d", argCount)
		args = append(args, filter.StartDate)
		argCount++
	}
	if filter.EndDate != "" {
		query += fmt.Sprintf(" AND expense_date <= $%d", argCount)
		args = append(args, filter.EndDate)
		argCount++
	}
	if filter.CategoryID > 0 {
		query += fmt.Sprintf(" AND category_id = $%d", argCount)
		args = append(args, filter.CategoryID)
		argCount++
	}

	var stats periodStats
	err := r.db.QueryRow(query, args...).Scan(&stats.total, &stats.count)
	return &stats, err
}

func (r *PostgresExpenseRepository) getTopCategories(filter models.ExpenseFilter) ([]models.CategorySpending, error) {
	query := `SELECT e.category_id, c.name, COALESCE(SUM(e.amount), 0) as total, COUNT(*) as cnt
	          FROM expenses e
	          JOIN categories c ON e.category_id = c.id
	          WHERE 1=1`
	var args []interface{}
	argCount := 1

	if filter.StartDate != "" {
		query += fmt.Sprintf(" AND e.expense_date >= $%d", argCount)
		args = append(args, filter.StartDate)
		argCount++
	}
	if filter.EndDate != "" {
		query += fmt.Sprintf(" AND e.expense_date <= $%d", argCount)
		args = append(args, filter.EndDate)
		argCount++
	}
	if filter.CategoryID > 0 {
		query += fmt.Sprintf(" AND e.category_id = $%d", argCount)
		args = append(args, filter.CategoryID)
		argCount++
	}

	query += " GROUP BY e.category_id, c.name ORDER BY total DESC LIMIT 5"

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []models.CategorySpending
	for rows.Next() {
		var cs models.CategorySpending
		if err := rows.Scan(&cs.CategoryID, &cs.CategoryName, &cs.TotalAmount, &cs.Count); err != nil {
			return nil, err
		}
		results = append(results, cs)
	}
	return results, nil
}

func (r *PostgresExpenseRepository) getSpendingByDay(filter models.ExpenseFilter) ([]models.DaySpending, error) {
	dayNames := []string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"}

	query := `SELECT EXTRACT(DOW FROM expense_date)::int as dow, 
	                 COALESCE(SUM(amount), 0) as total, COUNT(*) as cnt
	          FROM expenses WHERE 1=1`
	var args []interface{}
	argCount := 1

	if filter.StartDate != "" {
		query += fmt.Sprintf(" AND expense_date >= $%d", argCount)
		args = append(args, filter.StartDate)
		argCount++
	}
	if filter.EndDate != "" {
		query += fmt.Sprintf(" AND expense_date <= $%d", argCount)
		args = append(args, filter.EndDate)
		argCount++
	}
	if filter.CategoryID > 0 {
		query += fmt.Sprintf(" AND category_id = $%d", argCount)
		args = append(args, filter.CategoryID)
		argCount++
	}

	query += " GROUP BY dow ORDER BY cnt DESC"

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []models.DaySpending
	for rows.Next() {
		var ds models.DaySpending
		if err := rows.Scan(&ds.DayOfWeek, &ds.TotalAmount, &ds.Count); err != nil {
			return nil, err
		}
		ds.DayName = dayNames[ds.DayOfWeek]
		results = append(results, ds)
	}
	return results, nil
}

func (r *PostgresExpenseRepository) GetYearlyTotal(categoryID, year int) (float64, error) {
	query := `SELECT COALESCE(SUM(amount), 0) FROM expenses WHERE category_id = $1 AND EXTRACT(YEAR FROM expense_date) = $2`
	var total float64
	err := r.db.QueryRow(query, categoryID, year).Scan(&total)
	return total, err
}
