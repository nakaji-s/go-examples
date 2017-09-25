package main

import (
	"bytes"
	"encoding/json"
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
	// Retrieve Paramteters from request
	queryParams := c.QueryParams()
	b, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		return err
	}
	c.Request().Body = ioutil.NopCloser(bytes.NewBuffer(b)) // Reset
	matchedPath := c.Path()                                 // TODO: convert :id to {id}
	method := c.Request().Method

	// create swagger path object
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

	// validate each paramteres
	ret := validate.Result{}
	for _, param := range operation.OperationProps.Parameters {
		switch param.In {
		case "query":
			validator := validate.NewParamValidator(&param, strfmt.Default)

			// TODO: need to convert value from string to expected before validate
			ret.Merge(validator.Validate(strings.Join(queryParams[param.Name], ",")))
		case "body":
			// TODO: Get schema in a correct way
			var schema *spec.Schema
			refURL := param.Schema.Ref.Ref.GetURL()
			if refURL == nil {
				schema = param.Schema
			} else {
				tmp := v.swagger.Definitions[strings.TrimPrefix(refURL.Fragment, "/definitions/")]
				schema = &tmp
			}
			validator := validate.NewSchemaValidator(schema, nil, "", strfmt.Default)

			switch schema.Type[0] {
			case "object":
				var m interface{}
				if err := json.Unmarshal(b, &m); err != nil {
					return err
				}

				ret.Merge(validator.Validate(m))
			case "string":
				// TODO
				//body := string(b)
			}

		case "path":
			validator := validate.NewParamValidator(&param, strfmt.Default)

			// TODO: need to convert value from string to expected before validate
			ret.Merge(validator.Validate(c.Param(param.Name)))
		}
	}

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
