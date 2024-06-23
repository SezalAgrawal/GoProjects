package controller_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/goProjects/loan_app/app/model"
	"github.com/goProjects/loan_app/app/service"
	"github.com/goProjects/loan_app/app/test_helper"
	"github.com/goProjects/loan_app/lib/db"
	"github.com/stretchr/testify/assert"
)

func TestPayLoan(t *testing.T) {

	// TODO: add test for loan not found and invalid loan and partial and full payments

	t.Run("pay loan", func(t *testing.T) {
		test_helper.ClearDataFromPostgres()
		assert.Equal(t, nil, db.Get().Create(&model.Role{ID: "rol_123", Name: "USER"}).Error)
		user, err := service.NewUserService().CreateUser(context.Background(), db.Get(), "john", "john@123", "rol_123")
		assert.Nil(t, err)
		_, accessToken, err := service.NewUserService().GenerateAccessToken(context.Background(), db.Get(), "john", "john@123")
		assert.Nil(t, err)
		loan, err := service.NewLoanService().CreateLoan(context.Background(), db.Get(), user.ID, &service.CreateLoanParams{
			AmountInCents:   100,
			Term:            3,
			FrequencyInDays: 7,
		})
		assert.Nil(t, err)
		loan, err = service.NewLoanService().ApproveLoan(context.Background(), db.Get(), loan.ID)
		assert.Nil(t, err)

		pay := map[string]any{
			"one_time_settlement_id": "ots_123",
			"amount_in_cents":        50,
		}

		expected := map[string]any{
			"api_version": "1.0",
			"data": map[string]any{
				"loan": map[string]any{
					"id":                loan.ID,
					"amount_in_cents":   float64(100),
					"term":              float64(3),
					"frequency_in_days": float64(7),
					"status":            "APPROVED",
					"installments": []any{
						map[string]any{
							"loan_id":         loan.ID,
							"id":              loan.Installments[0].ID,
							"amount_in_cents": float64(33),
							"status":          "PAID",
							"serial_no":       float64(1),
							"payments": []any{
								map[string]any{
									"loan_id":                loan.ID,
									"installment_id":         loan.Installments[0].ID,
									"one_time_settlement_id": "ots_123",
									"amount_in_cents":        float64(33),
								},
							},
						},
						map[string]any{
							"loan_id":         loan.ID,
							"id":              loan.Installments[1].ID,
							"amount_in_cents": float64(33),
							"status":          "PENDING",
							"serial_no":       float64(2),
							"payments": []any{
								map[string]any{
									"loan_id":                loan.ID,
									"installment_id":         loan.Installments[1].ID,
									"one_time_settlement_id": "ots_123",
									"amount_in_cents":        float64(17),
								},
							},
						},
						map[string]any{
							"loan_id":         loan.ID,
							"id":              loan.Installments[2].ID,
							"amount_in_cents": float64(34),
							"status":          "PENDING",
							"serial_no":       float64(3),
							"payments":        nil,
						},
					},
				},
			},
			"success": true,
		}

		responseCode, responseBody := sendRequest(t, http.MethodPost, fmt.Sprintf("/api/v1/loans/%s/payment", loan.ID), pay, nil, map[string]string{"Access-Token": accessToken})

		assert.Equal(t, http.StatusOK, responseCode)
		test_helper.AssertEqualMap(t, expected, responseBody, []string{"id", "created_at", "updated_at", "due_date"})
	})
}
