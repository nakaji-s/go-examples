package main

import (
	"io/ioutil"
	"net/http"

	"fmt"

	"strings"

	"github.com/go-openapi/spec"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/validate"
	"github.com/labstack/echo"
)

type RequestValidator struct {
	swagger *spec.Swagger
}

func NewRequestValidator(filename string) (RequestValidator, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return RequestValidator{}, err
	}

	swagger := spec.Swagger{}
	swagger.UnmarshalJSON(data)
	if err != nil {
		return RequestValidator{}, err
	}

	return RequestValidator{&swagger}, nil
}

func (v RequestValidator) Validate(c echo.Context) error {
	queryParams := c.QueryParams()
	//headers := c.Request().Header
	//b, err := ioutil.ReadAll(c.Request().Body)
	//if err != nil {
	//	return err
	//}
	//c.Request().Body = ioutil.NopCloser(bytes.NewBuffer(body)) // Reset
	//body := string(body)
	//requestPath := c.Request().RequestURI
	matchedPath := c.Path()
	method := c.Request().Method

	path := v.swagger.Paths.Paths[matchedPath]
	var operation *spec.Operation
	switch method {
	case "GET":
		operation = path.Get
	case "PUT":
		operation = path.Put
	case "POST":
		operation = path.Post
	case "DELETE":
		operation = path.Delete
	}
	if operation == nil {
		return c.NoContent(http.StatusNotFound)
	}

	ret := validate.Result{}
	for _, param := range operation.OperationProps.Parameters {
		switch param.In {
		case "query":
			fmt.Println(param.Name, strings.Join(queryParams[param.Name], ","))
			validator := validate.NewParamValidator(&param, strfmt.Default)
			ret.Merge(validator.Validate(strings.Join(queryParams[param.Name], ",")))
		case "body":
		case "header":
		case "path":
		}
	}

	//param := &v.swagger.Paths.Paths["/pets"].Get.OperationProps.Parameters[0]
	//validator := validate.NewParamValidator(param, strfmt.Default)
	//
	//
	//ret.Merge(validator.Validate(33))
	//ret.Merge(validator.Validate("aaa"))
	//
	if ret.HasErrors() {
		return ret.AsError()
	}

	return nil
}

func main() {
	requestValidator, err := NewRequestValidator("petstore-expanded.json")
	if err != nil {
		panic(err)
	}

	// start server
	e := echo.New()

	// original middleware
	e.Use(
		func(next echo.HandlerFunc) echo.HandlerFunc {
			return func(c echo.Context) error {
				if err := requestValidator.Validate(c); err != nil {
					fmt.Println(err)
					return err
				}

				return next(c)
			}
		},
	)

	e.GET("/hello", func(c echo.Context) error {
		return c.String(http.StatusOK, "hello")
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
