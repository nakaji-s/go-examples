package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/jinzhu/configor"
)

var Config = struct {
	APPName string `default:"app name"`

	DB struct {
		Name     string
		User     string `default:"root"`
		Password string `required:"true" env:"DBPassword"`
		Port     uint   `default:"3306"`
		Abc      string `required:"true"`
	}

	Contacts []struct {
		Name  string
		Email string `required:"true"`
	}
}{}

func main() {
	os.Setenv("DBPassword", "envDBPass")
	if err := configor.Load(&Config, "config.yml"); err != nil {
		fmt.Println(err)
		fmt.Println()
	}

	if b, err := json.MarshalIndent(&Config, "", "    "); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(string(b))
	}

	//fmt.Printf("config: %#v", Config)
}
