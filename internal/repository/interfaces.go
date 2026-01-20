package repository

import (
	"expense-tracker/internal/models"
	"time"
)

type CategoryRepository interface {
	GetAll(activeOnly bool) ([]models.Category, error)
	GetByID(id int) (*models.Category, error)
	Create(name string, isActive bool) (*models.Category, error)
	Update(id int, name string, isActive bool) (*models.Category, error)
	ToggleStatus(id int) (*models.Category, error)
	ExistsByName(name string) (bool, error)
}

type BudgetRepository interface {
	GetAll(year int) ([]models.Budget, error)
	CreateOrUpdate(categoryID int, amount float64, year int) (*models.Budget, error)
	GetDashboardSummary(year int) (*models.BudgetDashboardSummary, error)
	GetByCategory(categoryID, year int) (*models.Budget, error)
	GetMonitoringData(year int) ([]models.BudgetMonitoringItem, error)
	ToggleLock(budgetID int, isLocked bool) error
	IsLocked(categoryID, year int) (bool, error)
}

type UserRepository interface {
	Create(user *models.User) error
	GetByID(id int) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	GetByDisplayID(displayID string) (*models.User, error)
	UpdatePassword(id int, passwordHash string) error
	UpdateRole(id int, role models.UserRole) error
	SetPasswordToken(email, token string, expiry time.Time) error
	GetByToken(token string) (*models.User, error)
	GetAll() ([]models.User, error)
}

type ExpenseRepository interface {
	Create(req models.ExpenseRequest) (*models.Expense, error)
	GetAll(filter models.ExpenseFilter) ([]models.Expense, error)
	Delete(id int) error
	GetInsights(filter models.ExpenseFilter) (*models.ExpenseInsights, error)
	GetYearlyTotal(categoryID, year int) (float64, error)
}
