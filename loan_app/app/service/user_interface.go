package service

import (
	"context"

	"github.com/goProjects/loan_app/app/model"
	"github.com/goProjects/loan_app/app/store"
	"gorm.io/gorm"
)

type UserInterface interface {
	CreateUser(ctx context.Context, db *gorm.DB, name, password, roleID string) (*model.User, error)
	GenerateAccessToken(ctx context.Context, db *gorm.DB, name, password string) (*model.User, string, error)
	DeleteAllAccessTokens(ctx context.Context, db *gorm.DB, userID string) error
}

type userService struct {
	userStore store.UserStore
}

func NewUserService() UserInterface {
	svc := &userService{
		userStore: store.NewUserStore(),
	}
	return svc
}
