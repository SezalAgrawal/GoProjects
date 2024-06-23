package service

import (
	"context"
	"errors"
	"fmt"

	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/goProjects/loan_app/app/model"
	"github.com/goProjects/loan_app/app/store"
	"github.com/goProjects/loan_app/lib/logger"
)

func (svc *loanService) PayLoan(ctx context.Context, db *gorm.DB, userID, loanID, otsID string, amountInCents int64) (*model.Loan, error) {
	loan, err := svc.loanStore.GetLoanByID(ctx, db, loanID)
	if err != nil {
		if !errors.Is(err, store.ErrNotFound) {
			logger.E(ctx, "error in fetching loan", zap.Error(err), logger.Field("loan_id", loanID))
		}

		return nil, fmt.Errorf("error in fetching loan: %w", err)
	}
	if loan.UserID != userID {
		return nil, fmt.Errorf("%w: loan does not belong to user", ErrInvalidLoan)
	}
	if loan.Status != model.LoanStatusApproved {
		return nil, fmt.Errorf("%w: loan is not in APPROVED state", ErrInvalidLoan)
	}
	if amountInCents > loan.AmountInCents {
		return nil, fmt.Errorf("%w: given amount is greater than loan amount", ErrInvalidLoan)
	}
	if amountInCents > loan.DueAmount() {
		return nil, fmt.Errorf("%w: given amount is greater than loan due amount", ErrInvalidLoan)
	}
	if !verifyPaymentSuccess(otsID) {
		return nil, fmt.Errorf("%w: given ots_id is not paid", ErrInvalidLoan)
	}

	err = db.Transaction(func(tx *gorm.DB) error {
		// get all installments not in paid state order by serial_no asc
		pendingInstallments := loan.GetPendingInstallments()
		paidInstallmentsNum := loan.Term - int64(len(pendingInstallments))

		// iterate over all pending installments and use the given amount to pay as many installments as possible
		for _, inst := range pendingInstallments {
			if amountInCents <= 0 {
				break
			}

			installmentAmountDue := inst.DueAmount()
			inst, err := svc.payInstallment(ctx, tx, loan, inst, amountInCents, otsID)
			if err != nil {
				return err
			}
			if inst.Status == model.InstallmentStatusPaid {
				paidInstallmentsNum++
			}
			amountInCents -= installmentAmountDue
		}

		// if all installments are paid, mark the loan as paid
		if paidInstallmentsNum == loan.Term {
			loan.Status = model.LoanStatusPaid
			if _, err = svc.loanStore.UpdateLoan(ctx, tx, loan); err != nil {
				logger.E(ctx, "update loan failed", zap.Error(err))
				return fmt.Errorf("error in updating loan: %w", err)
			}
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("error in paying loan: %w", err)
	}

	// refresh loan with preloaded installments
	loan, err = svc.loanStore.GetLoanByID(ctx, db, loanID)
	if err != nil {
		logger.E(ctx, "error in fetching loan", zap.Error(err), logger.Field("loan_id", loanID))
		return nil, fmt.Errorf("error in fetching loan: %w", err)
	}
	return loan, nil
}

func verifyPaymentSuccess(otsID string) bool {
	// TODO: call payments_api and check if the ots_id has a corresponding success payment
	// Schema will be something like this:
	// one_time_settlements:
	// id
	// payment_reference
	// amount_in_cents
	// user_id
	// status
	return true
}

func (svc *loanService) payInstallment(ctx context.Context, db *gorm.DB, loan *model.Loan, inst *model.Installment, givenAmount int64, otsID string) (*model.Installment, error) {
	var installmentPaymentAmount int64
	var installmentPaid bool
	installmentAmountDue := inst.DueAmount()
	if installmentAmountDue <= givenAmount {
		// full installment payment
		installmentPaymentAmount = installmentAmountDue
		installmentPaid = true
	} else {
		// partial installment payment
		installmentPaymentAmount = givenAmount
	}

	if _, err := svc.loanStore.CreateInstallmentPayment(ctx, db, &model.InstallmentPayment{
		LoanID:        loan.ID,
		InstallmentID: inst.ID,
		OtsID:         otsID,
		AmountInCents: installmentPaymentAmount,
	}); err != nil {
		logger.E(ctx, "create installment payment failed", zap.Error(err))
		return nil, fmt.Errorf("error in create installment payment: %w", err)
	}

	// if full payment of installment is received, mark installment as paid
	if installmentPaid {
		inst.Status = model.InstallmentStatusPaid
		var err error
		inst, err = svc.loanStore.UpdateInstallment(ctx, db, inst)
		if err != nil {
			logger.E(ctx, "update installment failed", zap.Error(err))
			return nil, fmt.Errorf("error in update installment: %w", err)
		}
	}

	return inst, nil
}
