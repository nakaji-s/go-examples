package main

import (
	"log/syslog"
	"os"

	"fmt"

	log "github.com/sirupsen/logrus"
	lSyslog "github.com/sirupsen/logrus/hooks/syslog"
)

type myHook struct{}

func (hook *myHook) Fire(entry *log.Entry) error {
	fmt.Println("message = " + entry.Message)
	return nil
}

func (hook *myHook) Levels() []log.Level {
	return log.AllLevels
}

func main() {
	// settings
	log.SetLevel(log.ErrorLevel)
	log.SetOutput(os.Stdout)

	// msg1
	log.Error("example error message1")

	// msg2 with JSON form
	log.SetFormatter(&log.JSONFormatter{})
	log.Error("example error message2")

	// msg3 to stdout and syslog
	hook, err := lSyslog.NewSyslogHook("", "", syslog.LOG_INFO, "")
	if err == nil {
		log.AddHook(hook)
	}

	// msg4 custom hook
	log.AddHook(&myHook{})
	log.Error("example error message4")
}
