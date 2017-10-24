package main

import (
	"bytes"
	"net/http"

	"io/ioutil"

	"fmt"

	"math/rand"

	"net/url"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/middleware"
)

func main() {
	// start server
	e := echo.New()

	e.Use(session.Middleware(sessions.NewCookieStore([]byte("secret"))))

	// original middleware for add header
	e.Use(
		func(next echo.HandlerFunc) echo.HandlerFunc {
			return func(c echo.Context) error {
				c.Request().Header.Set("dummy", "dummy")
				return next(c)
			}
		},
	)
	// reverse proxy
	url1, err := url.Parse("http://127.0.0.1:8082")
	if err != nil {
		e.Logger.Fatal(err)
	}
	e.Group("/blog", middleware.Proxy(
		&middleware.RoundRobinBalancer{
			Targets: []*middleware.ProxyTarget{
				&middleware.ProxyTarget{
					URL: url1,
				},
			},
		},
	))

	// second server
	go func() {
		e := echo.New()
		e.POST("/blog/echo", func(c echo.Context) error {
			body, err := ioutil.ReadAll(c.Request().Body)
			if err != nil {
				return err
			}
			return c.String(http.StatusOK, string(body)+c.Request().Header.Get("dummy"))
		})
		e.Start("127.0.0.1:8082")
	}()

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
				fmt.Println("body:", string(body))

				fmt.Println("path:", c.Request().RequestURI)
				fmt.Println("matchedPath:", c.Path())
				fmt.Println("method:", c.Request().Method)

				err = next(c)
				return err
			}
		},
	)

	e.Static("/static", "static")
	e.File("/", "static/index.html")

	e.GET("/hello", func(c echo.Context) error {
		sess, _ := session.Get("session", c)
		sess.Options = &sessions.Options{
			Path:     "/",
			MaxAge:   10,
			HttpOnly: true,
		}
		if _, ok := sess.Values["foo"]; ok == false {
			sess.Values["foo"] = rand.Int()
		}
		sess.Save(c.Request(), c.Response())

		return c.String(http.StatusOK, fmt.Sprint("hello", sess.Values["foo"]))
	})

	e.GET("/hello/:id", func(c echo.Context) error {
		return c.String(http.StatusOK, "hello "+c.Param("id"))
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
