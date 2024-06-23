package controller

import (
	"net/http"
	"time"

	"github.com/goProjects/loan_app/app/model"
	"github.com/goProjects/loan_app/app/service"
	"github.com/goProjects/loan_app/lib/db"
	"github.com/goProjects/loan_app/lib/utils"
	"github.com/goProjects/loan_app/lib/web"
)

type getLoanResponse struct {
	Loans []loan `json:"loans"`
}

type loan struct {
	ID              string         `json:"id"`
	AmountInCents   int64          `json:"amount_in_cents"`
	Term            int64          `json:"term"`
	FrequencyInDays int64          `json:"frequency_in_days"`
	Status          string         `json:"status"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	Installments    []*installment `json:"installments"`
}

type installment struct {
	ID                  string                `json:"id"`
	LoanID              string                `json:"loan_id"`
	AmountInCents       int64                 `json:"amount_in_cents"`
	Status              string                `json:"status"`
	SerialNo            int64                 `json:"serial_no"`
	DueDate             time.Time             `json:"due_date"`
	CreatedAt           time.Time             `json:"created_at"`
	UpdatedAt           time.Time             `json:"updated_at"`
	InstallmentPayments []*installmentPayment `json:"payments"`
}

type installmentPayment struct {
	ID            string    `json:"id"`
	LoanID        string    `json:"loan_id"`
	InstallmentID string    `json:"installment_id"`
	OtsID         string    `json:"one_time_settlement_id"`
	AmountInCents int64     `json:"amount_in_cents"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func GetAllLoans(request *web.Request) *web.APIResponse {
	userID := request.Value(utils.CurrentUserIDStoreKey).(string)

	loans, err := service.NewLoanService().GetAllLoans(request.Context(), db.Get(), userID)
	if err != nil {
		return web.ErrInternalServerError
	}

	loansResp := make([]loan, 0)
	for _, l := range loans {
		loansResp = append(loansResp, convertLoanModelToResp(l))
	}
	response := getLoanResponse{Loans: loansResp}
	return web.NewAPISuccessResponse(response, http.StatusOK)
}

func convertLoanModelToResp(l *model.Loan) loan {
	respLoan := loan{
		ID:              l.ID,
		AmountInCents:   l.AmountInCents,
		Term:            l.Term,
		FrequencyInDays: l.FrequencyInDays,
		Status:          string(l.Status),
		CreatedAt:       l.CreatedAt,
		UpdatedAt:       l.UpdatedAt,
	}

	var installments []*installment
	for _, inst := range l.Installments {
		installments = append(installments, convertInstallmentModelToResp(inst))
	}
	respLoan.Installments = installments
	return respLoan
}

func convertInstallmentModelToResp(e *model.Installment) *installment {
	respInstallment := &installment{
		ID:            e.ID,
		LoanID:        e.LoanID,
		AmountInCents: e.AmountInCents,
		Status:        string(e.Status),
		SerialNo:      e.SerialNo,
		DueDate:       e.DueDate,
		CreatedAt:     e.CreatedAt,
		UpdatedAt:     e.UpdatedAt,
	}
	var installmentPayments []*installmentPayment
	for _, pay := range e.InstallmentPayments {
		installmentPayments = append(installmentPayments, convertInstallmentPaymentModelToResp(pay))
	}
	respInstallment.InstallmentPayments = installmentPayments

	return respInstallment
}

func convertInstallmentPaymentModelToResp(e *model.InstallmentPayment) *installmentPayment {
	return &installmentPayment{
		ID:            e.ID,
		LoanID:        e.LoanID,
		InstallmentID: e.InstallmentID,
		OtsID:         e.OtsID,
		AmountInCents: e.AmountInCents,
		CreatedAt:     e.CreatedAt,
		UpdatedAt:     e.UpdatedAt,
	}
}
