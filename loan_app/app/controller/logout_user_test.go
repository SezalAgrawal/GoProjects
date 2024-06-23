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

func TestLogoutUser(t *testing.T) {

	t.Run("return error when access token is invalid", func(t *testing.T) {
		expected := map[string]any{
			"api_version": "1.0",
			"error": map[string]any{
				"code":    "unauthorized",
				"message": "unauthorized",
			},
			"success": false,
		}

		responseCode, responseBody := sendRequest(t, http.MethodDelete, "/api/v1/users/logout", nil, nil, nil)

		assert.Equal(t, http.StatusUnauthorized, responseCode)
		assert.Equal(t, expected, responseBody)
	})

	t.Run("delete all access tokens of the user", func(t *testing.T) {
		test_helper.ClearDataFromPostgres()
		assert.Equal(t, nil, db.Get().Create(&model.Role{ID: "rol_123", Name: "USER"}).Error)
		_, err := service.NewUserService().CreateUser(context.Background(), db.Get(), "john", "john@123", "rol_123")
		assert.Nil(t, err)
		_, accessToken, err := service.NewUserService().GenerateAccessToken(context.Background(), db.Get(), "john", "john@123")
		assert.Nil(t, err)

		expected := map[string]any{
			"api_version": "1.0",
			"success":     true,
		}

		responseCode, responseBody := sendRequest(t, http.MethodDelete, "/api/v1/users/logout", nil, nil, map[string]string{"Access-Token": accessToken})

		assert.Equal(t, http.StatusOK, responseCode)
		assert.Equal(t, expected, responseBody)

		// re-request with old access token fails
		expected = map[string]any{
			"api_version": "1.0",
			"error": map[string]any{
				"code":    "unauthorized",
				"message": "unauthorized",
			},
			"success": false,
		}

		responseCode, responseBody = sendRequest(t, http.MethodDelete, "/api/v1/users/logout", nil, nil, map[string]string{"ACCESS_TOKEN": accessToken})

		assert.Equal(t, http.StatusUnauthorized, responseCode)
		assert.Equal(t, expected, responseBody)
	})
}
