package controller

import (
	"errors"
	"net/http"

	"github.com/goProjects/loan_app/app/service"
	"github.com/goProjects/loan_app/app/store"
	"github.com/goProjects/loan_app/lib/db"
	"github.com/goProjects/loan_app/lib/utils"
	"github.com/goProjects/loan_app/lib/web"
)

type payLoanRequestParams struct {
	OtsID         string `json:"one_time_settlement_id" validate:"notblank"`
	AmountInCents int64  `json:"amount_in_cents" validate:"required,gt=0"`
}

type payLoanResponse struct {
	Loan loan `json:"loan"`
}

func PayLoan(request *web.Request) *web.APIResponse {
	userID := request.Value(utils.CurrentUserIDStoreKey).(string)

	data := &payLoanRequestParams{}

	if err := request.ValidateBodyToStruct(data); err != nil {
		return web.ErrBadRequest(err.Error())
	}

	loanID := request.GetPathParam("id")
	l, err := service.NewLoanService().PayLoan(request.Context(), db.Get(), userID, loanID, data.OtsID, data.AmountInCents)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			return web.ErrNotFound(err.Error())
		case errors.Is(err, service.ErrInvalidLoan):
			return web.ErrBadRequest(err.Error())
		default:
			return web.ErrInternalServerError
		}
	}

	response := payLoanResponse{
		Loan: convertLoanModelToResp(l),
	}
	return web.NewAPISuccessResponse(response, http.StatusOK)
}
