package model

import "time"

type InstallmentStatus string

const (
	InstallmentStatusPending        InstallmentStatus = "PENDING"
	InstallmentStatusPaymentPending InstallmentStatus = "PAYMENT_PENDING"
	InstallmentStatusPaid           InstallmentStatus = "PAID"
)

type Installment struct {
	ID            string
	LoanID        string
	AmountInCents int64
	Status        InstallmentStatus
	SerialNo      int64
	DueDate       time.Time
	CreatedAt     time.Time
	UpdatedAt     time.Time

	InstallmentPayments []*InstallmentPayment
}

func (i *Installment) DueAmount() int64 {
	amount := i.AmountInCents
	for _, pay := range i.InstallmentPayments {
		amount -= pay.AmountInCents
	}
	return amount
}
