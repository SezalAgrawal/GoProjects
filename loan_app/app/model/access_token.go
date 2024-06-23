package model

import "time"

type AccessToken struct {
	ID        string
	UserID    string
	Token     string
	Deleted   bool
	DeletedAt time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}
