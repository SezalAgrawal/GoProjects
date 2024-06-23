package app

import (
	"fmt"
	"net/http"

	"github.com/goProjects/loan_app/app/controller"
	"github.com/julienschmidt/httprouter"
)

func InitRoutes(r *httprouter.Router) {
	r.GET("/ping", func(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
		writer.Header().Add("Content-Type", "application/json")
		_, _ = fmt.Fprint(writer, "{ \"message\":\"pong\"}")
	})

	// user
	r.POST("/api/v1/users", ServeEndpoint(controller.CreateUser))
	r.POST("/api/v1/users/login", ServeEndpoint(controller.LoginUser))
	r.DELETE("/api/v1/users/logout", ServeEndpoint(validateAccessTokenMiddleware([]role{}, controller.LogoutUser)))

	// loan
	r.POST("/api/v1/loans", ServeEndpoint(validateAccessTokenMiddleware([]role{}, controller.CreateLoan)))
	r.GET("/api/v1/loans", ServeEndpoint(validateAccessTokenMiddleware([]role{}, controller.GetAllLoans)))
	r.POST("/api/v1/loans/:id/approve", ServeEndpoint(validateAccessTokenMiddleware([]role{adminRole}, controller.ApproveLoan)))
	r.POST("/api/v1/loans/:id/payment", ServeEndpoint(validateAccessTokenMiddleware([]role{}, controller.PayLoan)))
}
