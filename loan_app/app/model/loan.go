package model

import (
	"sort"
	"time"
)

type LoanStatus string

const (
	LoanStatusPending  LoanStatus = "PENDING"
	LoanStatusApproved LoanStatus = "APPROVED"
	LoanStatusPaid     LoanStatus = "PAID"
)

type Loan struct {
	ID              string
	UserID          string
	AmountInCents   int64
	Term            int64
	FrequencyInDays int64
	Status          LoanStatus
	CreatedAt       time.Time
	UpdatedAt       time.Time

	Installments []*Installment
}

func (l *Loan) DueAmount() int64 {
	var amount int64
	for _, inst := range l.Installments {
		amount += inst.DueAmount()
	}
	return amount
}

func (l *Loan) GetPendingInstallments() []*Installment {
	var pendingInstallments []*Installment
	for _, inst := range l.Installments {
		if inst.Status != InstallmentStatusPaid {
			pendingInstallments = append(pendingInstallments, inst)
		}
	}

	sort.Slice(pendingInstallments, func(i, j int) bool {
		// sort the installments in ascending order of serial_no
		return pendingInstallments[i].SerialNo < pendingInstallments[j].SerialNo
	})
	return pendingInstallments
}
