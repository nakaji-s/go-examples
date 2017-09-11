package main

import (
	"github.com/labstack/gommon/log"
	"os"
)

func main() {
	// settings
	log.SetLevel(log.ERROR)
	log.SetOutput(os.Stdout)
	log.SetPrefix("-")

	// msg1
	log.Error("example error message1")

	// msg2
	log.SetHeader("--------")
	log.Error("example error message2")
}