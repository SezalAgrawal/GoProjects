package controller

import (
	"errors"
	"net/http"

	"github.com/goProjects/loan_app/app/service"
	"github.com/goProjects/loan_app/app/store"
	"github.com/goProjects/loan_app/lib/db"
	"github.com/goProjects/loan_app/lib/web"
)

type approveLoanResponse struct {
	Loan loan `json:"loan"`
}

func ApproveLoan(request *web.Request) *web.APIResponse {
	loanID := request.GetPathParam("id")
	l, err := service.NewLoanService().ApproveLoan(request.Context(), db.Get(), loanID)
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

	response := approveLoanResponse{
		Loan: convertLoanModelToResp(l),
	}
	return web.NewAPISuccessResponse(response, http.StatusOK)
}
