package errorHandler

import (
	"fmt"
	"github.com/maleck13/photoProcessor/logger"
)

func FailOnError(err error, msg string) {
	if err != nil {
		logger.ErrorLog.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}

func LogOnError(err error, msg string) {
	if err != nil {
		logger.ErrorLog.Println(err, msg)
	}
}
