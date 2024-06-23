package controller

import (
	"net/http"

	"github.com/goProjects/loan_app/app/service"
	"github.com/goProjects/loan_app/lib/db"
	"github.com/goProjects/loan_app/lib/utils"
	"github.com/goProjects/loan_app/lib/web"
)

type createLoanRequestParams struct {
	AmountInCents   int64 `json:"amount_in_cents" validate:"required,gt=0"`
	Term            int64 `json:"term" validate:"required,gt=0"`
	FrequencyInDays int64 `json:"frequency_in_days" validate:"required,gt=0"`
}

type createLoanResponse struct {
	Loan loan `json:"loan"`
}

func CreateLoan(request *web.Request) *web.APIResponse {
	userID := request.Value(utils.CurrentUserIDStoreKey).(string)

	data := &createLoanRequestParams{}

	if err := request.ValidateBodyToStruct(data); err != nil {
		return web.ErrBadRequest(err.Error())
	}

	l, err := service.NewLoanService().CreateLoan(request.Context(), db.Get(), userID, &service.CreateLoanParams{
		AmountInCents:   data.AmountInCents,
		Term:            data.Term,
		FrequencyInDays: data.FrequencyInDays,
	})
	if err != nil {
		return web.ErrInternalServerError
	}

	response := createLoanResponse{
		Loan: convertLoanModelToResp(l),
	}
	return web.NewAPISuccessResponse(response, http.StatusOK)
}
