package logger

import (
	"log"
	"os"
)

var INFO *log.Logger
var DEBUG *log.Logger
var ERROR *log.Logger

var file *os.File

func InitLogger(logName string) {
	file, err := os.OpenFile(logName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		panic(err)
	}

	INFO = log.New(file, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	DEBUG = log.New(file, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile)
	ERROR = log.New(file, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}

func CloseLogger() {
	file.Close()
}
