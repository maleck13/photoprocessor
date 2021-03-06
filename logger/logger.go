package logger

import (
	"fmt"
	"log"
	"os"
)

var (
	TraceLog   *log.Logger
	InfoLog    *log.Logger
	WarningLog *log.Logger
	ErrorLog   *log.Logger
)

const (
	LOG_FILE  = "/var/log/photoprocessor/log.log"
	ERROR_LOG = "/var/log/photoprocessor/error.log"
)

func InitLogger() {

	f, e := os.OpenFile(LOG_FILE, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	ef, e := os.OpenFile(ERROR_LOG, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	if e != nil {
		fmt.Printf("error opening log files " + e.Error())
		panic(e)
	}

	TraceLog = log.New(f,
		"TRACE: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	InfoLog = log.New(f,
		"INFO: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	WarningLog = log.New(f,
		"WARNING: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	ErrorLog = log.New(ef,
		"ERROR: ",
		log.Ldate|log.Ltime|log.Lshortfile)

}
