package main

import (
	"os"
	"time"

	"log/syslog"

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

	// msg4(to syslog)
	log.SetHeader(`{"time":"${time_rfc3339_nano}","level":"${level}","prefix":"${prefix}",` +
		`"file":"${short_file}","line":"${line}"}`)
	syslog, err := syslog.New(syslog.LOG_NOTICE|syslog.LOG_USER, "syslog-example")
	if err != nil {
		panic(err)
	}
	log.SetOutput(syslog)
	log.Error("exmaple error message4")
}
