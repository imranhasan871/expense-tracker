package repository

import (
	"database/sql"
	"errors"
	"strings"

	"expense-tracker/internal/models"
)

// CategoryRepository handles database operations for categories
type CategoryRepository struct {
	db *sql.DB
}

// NewCategoryRepository creates a new category repository
func NewCategoryRepository(db *sql.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

// GetAll retrieves all categories
func (r *CategoryRepository) GetAll(activeOnly bool) ([]models.Category, error) {
	var query string
	if activeOnly {
		query = "SELECT id, name, is_active, created_at, updated_at FROM categories WHERE is_active = true ORDER BY name ASC"
	} else {
		query = "SELECT id, name, is_active, created_at, updated_at FROM categories ORDER BY name ASC"
	}

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	categories := []models.Category{}
	for rows.Next() {
		var category models.Category
		if err := rows.Scan(&category.ID, &category.Name, &category.IsActive, &category.CreatedAt, &category.UpdatedAt); err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}

	return categories, rows.Err()
}

// GetByID retrieves a category by ID
func (r *CategoryRepository) GetByID(id int) (*models.Category, error) {
	var category models.Category
	query := "SELECT id, name, is_active, created_at, updated_at FROM categories WHERE id = $1"

	err := r.db.QueryRow(query, id).Scan(
		&category.ID,
		&category.Name,
		&category.IsActive,
		&category.CreatedAt,
		&category.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("category not found")
	} else if err != nil {
		return nil, err
	}

	return &category, nil
}

// Create creates a new category
func (r *CategoryRepository) Create(name string, isActive bool) (*models.Category, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, errors.New("category name is required")
	}

	// Check if category already exists (case-insensitive)
	exists, err := r.ExistsByName(name)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("category with this name already exists")
	}

	var category models.Category
	query := `INSERT INTO categories (name, is_active) VALUES ($1, $2) 
	          RETURNING id, name, is_active, created_at, updated_at`

	err = r.db.QueryRow(query, name, isActive).Scan(
		&category.ID,
		&category.Name,
		&category.IsActive,
		&category.CreatedAt,
		&category.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &category, nil
}

// ExistsByName checks if a category with the given name exists (case-insensitive)
func (r *CategoryRepository) ExistsByName(name string) (bool, error) {
	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM categories WHERE LOWER(name) = LOWER($1))"
	err := r.db.QueryRow(query, name).Scan(&exists)
	return exists, err
}

// Update modifies an existing category
func (r *CategoryRepository) Update(id int, name string, isActive bool) (*models.Category, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, errors.New("category name is required")
	}

	// Check if name is already used by another category
	var duplicateExists bool
	err := r.db.QueryRow("SELECT EXISTS(SELECT 1 FROM categories WHERE LOWER(name) = LOWER($1) AND id != $2)", name, id).Scan(&duplicateExists)
	if err != nil {
		return nil, err
	}
	if duplicateExists {
		return nil, errors.New("another category with this name already exists")
	}

	var category models.Category
	query := `UPDATE categories SET name = $1, is_active = $2, updated_at = CURRENT_TIMESTAMP 
	          WHERE id = $3 RETURNING id, name, is_active, created_at, updated_at`

	err = r.db.QueryRow(query, name, isActive, id).Scan(
		&category.ID,
		&category.Name,
		&category.IsActive,
		&category.CreatedAt,
		&category.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("category not found")
	} else if err != nil {
		return nil, err
	}

	return &category, nil
}

// ToggleStatus switches the active state of a category
func (r *CategoryRepository) ToggleStatus(id int) (*models.Category, error) {
	var category models.Category
	query := `UPDATE categories SET is_active = NOT is_active, updated_at = CURRENT_TIMESTAMP 
	          WHERE id = $1 RETURNING id, name, is_active, created_at, updated_at`

	err := r.db.QueryRow(query, id).Scan(
		&category.ID,
		&category.Name,
		&category.IsActive,
		&category.CreatedAt,
		&category.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("category not found")
	} else if err != nil {
		return nil, err
	}

	return &category, nil
}

// InitializeDefaults creates default categories if they don't exist
func (r *CategoryRepository) InitializeDefaults() error {
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
		exists, err := r.ExistsByName(name)
		if err != nil {
			return err
		}

		if !exists {
			_, err := r.Create(name, true)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
