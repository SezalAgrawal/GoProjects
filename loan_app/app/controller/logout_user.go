package controller

import (
	"net/http"

	"github.com/goProjects/loan_app/app/service"
	"github.com/goProjects/loan_app/lib/db"
	"github.com/goProjects/loan_app/lib/utils"
	"github.com/goProjects/loan_app/lib/web"
)

func LogoutUser(request *web.Request) *web.APIResponse {
	userID := request.Value(utils.CurrentUserIDStoreKey).(string)
	if err := service.NewUserService().DeleteAllAccessTokens(request.Context(), db.Get(), userID); err != nil {
		return web.ErrInternalServerError
	}

	return web.NewAPISuccessResponse(nil, http.StatusOK)
}
