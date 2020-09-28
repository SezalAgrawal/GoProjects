package simplpay

import (
	"context"
	"testing"

	"github.com/goProjects/simplpay/errorops"
	"github.com/stretchr/testify/assert"
)

func TestCreateUser(t *testing.T) {
	cases := []struct {
		name    string
		user    *User
		wantErr *errorops.Error
	}{
		{
			name: "created user",
			user: &User{
				Name: "abc",
			},
			wantErr: nil,
		},
		{
			name: "duplicate user",
			user: &User{
				Name: "abc",
			},
			wantErr: &errorops.Error{
				Code: errorops.Conflict,
			},
		},
	}

	s := &Service{
		DB: &mockDB{},
	}

	for _, c := range cases {
		gotErr := CreateUser(context.Background(), s, c.user)
		assert.Equal(t, c.wantErr, gotErr)
	}
}
