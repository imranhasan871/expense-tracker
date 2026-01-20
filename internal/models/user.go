package models

import (
	"time"
)

type UserRole string

const (
	RoleAdmin      UserRole = "admin"
	RoleManagement UserRole = "management"
	RoleExecutive  UserRole = "executive"
)

type User struct {
	ID                     int        `json:"id"`
	Username               string     `json:"username"`
	UserDisplayID          string     `json:"user_display_id"`
	Email                  string     `json:"email"`
	PasswordHash           string     `json:"-"`
	Role                   UserRole   `json:"role"`
	IsActive               bool       `json:"is_active"`
	PasswordSetToken       *string    `json:"-"`
	PasswordSetTokenExpiry *time.Time `json:"-"`
	CreatedAt              time.Time  `json:"created_at"`
	UpdatedAt              time.Time  `json:"updated_at"`
}

func (u *User) IsAdmin() bool {
	return u.Role == RoleAdmin
}

func (u *User) CanManage() bool {
	return u.Role == RoleAdmin || u.Role == RoleManagement
}

func (u *User) IsExecutive() bool {
	return u.Role == RoleExecutive
}

func (u *User) CanEnterExpenses() bool {
	return u.Role == RoleExecutive
}

func (u *User) CanViewAllExpenses() bool {
	return u.Role == RoleManagement || u.Role == RoleAdmin
}
