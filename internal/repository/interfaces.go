package repository

import "expense-tracker/internal/models"

type BudgetRepository interface {
	GetAll(year int) ([]models.Budget, error)
	CreateOrUpdate(categoryID int, amount float64, year int) (*models.Budget, error)
	GetDashboardSummary(year int) (*models.BudgetDashboardSummary, error)
	GetByCategory(categoryID, year int) (*models.Budget, error)
	GetMonitoringData(year int) ([]models.BudgetMonitoringItem, error)
	ToggleLock(budgetID int, isLocked bool) error
	IsLocked(categoryID, year int) (bool, error)
}

type ExpenseRepository interface {
	Create(req models.ExpenseRequest) (*models.Expense, error)
	GetAll(filter models.ExpenseFilter) ([]models.Expense, error)
	Delete(id int) error
	GetInsights(filter models.ExpenseFilter) (*models.ExpenseInsights, error)
	GetYearlyTotal(categoryID, year int) (float64, error)
}
