package main

import (
	"io/ioutil"
	"net/http"

	"fmt"

	"regexp"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/spec"
	"github.com/go-openapi/strfmt"
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

func ReadFile(filename string) []byte {
	data, _ := ioutil.ReadFile(filename)
	return data
}

func (v RequestValidator) Validate(c echo.Context) error {
	// retrieve paramteters from request
	const sentinel = "/"
	re := regexp.MustCompile(`:(.+?)/`)
	matchedPathWithSentinel := re.ReplaceAllString(c.Path()+sentinel, `{$1}/`)
	matchedPath := matchedPathWithSentinel[:len(matchedPathWithSentinel)-1]

	// create swagger path object
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

	// create requestValidator
	m := map[string]spec.Parameter{}
	for i, param := range operation.OperationProps.Parameters {
		m[fmt.Sprint(i)] = param
	}
	binder := middleware.NewUntypedRequestBinder(m, v.swagger, strfmt.Default)

	// get PathParams from request and set for validate
	pathParams := middleware.RouteParams{}
	for _, paramName := range c.ParamNames() {
		pathParams = append(pathParams, middleware.RouteParam{paramName, c.Param(paramName)})
	}

	// set Params if default value is needed
	data := map[string]interface{}{}

	// validate request and set defalut value to data
	err := binder.Bind(c.Request(), pathParams, runtime.JSONConsumer(), &data)
	if err != nil {
		return err
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

	e.GET("/pets/:id", func(c echo.Context) error {
		return c.String(http.StatusOK, "pets "+c.Param("id"))
	})

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
