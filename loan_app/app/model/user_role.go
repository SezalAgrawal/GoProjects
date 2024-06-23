package model

import "time"

type UserRole struct {
	ID        string
	UserID    string
	RoleID    string
	CreatedAt time.Time
	UpdatedAt time.Time

	Role *Role
}
