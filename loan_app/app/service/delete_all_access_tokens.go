package service

import (
	"context"
	"fmt"

	"gorm.io/gorm"
)

func (svc *userService) DeleteAllAccessTokens(ctx context.Context, db *gorm.DB, userID string) error {
	if err := svc.userStore.DeleteAllAccessTokens(ctx, db, userID); err != nil {
		return fmt.Errorf("unable to delete all access tokens: %w", err)
	}

	return nil
}
