package service

import (
	"context"
	"fmt"

	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/goProjects/loan_app/app/model"
	"github.com/goProjects/loan_app/lib/logger"
)

func (svc *loanService) GetAllLoans(ctx context.Context, db *gorm.DB, userID string) ([]*model.Loan, error) {
	loans, err := svc.loanStore.GetAllLoans(ctx, db, userID)
	if err != nil {
		logger.E(ctx, "getting loan failed", zap.Error(err))
		return nil, fmt.Errorf("getting loan in service: %w", err)
	}

	return loans, nil
}
