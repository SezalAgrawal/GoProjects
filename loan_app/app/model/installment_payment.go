package model

import "time"

type InstallmentPayment struct {
	ID            string
	LoanID        string
	InstallmentID string
	AmountInCents int64
	OtsID         string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
