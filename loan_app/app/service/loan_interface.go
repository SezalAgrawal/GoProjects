package service

import (
	"context"

	"github.com/goProjects/loan_app/app/model"
	"github.com/goProjects/loan_app/app/store"
	"gorm.io/gorm"
)

type LoanInterface interface {
	CreateLoan(ctx context.Context, db *gorm.DB, userID string, params *CreateLoanParams) (*model.Loan, error)
	GetAllLoans(ctx context.Context, db *gorm.DB, userID string) ([]*model.Loan, error)
	ApproveLoan(ctx context.Context, db *gorm.DB, loanID string) (*model.Loan, error)
	PayLoan(ctx context.Context, db *gorm.DB, userID, loanID, otsID string, amountInCents int64) (*model.Loan, error)
}

type loanService struct {
	loanStore store.LoanStore
}

func NewLoanService() LoanInterface {
	svc := &loanService{
		loanStore: store.NewLoanStore(),
	}
	return svc
}
