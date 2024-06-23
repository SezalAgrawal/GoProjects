package controller

import (
	"net/http"
	"time"

	"github.com/goProjects/loan_app/app/service"
	"github.com/goProjects/loan_app/lib/db"
	"github.com/goProjects/loan_app/lib/web"
)

type createUserRequestParams struct {
	Name     string `json:"name" validate:"notblank"`
	Password string `json:"password" validate:"notblank"`
	RoleID   string `json:"role_id" validate:"notblank"`
}

type createUserResponse struct {
	User struct {
		ID        string    `json:"id"`
		Name      string    `json:"name"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	} `json:"user"`
}

func CreateUser(request *web.Request) *web.APIResponse {
	data := &createUserRequestParams{}

	if err := request.ValidateBodyToStruct(data); err != nil {
		return web.ErrBadRequest(err.Error())
	}

	user, err := service.NewUserService().CreateUser(request.Context(), db.Get(), data.Name, data.Password, data.RoleID)
	if err != nil {
		return web.ErrInternalServerError
	}

	response := createUserResponse{}
	response.User.ID = user.ID
	response.User.Name = user.Name
	response.User.CreatedAt = user.CreatedAt
	response.User.UpdatedAt = user.UpdatedAt
	return web.NewAPISuccessResponse(response, http.StatusOK)
}
