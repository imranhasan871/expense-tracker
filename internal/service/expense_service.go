package service

import (
	"errors"
	"strconv"

	"expense-tracker/internal/models"
)

type ExpenseRepositoryInterface interface {
	Create(req models.ExpenseRequest) (*models.Expense, error)
	GetAll(filter models.ExpenseFilter) ([]models.Expense, error)
	Delete(id int) error
	GetInsights(filter models.ExpenseFilter) (*models.ExpenseInsights, error)
}

type BudgetRepositoryInterface interface {
	IsLocked(categoryID, year int) (bool, error)
}

type ExpenseService struct {
	repo       ExpenseRepositoryInterface
	budgetRepo BudgetRepositoryInterface
}

func NewExpenseService(repo ExpenseRepositoryInterface, budgetRepo BudgetRepositoryInterface) *ExpenseService {
	return &ExpenseService{
		repo:       repo,
		budgetRepo: budgetRepo,
	}
}

func (s *ExpenseService) Create(req models.ExpenseRequest, user *models.User) (*models.Expense, error) {
	// Fix boolean logic: Use OR for allowed roles, then negate
	if !(user.Role == models.RoleExecutive || user.Role == models.RoleAdmin) {
		return nil, errors.New("only executives and admins can enter expenses")
	}

	req.UserID = user.ID // Ensure the expense is owned by the requester
	// Validate request
	if req.CategoryID <= 0 {
		return nil, errors.New("category ID is required")
	}
	if req.Amount <= 0 {
		return nil, errors.New("amount must be greater than 0")
	}
	if req.ExpenseDate == "" {
		return nil, errors.New("expense date is required")
	}

	// Check if budget is locked for this category and year
	if len(req.ExpenseDate) >= 4 {
		year, err := strconv.Atoi(req.ExpenseDate[:4])
		if err == nil {
			isLocked, err := s.budgetRepo.IsLocked(req.CategoryID, year)
			if err != nil {
				return nil, errors.New("failed to check budget lock status")
			}
			if isLocked {
				return nil, errors.New("spending is temporarily locked for this category")
			}
		}
	}

	return s.repo.Create(req)
}

func (s *ExpenseService) GetAll(filter models.ExpenseFilter, user *models.User) ([]models.Expense, error) {
	if user.Role == models.RoleExecutive {
		filter.UserID = user.ID
	}

	if err := filter.Validate(); err != nil {
		return nil, err
	}
	return s.repo.GetAll(filter)
}

func (s *ExpenseService) GetInsights(filter models.ExpenseFilter, user *models.User) (*models.ExpenseInsights, error) {
	if user.Role == models.RoleExecutive {
		filter.UserID = user.ID
	}

	if err := filter.Validate(); err != nil {
		return nil, err
	}
	return s.repo.GetInsights(filter)
}

func (s *ExpenseService) Delete(id int) error {
	if id <= 0 {
		return errors.New("expense ID must be greater than 0")
	}
	return s.repo.Delete(id)
}
