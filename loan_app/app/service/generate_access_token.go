package service

import (
	"context"
	"fmt"

	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/goProjects/loan_app/app/model"
	"github.com/goProjects/loan_app/app/store"
	"github.com/goProjects/loan_app/lib/logger"
	"github.com/goProjects/loan_app/lib/utils"
)

func (svc *userService) GenerateAccessToken(ctx context.Context, db *gorm.DB, name, password string) (*model.User, string, error) {
	hashedPassword := utils.Hash(password)
	user, err := svc.userStore.GetUserByUsernamePassword(ctx, db, name, hashedPassword)
	if err != nil {
		if err != store.ErrNotFound {
			logger.E(ctx, "error fetching user", zap.Error(err), zap.String("name", name))
		}
		return nil, "", fmt.Errorf("unable to fetch user: %w", err)
	}

	accessToken, err := svc.userStore.GenerateAccessToken(ctx, db, user.ID)
	if err != nil {
		logger.E(ctx, "create access token failed", zap.Error(err))
		return nil, "", fmt.Errorf("creating access token in service: %w", err)
	}

	return user, accessToken.Token, nil
}
