package service

import (
	"errors"
	"strings"

	"expense-tracker/internal/models"
)

type CategoryRepository interface {
	GetAll(activeOnly bool) ([]models.Category, error)
	GetByID(id int) (*models.Category, error)
	Create(name string, isActive bool) (*models.Category, error)
	Update(id int, name string, isActive bool) (*models.Category, error)
	ToggleStatus(id int) (*models.Category, error)
	ExistsByName(name string) (bool, error)
}

type CategoryService struct {
	repo CategoryRepository
}

func NewCategoryService(repo CategoryRepository) *CategoryService {
	return &CategoryService{repo: repo}
}

func (s *CategoryService) GetAll(activeOnly bool) ([]models.Category, error) {
	return s.repo.GetAll(activeOnly)
}

func (s *CategoryService) GetByID(id int) (*models.Category, error) {
	return s.repo.GetByID(id)
}

func (s *CategoryService) Create(name string, isActive bool) (*models.Category, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, errors.New("category name is required")
	}
	return s.repo.Create(name, isActive)
}

func (s *CategoryService) Update(id int, name string, isActive bool) (*models.Category, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, errors.New("category name is required")
	}
	return s.repo.Update(id, name, isActive)
}

func (s *CategoryService) ToggleStatus(id int) (*models.Category, error) {
	return s.repo.ToggleStatus(id)
}

func (s *CategoryService) InitializeDefaults() error {
	defaultCategories := []string{
		"Food",
		"Transport",
		"Rent",
		"Utilities",
		"Marketing",
		"Salary",
		"Office Rent",
		"HR Development",
		"Entertainment",
	}

	for _, name := range defaultCategories {
		exists, err := s.repo.ExistsByName(name)
		if err != nil {
			return err
		}

		if !exists {
			_, err := s.repo.Create(name, true)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
