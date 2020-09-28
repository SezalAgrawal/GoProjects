package simplpay

import (
	"context"
	"time"

	"github.com/goProjects/simplpay/errorops"
)

var (
	usersCollection = make(map[string]*User)
)

type mockDB struct{}

func cleanUsersCollection() {
	usersCollection = make(map[string]*User)
}

func (db *mockDB) CreateUser(ctx context.Context, u *User) *errorops.Error {
	for _, u := range usersCollection {
		if u.Name == u.Name {
			return &errorops.Error{Code: errorops.Conflict}
		}
	}
	usersCollection[u.Name] = u
	return nil
}

func (db *mockDB) CreateMerchant(ctx context.Context, m *Merchant) *errorops.Error {
	return nil
}

func (db *mockDB) CreateTransaction(ctx context.Context, t *Transaction, userOwnedAmount float64, updatedOn time.Time) *errorops.Error {
	return nil
}

func (db *mockDB) GetUser(ctx context.Context, name string) (*User, *errorops.Error) {
	return nil, nil
}

func (db *mockDB) GetMerchant(ctx context.Context, name string) (*Merchant, *errorops.Error) {
	return nil, nil
}
