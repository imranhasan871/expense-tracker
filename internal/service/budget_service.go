package service

import (
	"database/sql"
	"errors"

	"expense-tracker/internal/models"
)

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
	GetYearlyTotal(categoryID, year int) (float64, error)
}

type BudgetService struct {
	repo        BudgetRepository
	expenseRepo ExpenseRepository
}

func NewBudgetService(repo BudgetRepository, expenseRepo ExpenseRepository) *BudgetService {
	return &BudgetService{
		repo:        repo,
		expenseRepo: expenseRepo,
	}
}

func (s *BudgetService) GetAll(year int) ([]models.Budget, error) {
	if year <= 0 {
		return nil, errors.New("year must be greater than 0")
	}
	return s.repo.GetAll(year)
}

func (s *BudgetService) GetDashboardSummary(year int) (*models.BudgetDashboardSummary, error) {
	if year <= 0 {
		return nil, errors.New("year must be greater than 0")
	}
	return s.repo.GetDashboardSummary(year)
}

func (s *BudgetService) GetStatus(categoryID, year int) (*models.BudgetStatus, error) {
	if categoryID <= 0 {
		return nil, errors.New("category ID must be greater than 0")
	}
	if year <= 0 {
		return nil, errors.New("year must be greater than 0")
	}

	budget, err := s.repo.GetByCategory(categoryID, year)
	if err != nil {
		if err == sql.ErrNoRows {
			return &models.BudgetStatus{
				Allocated: 0,
				Spent:     0,
				Remaining: 0,
				Percent:   0,
				IsLocked:  false,
			}, nil
		}
		return nil, err
	}

	spent, err := s.expenseRepo.GetYearlyTotal(categoryID, year)
	if err != nil {
		return nil, err
	}

	remaining := budget.Amount - spent
	var percent float64
	if budget.Amount > 0 {
		percent = (spent / budget.Amount) * 100
	}

	return &models.BudgetStatus{
		Allocated: budget.Amount,
		Spent:     spent,
		Remaining: remaining,
		Percent:   percent,
		IsLocked:  budget.IsLocked,
	}, nil
}

func (s *BudgetService) CreateOrUpdate(categoryID int, amount float64, year int) (*models.Budget, error) {
	if categoryID <= 0 {
		return nil, errors.New("category ID must be greater than 0")
	}
	if amount <= 0 {
		return nil, errors.New("amount must be greater than 0")
	}
	if year <= 0 {
		return nil, errors.New("year must be greater than 0")
	}

	return s.repo.CreateOrUpdate(categoryID, amount, year)
}

func (s *BudgetService) ToggleLock(budgetID int, isLocked bool) error {
	if budgetID <= 0 {
		return errors.New("budget ID must be greater than 0")
	}
	return s.repo.ToggleLock(budgetID, isLocked)
}

func (s *BudgetService) GetMonitoringData(year int) ([]models.BudgetMonitoringItem, error) {
	if year <= 0 {
		return nil, errors.New("year must be greater than 0")
	}
	return s.repo.GetMonitoringData(year)
}

func (s *BudgetService) IsLocked(categoryID, year int) (bool, error) {
	if categoryID <= 0 {
		return false, errors.New("category ID must be greater than 0")
	}
	if year <= 0 {
		return false, errors.New("year must be greater than 0")
	}
	return s.repo.IsLocked(categoryID, year)
}
