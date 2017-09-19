package main

import (
	"net/http"

	"io/ioutil"

	"github.com/labstack/echo"
)

func main() {
	// start server
	e := echo.New()

	e.GET("/hello", func(c echo.Context) error {
		return c.String(http.StatusOK, "hello")
	})

	e.POST("/echo", func(c echo.Context) error {
		body, err := ioutil.ReadAll(c.Request().Body)
		if err != nil {
			return err
		}
		return c.String(http.StatusOK, string(body))
	})

	e.Logger.Fatal(e.Start("127.0.0.1:8080"))
}
