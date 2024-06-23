package model

import "time"

type User struct {
	ID             string
	Name           string
	HashedPassword string
	CreatedAt      time.Time
	UpdatedAt      time.Time

	UserRoles []*UserRole
}
