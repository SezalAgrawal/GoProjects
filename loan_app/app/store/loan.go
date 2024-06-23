package store

import (
	"context"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"

	"github.com/goProjects/loan_app/app/model"
	"github.com/goProjects/loan_app/lib/utils"
)

type LoanStore interface {
	CreateLoan(ctx context.Context, db *gorm.DB, loan *model.Loan) (*model.Loan, error)
	GetAllLoans(ctx context.Context, db *gorm.DB, userID string) ([]*model.Loan, error)
	GetLoanByID(ctx context.Context, db *gorm.DB, loanID string) (*model.Loan, error)
	UpdateLoan(ctx context.Context, db *gorm.DB, loan *model.Loan) (*model.Loan, error)
	CreateInstallment(ctx context.Context, db *gorm.DB, installment *model.Installment) (*model.Installment, error)
	UpdateInstallment(ctx context.Context, db *gorm.DB, installment *model.Installment) (*model.Installment, error)
	CreateInstallmentPayment(ctx context.Context, db *gorm.DB, pay *model.InstallmentPayment) (*model.InstallmentPayment, error)
}

type loanStore struct {
}

func NewLoanStore() LoanStore {
	return &loanStore{}
}

func (e *loanStore) CreateLoan(ctx context.Context, db *gorm.DB, loan *model.Loan) (*model.Loan, error) {
	l := loanModelToDB(loan)

	if err := db.WithContext(ctx).Create(&l).Error; err != nil {
		return nil, fmt.Errorf("creating loan in store %w", err)
	}

	return l.dbToModel(), nil
}

func (e *loanStore) GetAllLoans(ctx context.Context, db *gorm.DB, userID string) ([]*model.Loan, error) {
	var dbLoans []*loan
	var modelLoans []*model.Loan

	if err := db.WithContext(ctx).Model(&loan{}).
		Preload("Installments").
		Preload("Installments.InstallmentPayments").
		Where("user_id = ?", userID).Find(&dbLoans).Error; err != nil {
		return nil, err
	}

	for _, dbLoan := range dbLoans {
		modelLoans = append(modelLoans, dbLoan.dbToModel())
	}

	return modelLoans, nil
}

func (e *loanStore) GetLoanByID(ctx context.Context, db *gorm.DB, loanID string) (*model.Loan, error) {
	l := new(loan)

	if err := db.WithContext(ctx).Model(&loan{}).
		Preload("Installments").
		Preload("Installments.InstallmentPayments").
		Where("id = ?", loanID).
		First(l).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}

		return nil, fmt.Errorf("fetching loan in store %w", err)
	}

	return l.dbToModel(), nil
}

func (e *loanStore) UpdateLoan(ctx context.Context, db *gorm.DB, l *model.Loan) (*model.Loan, error) {
	loanDB := loanModelToDB(l)

	if err := db.WithContext(ctx).Updates(loanDB).Error; err != nil {
		return nil, fmt.Errorf("updating loan in store %w", err)
	}

	return e.GetLoanByID(ctx, db, l.ID)
}

func (e *loanStore) CreateInstallment(ctx context.Context, db *gorm.DB, inst *model.Installment) (*model.Installment, error) {
	i := installmentModelToDB(inst)

	if err := db.WithContext(ctx).Create(&i).Error; err != nil {
		return nil, fmt.Errorf("creating installment in store %w", err)
	}

	return i.dbToModel(), nil
}

func (e *loanStore) UpdateInstallment(ctx context.Context, db *gorm.DB, inst *model.Installment) (*model.Installment, error) {
	instDB := installmentModelToDB(inst)

	if err := db.WithContext(ctx).Updates(instDB).Error; err != nil {
		return nil, fmt.Errorf("updating installment in store %w", err)
	}

	return instDB.dbToModel(), nil
}

func (e *loanStore) CreateInstallmentPayment(ctx context.Context, db *gorm.DB, pay *model.InstallmentPayment) (*model.InstallmentPayment, error) {
	i := installmentPaymentModelToDB(pay)

	if err := db.WithContext(ctx).Create(&i).Error; err != nil {
		return nil, fmt.Errorf("creating installment payment in store %w", err)
	}

	return i.dbToModel(), nil
}

type loan struct {
	ID              string         `gorm:"column:id"`
	UserID          string         `gorm:"column:user_id"`
	AmountInCents   int64          `gorm:"column:amount_in_cents"`
	Term            int64          `gorm:"column:term"`
	FrequencyInDays int64          `gorm:"column:frequency_in_days"`
	Status          string         `gorm:"column:status"`
	CreatedAt       time.Time      `gorm:"column:created_at"`
	UpdatedAt       time.Time      `gorm:"column:updated_at"`
	Installments    []*installment `gorm:"foreignkey:LoanID"`
}

