package main

import (
	"fmt"

	"github.com/pkg/errors"
)

var (
	err1 = errors.New("error message1")
	err2 = errors.New("error message2")
)

func handle(err error) {
	switch errors.Cause(err) {
	case err1:
		fmt.Printf("err1: %v\n", err)
	case err2:
		fmt.Printf("err2: %v\n", err)
	default:
		fmt.Println(err)
	}

	if err != nil {
		// Trace top
		type stackTracer interface {
			StackTrace() errors.StackTrace
		}
		err, ok := errors.Cause(err).(stackTracer)
		if !ok {
			panic("oops, err does not implement stackTracer")
		}
		st := err.StackTrace()
		fmt.Printf("%+v\n", st[0])

		// SimpleTrace
		//fmt.Printf("%+v\n", err)
	}
}

func main() {
	var err error

	err = err1
	err = errors.Wrap(err, "wrapped")
	handle(err)

	fmt.Println()

	err = err2
	err = errors.Wrap(err, "wrapped2")
	handle(err)
}
