package controller_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/goProjects/loan_app/app/model"
	"github.com/goProjects/loan_app/app/service"
	"github.com/goProjects/loan_app/app/test_helper"
	"github.com/goProjects/loan_app/lib/db"
	"github.com/stretchr/testify/assert"
)

func TestCreateLoan(t *testing.T) {

	t.Run("mandatory fields missing", func(t *testing.T) {
		test_helper.ClearDataFromPostgres()
		assert.Equal(t, nil, db.Get().Create(&model.Role{ID: "rol_123", Name: "USER"}).Error)
		_, err := service.NewUserService().CreateUser(context.Background(), db.Get(), "john", "john@123", "rol_123")
		assert.Nil(t, err)
		_, accessToken, err := service.NewUserService().GenerateAccessToken(context.Background(), db.Get(), "john", "john@123")
		assert.Nil(t, err)

		loan := map[string]any{
			"amount_in_cents":   0,
			"term":              0,
			"frequency_in_days": 0,
		}

		expected := map[string]any{
			"api_version": "1.0",
			"error": map[string]any{
				"code":    "bad_request",
				"message": "InvalidValue: amount_in_cents is a required field, term is a required field, frequency_in_days is a required field",
			},
			"success": false,
		}

		responseCode, responseBody := sendRequest(t, http.MethodPost, "/api/v1/loans", loan, nil, map[string]string{"Access-Token": accessToken})

		assert.Equal(t, http.StatusBadRequest, responseCode)
		assert.Equal(t, expected, responseBody)
	})

	t.Run("successfully created loan", func(t *testing.T) {
		test_helper.ClearDataFromPostgres()
		assert.Equal(t, nil, db.Get().Create(&model.Role{ID: "rol_123", Name: "USER"}).Error)
		_, err := service.NewUserService().CreateUser(context.Background(), db.Get(), "john", "john@123", "rol_123")
		assert.Nil(t, err)
		_, accessToken, err := service.NewUserService().GenerateAccessToken(context.Background(), db.Get(), "john", "john@123")
		assert.Nil(t, err)

		loan := map[string]any{
			"amount_in_cents":   100,
			"term":              3,
			"frequency_in_days": 7,
		}

		expected := map[string]any{
			"api_version": "1.0",
			"data": map[string]any{
				"loan": map[string]any{
					"amount_in_cents":   float64(100),
					"term":              float64(3),
					"frequency_in_days": float64(7),
					"status":            "PENDING",
					"installments":      nil,
				},
			},
			"success": true,
		}

		responseCode, responseBody := sendRequest(t, http.MethodPost, "/api/v1/loans", loan, nil, map[string]string{"Access-Token": accessToken})

		assert.Equal(t, http.StatusOK, responseCode)
		test_helper.AssertEqualMap(t, expected, responseBody, []string{"id", "created_at", "updated_at"})
	})
}
