package service

import (
	"context"
	"fmt"

	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/goProjects/loan_app/app/model"
	"github.com/goProjects/loan_app/lib/logger"
	"github.com/goProjects/loan_app/lib/utils"
)

func (svc *userService) CreateUser(ctx context.Context, db *gorm.DB, name, password, roleID string) (*model.User, error) {
	var user *model.User
	var err error

	err = db.Transaction(func(tx *gorm.DB) error {
		user = &model.User{
			Name:           name,
			HashedPassword: utils.Hash(password),
		}
		user, err = svc.userStore.CreateUser(ctx, tx, user)
		if err != nil {
			logger.E(ctx, "create user failed", zap.Error(err))
			return fmt.Errorf("creating user in service: %w", err)
		}

		if _, err := svc.userStore.CreateUserRole(ctx, tx, user.ID, roleID); err != nil {
			logger.E(ctx, "create user role failed", zap.Error(err))
			return fmt.Errorf("creating user role in service: %w", err)
		}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("error in creating user: %w", err)
	}

	return user, nil
}
