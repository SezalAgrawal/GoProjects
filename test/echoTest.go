package main

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"net/http"
)

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Logger.SetLevel(log.INFO)

	e.GET("/test", func(c echo.Context) error {
		return c.String(http.StatusOK, "GET /test")
	}).Name = "get-test"
	e.GET("/test/:test", func(c echo.Context) error {
		fmt.Printf("%#v\n", c.ParamValues())
		fmt.Printf("%#v\n", c.ParamNames()) // return []string{"test"}
		return c.String(http.StatusOK, "GET /test/:test="+c.Param("test"))
	}).Name = "get-list-test"
	e.PUT("/test/:id", func(c echo.Context) error {
		fmt.Printf("%#v\n", c.ParamValues())
		fmt.Printf("%#v\n", c.ParamNames()) // return []string{"test"}
		return c.String(http.StatusOK, "PUT /test/:id="+c.Param("test"))
	}).Name = "put-test"

	e.Logger.Fatal(e.Start(":1323"))
}
