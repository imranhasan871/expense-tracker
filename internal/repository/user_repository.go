package repository

import (
	"database/sql"
	"errors"
	"expense-tracker/internal/models"
	"time"
)

type sqlUserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &sqlUserRepository{db: db}
}

func (r *sqlUserRepository) Create(user *models.User) error {
	query := `INSERT INTO users (username, user_display_id, email, password_hash, role, is_active, password_set_token, password_set_token_expiry) 
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id, created_at, updated_at`

	return r.db.QueryRow(
		query,
		user.Username,
		user.UserDisplayID,
		user.Email,
		user.PasswordHash,
		user.Role,
		user.IsActive,
		user.PasswordSetToken,
		user.PasswordSetTokenExpiry,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
}

func (r *sqlUserRepository) GetByID(id int) (*models.User, error) {
	var user models.User
	query := `SELECT id, username, user_display_id, email, password_hash, role, is_active, password_set_token, password_set_token_expiry, created_at, updated_at 
	          FROM users WHERE id = $1`

	err := r.db.QueryRow(query, id).Scan(
		&user.ID, &user.Username, &user.UserDisplayID, &user.Email, &user.PasswordHash,
		&user.Role, &user.IsActive, &user.PasswordSetToken, &user.PasswordSetTokenExpiry,
		&user.CreatedAt, &user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("user not found")
	}
	return &user, err
}

func (r *sqlUserRepository) GetByEmail(email string) (*models.User, error) {
	var user models.User
	query := `SELECT id, username, user_display_id, email, password_hash, role, is_active, password_set_token, password_set_token_expiry, created_at, updated_at 
	          FROM users WHERE email = $1`

	err := r.db.QueryRow(query, email).Scan(
		&user.ID, &user.Username, &user.UserDisplayID, &user.Email, &user.PasswordHash,
		&user.Role, &user.IsActive, &user.PasswordSetToken, &user.PasswordSetTokenExpiry,
		&user.CreatedAt, &user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("user not found")
	}
	return &user, err
}

func (r *sqlUserRepository) GetByDisplayID(displayID string) (*models.User, error) {
	var user models.User
	query := `SELECT id, username, user_display_id, email, password_hash, role, is_active, password_set_token, password_set_token_expiry, created_at, updated_at 
	          FROM users WHERE user_display_id = $1`

	err := r.db.QueryRow(query, displayID).Scan(
		&user.ID, &user.Username, &user.UserDisplayID, &user.Email, &user.PasswordHash,
		&user.Role, &user.IsActive, &user.PasswordSetToken, &user.PasswordSetTokenExpiry,
		&user.CreatedAt, &user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("user not found")
	}
	return &user, err
}

func (r *sqlUserRepository) UpdatePassword(id int, passwordHash string) error {
	query := `UPDATE users SET password_hash = $1, is_active = TRUE, password_set_token = NULL, password_set_token_expiry = NULL, updated_at = CURRENT_TIMESTAMP WHERE id = $2`
	_, err := r.db.Exec(query, passwordHash, id)
	return err
}

func (r *sqlUserRepository) SetPasswordToken(email, token string, expiry time.Time) error {
	query := `UPDATE users SET password_set_token = $1, password_set_token_expiry = $2, updated_at = CURRENT_TIMESTAMP WHERE email = $3`
	_, err := r.db.Exec(query, token, expiry, email)
	return err
}

func (r *sqlUserRepository) GetByToken(token string) (*models.User, error) {
	var user models.User
	query := `SELECT id, username, user_display_id, email, password_hash, role, is_active, password_set_token, password_set_token_expiry, created_at, updated_at 
	          FROM users WHERE password_set_token = $1 AND password_set_token_expiry > CURRENT_TIMESTAMP`

	err := r.db.QueryRow(query, token).Scan(
		&user.ID, &user.Username, &user.UserDisplayID, &user.Email, &user.PasswordHash,
		&user.Role, &user.IsActive, &user.PasswordSetToken, &user.PasswordSetTokenExpiry,
		&user.CreatedAt, &user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("invalid or expired token")
	}
	return &user, err
}

func (r *sqlUserRepository) GetAll() ([]models.User, error) {
	query := `SELECT id, username, user_display_id, email, role, is_active, created_at, updated_at FROM users ORDER BY created_at DESC`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []models.User{}
	for rows.Next() {
		var user models.User
		err := rows.Scan(
			&user.ID, &user.Username, &user.UserDisplayID, &user.Email,
			&user.Role, &user.IsActive, &user.CreatedAt, &user.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}
