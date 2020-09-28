package db

import (
	"context"
	"testing"

	"github.com/goProjects/simplpay/errorops"
	"github.com/goProjects/simplpay/simplpay"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateUser(t *testing.T) {
	teardownFn(t, testDBName, beforeEachDBClean)
	cases := []struct {
		name    string
		user    *simplpay.User
		wantErr string
	}{
		{
			name: "Successfully created user",
			user: &simplpay.User{
				Name: "abc",
			},
			wantErr: "",
		},
		{
			name: "Duplicate user",
			user: &simplpay.User{
				Name: "abc",
			},
			wantErr: errorops.Conflict,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ctx := context.Background()
			gotErr := testDB.CreateUser(ctx, c.user)
			if c.wantErr != "" {
				require.NotNil(t, gotErr)
				assert.Equal(t, c.wantErr, gotErr.Code)
			} else {
				require.Nil(t, gotErr)
			}
		})
	}
}
