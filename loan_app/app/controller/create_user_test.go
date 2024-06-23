package controller_test

import (
	"net/http"
	"testing"

	"github.com/goProjects/loan_app/app/model"
	"github.com/goProjects/loan_app/app/test_helper"
	"github.com/goProjects/loan_app/lib/db"
	"github.com/stretchr/testify/assert"
)

func TestCreateUser(t *testing.T) {

	t.Run("mandatory fields missing", func(t *testing.T) {
		user := map[string]any{
			"name":     "",
			"password": "",
			"role_id":  "",
		}

		expected := map[string]any{
			"api_version": "1.0",
			"error": map[string]any{
				"code":    "bad_request",
				"message": "InvalidValue: name should not be empty, password should not be empty, role_id should not be empty",
			},
			"success": false,
		}

		responseCode, responseBody := sendRequest(t, http.MethodPost, "/api/v1/users", user, nil, nil)

		assert.Equal(t, http.StatusBadRequest, responseCode)
		assert.Equal(t, expected, responseBody)
	})

	t.Run("successfully created user", func(t *testing.T) {
		test_helper.ClearDataFromPostgres()
		assert.Equal(t, nil, db.Get().Create(&model.Role{ID: "rol_123", Name: "USER"}).Error)
		user := map[string]any{
			"name":     "john",
			"password": "john@123",
			"role_id":  "rol_123",
		}

		expected := map[string]any{
			"api_version": "1.0",
			"data": map[string]any{
				"user": map[string]any{
					"name": "john",
				},
			},
			"success": true,
		}

		responseCode, responseBody := sendRequest(t, http.MethodPost, "/api/v1/users", user, nil, nil)

		assert.Equal(t, http.StatusOK, responseCode)
		test_helper.AssertEqualMap(t, expected, responseBody, []string{"id", "created_at", "updated_at"})
	})
}
