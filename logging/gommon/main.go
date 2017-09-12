package main

import (
	"os"
	"time"

	"github.com/labstack/gommon/log"
	"github.com/lestrrat/go-file-rotatelogs"
)

func main() {
	// settings
	log.SetLevel(log.ERROR)
	log.SetOutput(os.Stdout)
	log.SetPrefix("-")

	// msg1
	log.Error("example error message1")

	// msg2(set header)
	log.SetHeader("--------")
	log.Error("example error message2")

	// msg3(file with rotate)
	logf, err := rotatelogs.New(
		"./log.%Y%m%d%H%M",
		rotatelogs.WithLinkName("./log"),
		rotatelogs.WithMaxAge(24*time.Hour),
		rotatelogs.WithRotationTime(time.Hour),
	)
	if err != nil {
		log.Printf("failed to create rotatelogs: %s", err)
		return
	}
	log.SetOutput(logf)
	log.Error("example error message3")
}
