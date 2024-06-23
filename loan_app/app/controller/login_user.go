package controller

import (
	"errors"
	"net/http"

	"github.com/goProjects/loan_app/app/service"
	"github.com/goProjects/loan_app/app/store"
	"github.com/goProjects/loan_app/lib/db"
	"github.com/goProjects/loan_app/lib/web"
)

type loginUserRequestParams struct {
	Name     string `json:"name" validate:"notblank"`
	Password string `json:"password" validate:"notblank"`
}

type loginUserResponse struct {
	User struct {
		ID          string `json:"id"`
		AccessToken string `json:"access_token"`
	} `json:"user"`
}

func LoginUser(request *web.Request) *web.APIResponse {
	data := &loginUserRequestParams{}

	if err := request.ValidateBodyToStruct(data); err != nil {
		return web.ErrBadRequest(err.Error())
	}

	user, accessToken, err := service.NewUserService().GenerateAccessToken(request.Context(), db.Get(), data.Name, data.Password)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			return web.ErrNotFound(err.Error())
		default:
			return web.ErrInternalServerError
		}
	}

	response := loginUserResponse{}
	response.User.ID = user.ID
	response.User.AccessToken = accessToken
	return web.NewAPISuccessResponse(response, http.StatusOK)
}
