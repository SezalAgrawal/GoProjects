package simplpay

import (
	"context"
	"time"

	"github.com/goProjects/simplpay/errorops"
)

// MerchantStore is the interface that provide basic merchant related functionality.
type MerchantStore interface {
	CreateMerchant(ctx context.Context, m *Merchant) *errorops.Error
	GetMerchant(ctx context.Context, name string) (*Merchant, *errorops.Error)
	//UpdateMerchantDiscountPercent(ctx context.Context, id string, discountPercent int, updatedOn time.Time) *errorops.Error
}

// Merchant contains details of a merchant.
type Merchant struct {
	Name            string    `bson:"name"`
	DiscountPercent float64   `bson:"discount_percent"`
	CreatedOn       time.Time `bson:"created_on"`
	UpdatedOn       time.Time `bson:"updated_on"`
}

// cleans the request like trims end spaces and casts lower.
func (m *Merchant) sanitize() {

}

// merchant validation.
func (m *Merchant) validate() {

}

// sets default for a merchant.
func (m *Merchant) setDefault() {

}

// CreateMerchant validates a merchant and saves in database.
func CreateMerchant(ctx context.Context, s *Service, m *Merchant) *errorops.Error {
	// sanitizes
	// validate merchant

	// set defaults
	m.CreatedOn = Now().UTC()
	m.UpdatedOn = Now().UTC()

	// store in db
	return s.DB.CreateMerchant(ctx, m)
}

// UpdateMerchant updates a merchant and saves in database.
func UpdateMerchant(ctx context.Context) {
	// parse args
	// create merchant
	// sanitizes
	// validate merchant
	// set defaults > set the updateOn time
	// store in db
}
