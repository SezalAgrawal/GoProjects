package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/goProjects/loan_app/app/model"
	"github.com/goProjects/loan_app/app/store"
	"github.com/goProjects/loan_app/lib/logger"
)

var (
	ErrInvalidLoan = errors.New("loan validation error")
)

func (svc *loanService) ApproveLoan(ctx context.Context, db *gorm.DB, loanID string) (*model.Loan, error) {
	loan, err := svc.loanStore.GetLoanByID(ctx, db, loanID)
	if err != nil {
		if !errors.Is(err, store.ErrNotFound) {
			logger.E(ctx, "error in fetching loan", zap.Error(err), logger.Field("loan_id", loanID))
		}

		return nil, fmt.Errorf("error in fetching loan: %w", err)
	}
	if loan.Status != model.LoanStatusPending {
		return nil, fmt.Errorf("%w: loan is not in PENDING state", ErrInvalidLoan)
	}

	err = db.Transaction(func(tx *gorm.DB) error {
		loan.Status = model.LoanStatusApproved
		loan, err = svc.loanStore.UpdateLoan(ctx, tx, loan)
		if err != nil {
			logger.E(ctx, "update loan failed", zap.Error(err))
			return fmt.Errorf("error in updating loan: %w", err)
		}

		installments := calculateInstallmentSchedule(loan.AmountInCents, loan.Term, loan.FrequencyInDays)
		for _, inst := range installments {
			inst.Status = model.InstallmentStatusPending
			inst.LoanID = loan.ID
			if _, err := svc.loanStore.CreateInstallment(ctx, tx, &inst); err != nil {
				logger.E(ctx, "create installment failed", zap.Error(err))
				return fmt.Errorf("error in creating loan: %w", err)
			}
		}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("error in approving loan: %w", err)
	}

	// refresh loan with preloaded installments
	loan, err = svc.loanStore.GetLoanByID(ctx, db, loanID)
	if err != nil {
		logger.E(ctx, "error in fetching loan", zap.Error(err), logger.Field("loan_id", loanID))
		return nil, fmt.Errorf("error in fetching loan: %w", err)
	}
	return loan, nil
}

func calculateInstallmentSchedule(amount int64, term int64, frequencyInDays int64) []model.Installment {
	installments := make([]model.Installment, term)
	installmentAmount := amount / term
	remainder := int64(amount) % term

	var i int64
	for i = 0; i < term; i++ {
		installments[i].AmountInCents = installmentAmount
		if i >= term-remainder {
			installments[i].AmountInCents += 1
		}
		installments[i].DueDate = time.Now().AddDate(0, 0, int((i+1)*frequencyInDays))
		installments[i].SerialNo = i + 1
	}

	return installments
}
