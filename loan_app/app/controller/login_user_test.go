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

func TestLoginUser(t *testing.T) {

	t.Run("mandatory fields missing", func(t *testing.T) {
		user := map[string]any{
			"name":     "",
			"password": "",
		}

		expected := map[string]any{
			"api_version": "1.0",
			"error": map[string]any{
				"code":    "bad_request",
				"message": "InvalidValue: name should not be empty, password should not be empty",
			},
			"success": false,
		}

		responseCode, responseBody := sendRequest(t, http.MethodPost, "/api/v1/users/login", user, nil, nil)

		assert.Equal(t, http.StatusBadRequest, responseCode)
		assert.Equal(t, expected, responseBody)
	})

	t.Run("return error when password is invalid", func(t *testing.T) {
		test_helper.ClearDataFromPostgres()
		assert.Equal(t, nil, db.Get().Create(&model.Role{ID: "rol_123", Name: "USER"}).Error)
		_, err := service.NewUserService().CreateUser(context.Background(), db.Get(), "john", "john@123", "rol_123")
		assert.Nil(t, err)

		login := map[string]any{
			"name":     "john",
			"password": "john123",
		}

		expected := map[string]any{
			"api_version": "1.0",
			"error": map[string]any{
				"code":    "not_found",
				"message": "unable to fetch user: record not found",
			},
			"success": false,
		}

		responseCode, responseBody := sendRequest(t, http.MethodPost, "/api/v1/users/login", login, nil, nil)

		assert.Equal(t, http.StatusNotFound, responseCode)
		assert.Equal(t, expected, responseBody)
	})

	t.Run("successfully login user", func(t *testing.T) {
		test_helper.ClearDataFromPostgres()
		assert.Equal(t, nil, db.Get().Create(&model.Role{ID: "rol_123", Name: "USER"}).Error)
		user, err := service.NewUserService().CreateUser(context.Background(), db.Get(), "john", "john@123", "rol_123")
		assert.Nil(t, err)

		login := map[string]any{
			"name":     "john",
			"password": "john@123",
		}

		expected := map[string]any{
			"api_version": "1.0",
			"data": map[string]any{
				"user": map[string]any{
					"id": user.ID,
				},
			},
			"success": true,
		}

		responseCode, responseBody := sendRequest(t, http.MethodPost, "/api/v1/users/login", login, nil, nil)

		assert.Equal(t, http.StatusOK, responseCode)
		test_helper.AssertEqualMap(t, expected, responseBody, []string{"access_token"})
	})
}
