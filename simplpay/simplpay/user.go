package simplpay

import (
	"context"
	"time"

	"github.com/goProjects/simplpay/errorops"
)

// UserStore is the interface that provide basic user related functionality.
type UserStore interface {
	CreateUser(ctx context.Context, u *User) *errorops.Error
	GetUser(ctx context.Context, name string) (*User, *errorops.Error)
	// UpdateUserOwnedAmount(ctx context.Context, id string, ownedAmount int, updatedOn time.Time) *errorops.Error
	// GetUser(ctx context.Context, userID string) (*User, *errorops.Error)
	// GetCreditLimitUsers(ctx context.Context) ([]User, *errorops.Error)
	// GetDueUsers(ctx context.Context) ([]User, *errorops.Error)
}

// User contains details of a user.
type User struct {
	Name        string    `bson:"name"`
	Email       string    `bson:"email"`
	CreditLimit float64   `bson:"credit_limit"`
	OwnedAmount float64   `bson:"owned_amount"`
	CreatedOn   time.Time `bson:"created_on"`
	UpdatedOn   time.Time `bson:"updated_on"`
}

// cleans the request like trims end spaces and casts lower.
func (u *User) sanitize() {

}

// user validation.
func (u *User) validate() {

}

// sets default for a user.
func (u *User) setDefault() {
	u.CreatedOn = Now().UTC()
	u.UpdatedOn = Now().UTC()
}

// CreateUser validates a user and saves in database.
func CreateUser(ctx context.Context, s *Service, u *User) *errorops.Error {
	// sanitizes
	// validate user

	// set defaults
	u.setDefault()

	// store in db
	return s.DB.CreateUser(ctx, u)
}

// Payback is called when a user paybacks his credit full/partial.
func (s *Service) Payback(ctx context.Context) {
	// parse args
	// create user
	// sanitizes
	// validate amount
	// set defaults
	// store in db
}

// GetDues gets the total dues pending for the user to payback.
func (s *Service) GetDues(ctx context.Context) {
	// parse args
	// get user info
	// return ownedAmount
}

// GetCreditLimitUsers gets the users who have reached the credit limit.
func (s *Service) GetCreditLimitUsers(ctx context.Context) {
	// get list of all users where creditLimit == ownedAmount
	// return the list
}

// GetTotalDues gets the total dues pending to be payed backed from all the users.
func (s *Service) GetTotalDues(ctx context.Context) {
	// get list of all users where ownedAmount > 0
	// return the list and total dues
}
