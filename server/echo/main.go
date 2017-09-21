package main

import (
	"bytes"
	"net/http"

	"io/ioutil"

	"fmt"

	"github.com/labstack/echo"
)

func main() {
	// start server
	e := echo.New()

	// original middleware for dump queryParam,header,body
	e.Use(
		func(next echo.HandlerFunc) echo.HandlerFunc {
			return func(c echo.Context) error {
				fmt.Println("queryParam:", c.QueryParams())
				fmt.Println("header:", c.Request().Header)
				body, err := ioutil.ReadAll(c.Request().Body)
				if err != nil {
					return err
				}
				c.Request().Body = ioutil.NopCloser(bytes.NewBuffer(body)) // Reset
				fmt.Println("body", string(body))

				err = next(c)
				return err
			}
		},
	)

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

	e.Logger.Fatal(e.Start("127.0.0.1:8081"))
}