type installment struct {
	ID                  string                `gorm:"column:id"`
	LoanID              string                `gorm:"column:loan_id"`
	AmountInCents       int64                 `gorm:"column:amount_in_cents"`
	Status              string                `gorm:"column:status"`
	SerialNo            int64                 `gorm:"column:serial_no"`
	DueDate             time.Time             `gorm:"column:due_date"`
	CreatedAt           time.Time             `gorm:"column:created_at"`
	UpdatedAt           time.Time             `gorm:"column:updated_at"`
	InstallmentPayments []*installmentPayment `gorm:"foreignkey:InstallmentID"`
}

type installmentPayment struct {
	ID            string    `gorm:"column:id"`
	LoanID        string    `gorm:"column:loan_id"`
	InstallmentID string    `gorm:"column:installment_id"`
	OtsID         string    `gorm:"column:one_time_settlement_id"`
	AmountInCents int64     `gorm:"column:amount_in_cents"`
	CreatedAt     time.Time `gorm:"column:created_at"`
	UpdatedAt     time.Time `gorm:"column:updated_at"`
}

func (e *loan) BeforeCreate(_ *gorm.DB) error {
	e.ID = "loan_" + utils.NewKSUID()
	return nil
}

func (e *installment) BeforeCreate(_ *gorm.DB) error {
	e.ID = "inst_" + utils.NewKSUID()
	return nil
}

func (e *installmentPayment) BeforeCreate(_ *gorm.DB) error {
	e.ID = "pay_" + utils.NewKSUID()
	return nil
}

func (e *loan) dbToModel() *model.Loan {
	loan := &model.Loan{
		ID:              e.ID,
		UserID:          e.UserID,
		AmountInCents:   e.AmountInCents,
		Term:            e.Term,
		FrequencyInDays: e.FrequencyInDays,
		Status:          model.LoanStatus(e.Status),
		CreatedAt:       e.CreatedAt,
		UpdatedAt:       e.UpdatedAt,
	}
	var installments []*model.Installment
	for _, inst := range e.Installments {
		installments = append(installments, inst.dbToModel())
	}
	loan.Installments = installments

	return loan
}

func loanModelToDB(model *model.Loan) *loan {
	return &loan{
		ID:              model.ID,
		UserID:          model.UserID,
		AmountInCents:   model.AmountInCents,
		Term:            model.Term,
		FrequencyInDays: model.FrequencyInDays,
		Status:          string(model.Status),
		CreatedAt:       model.CreatedAt,
		UpdatedAt:       model.UpdatedAt,
	}
}

func (e *installment) dbToModel() *model.Installment {
	inst := &model.Installment{
		ID:            e.ID,
		LoanID:        e.LoanID,
		AmountInCents: e.AmountInCents,
		Status:        model.InstallmentStatus(e.Status),
		SerialNo:      e.SerialNo,
		DueDate:       e.DueDate,
		CreatedAt:     e.CreatedAt,
		UpdatedAt:     e.UpdatedAt,
	}
	var installmentPayments []*model.InstallmentPayment
	for _, pay := range e.InstallmentPayments {
		installmentPayments = append(installmentPayments, pay.dbToModel())
	}
	inst.InstallmentPayments = installmentPayments

	return inst
}

func installmentModelToDB(model *model.Installment) *installment {
	return &installment{
		ID:            model.ID,
		LoanID:        model.LoanID,
		AmountInCents: model.AmountInCents,
		Status:        string(model.Status),
		SerialNo:      model.SerialNo,
		DueDate:       model.DueDate,
		CreatedAt:     model.CreatedAt,
		UpdatedAt:     model.UpdatedAt,
	}
}

func (e *installmentPayment) dbToModel() *model.InstallmentPayment {
	return &model.InstallmentPayment{
		ID:            e.ID,
		LoanID:        e.LoanID,
		InstallmentID: e.InstallmentID,
		OtsID:         e.OtsID,
		AmountInCents: e.AmountInCents,
		CreatedAt:     e.CreatedAt,
		UpdatedAt:     e.UpdatedAt,
	}
}

func installmentPaymentModelToDB(model *model.InstallmentPayment) *installmentPayment {
	return &installmentPayment{
		ID:            model.ID,
		LoanID:        model.LoanID,
		InstallmentID: model.InstallmentID,
		OtsID:         model.OtsID,
		AmountInCents: model.AmountInCents,
		CreatedAt:     model.CreatedAt,
		UpdatedAt:     model.UpdatedAt,
	}
}
