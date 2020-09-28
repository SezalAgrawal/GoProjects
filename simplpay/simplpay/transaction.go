package simplpay

import (
	"context"
	"time"

	"github.com/goProjects/simplpay/errorops"
	"github.com/google/uuid"
)

// TransactionStore is the interface that provide basic transaction related functionality.
type TransactionStore interface {
	CreateTransaction(ctx context.Context, t *Transaction, userOwnedAmount float64, updatedOn time.Time) *errorops.Error
	//GetMerchantDiscountAmount(ctx context.Context, merchantID string) (float64, *errorops.Error)
}

// Transaction contains details of a tx.
type Transaction struct {
	ID             string    `bson:"id"`
	UserName       string    `bson:"user_name"`
	MerchantName   string    `bson:"merchant_name"`
	Amount         float64   `bson:"amount"`
	DiscountAmount float64   `bson:"discount_amount"`
	CreatedOn      time.Time `bson:"created_on"`
}

// cleans the request like trims end spaces and casts lower.
func (t *Transaction) sanitize() {

}

// transaction validation.
func (t *Transaction) validate() {

}

// sets default for a transaction.
func (t *Transaction) setDefault() {

}

// CreateTransaction validates a tx and saves in database.
func CreateTransaction(ctx context.Context, s *Service, t *Transaction) *errorops.Error {
	// parse args
	// create transaction
	// sanitizes
	// validate transaction

	// set defaults
	t.ID = uuid.New().String()
	t.CreatedOn = Now().UTC()
	// set discount amount
	merchant, err := s.DB.GetMerchant(ctx, t.MerchantName)
	if err != nil {
		return err
	}
	discount := merchant.DiscountPercent
	t.DiscountAmount = t.Amount - (t.Amount*discount)/100

	// set updated owned amount for the user
	user, err := s.DB.GetUser(ctx, t.UserName)
	if err != nil {
		return err
	}
	ownedAmount := user.OwnedAmount + t.Amount

	// store in db
	return s.DB.CreateTransaction(ctx, t, ownedAmount, Now().UTC())
}

// GetMerchantDiscountAmount gets the total discount received from a merchant.
func GetMerchantDiscountAmount(ctx context.Context) {
	// parse args
	// get list of all txs of that merchantID
	// add the discountAmount
	// return the result
}
