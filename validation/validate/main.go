package main

import (
	"io/ioutil"

	"fmt"

	"github.com/go-openapi/spec"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/validate"
)

func newSwaggerFromJson(filename string) (spec.Swagger, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return spec.Swagger{}, err
	}

	swagger := spec.Swagger{}
	swagger.UnmarshalJSON(data)
	if err != nil {
		return spec.Swagger{}, err
	}

	return swagger, nil
}

func main() {
	swagger, err := newSwaggerFromJson("petstore.json")
	if err != nil {
		panic(err)
	}

	param := &swagger.Paths.Paths["/pets"].Get.OperationProps.Parameters[0]
	//pretty.Println(param)
	validator := validate.NewParamValidator(param, strfmt.Default)
	ret := validator.Validate(33)
	if ret != nil {
		fmt.Println(ret.AsError())
	}

	ret = validator.Validate("aaa")
	if ret != nil {
		fmt.Println(ret.AsError())
	}
}
