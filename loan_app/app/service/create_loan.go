package service

import (
	"context"
	"fmt"

	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/goProjects/loan_app/app/model"
	"github.com/goProjects/loan_app/lib/logger"
)

type CreateLoanParams struct {
	AmountInCents   int64
	Term            int64
	FrequencyInDays int64
}

func (svc *loanService) CreateLoan(ctx context.Context, db *gorm.DB, userID string, params *CreateLoanParams) (*model.Loan, error) {
	loan := &model.Loan{
		UserID:          userID,
		AmountInCents:   params.AmountInCents,
		Term:            params.Term,
		FrequencyInDays: params.FrequencyInDays,
		Status:          model.LoanStatusPending,
	}
	createdLoan, err := svc.loanStore.CreateLoan(ctx, db, loan)
	if err != nil {
		logger.E(ctx, "create loan failed", zap.Error(err))
		return nil, fmt.Errorf("creating loan in service: %w", err)
	}

	return createdLoan, nil
}
