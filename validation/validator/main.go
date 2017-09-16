package main

import (
	"fmt"

	"gopkg.in/go-playground/validator.v9"
)

type User struct {
	FirstName      string     `validate:"required"`
	LastName       string     `validate:"required"`
	Age            uint8      `validate:"gte=0,lte=130"`
	Email          string     `validate:"required,email"`
	FavouriteColor string     `validate:"iscolor"`                // alias for 'hexcolor|rgb|rgba|hsl|hsla'
	Addresses      []*Address `validate:"required,dive,required"` // a person can have a home and cottage...
}

type Address struct {
	Street string `validate:"required"`
	City   string `validate:"required"`
	Planet string `validate:"required"`
	Phone  string `validate:"required,mytag"`
}

// use a single instance of Validate, it caches struct info
var validate = validator.New()

func myValidation(fl validator.FieldLevel) bool {
	fmt.Println("value = ", fl.Field())
	return false
}

func main() {
	// set custom validation tag
	validate.RegisterValidation("mytag", myValidation)

	// validate struct
	validateStruct()
	fmt.Println()

	// validate var
	validateVariable()
}

func validateStruct() {
	user := &User{
		FirstName:      "Badger",
		LastName:       "Smith",
		Age:            135,
		Email:          "Badger.Smith@gmail.com",
		FavouriteColor: "#000-",
		Addresses: []*Address{&Address{
			Street: "Eavesdown Docks",
			Planet: "Persphone",
			Phone:  "none",
		}},
	}

	if err := validate.Struct(user); err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			fmt.Println(err)
		}
	}

	// validate partial struct
	fmt.Println()
	if err := validate.StructPartial(user, "Age"); err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			fmt.Println(err)
		}
	}
}

func validateVariable() {
	if err := validate.Var("joeybloggs.gmail.com", "required,email"); err != nil {
		// output: Key: "" Error:Field validation for "" failed on the "email" tag
		fmt.Println(err)
	}

	if err := validate.Var("aaa@bbb.ccc", "required,email"); err != nil {
		fmt.Println(err)
	}
}
