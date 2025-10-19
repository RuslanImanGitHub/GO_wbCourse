package logger

import (
	"io"
	"log"
	"os"
)

var (
	Info  *log.Logger
	Error *log.Logger
)

func Init(infoHandler io.Writer, errorHandler io.Writer) {
	Info = log.New(infoHandler, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	Error = log.New(errorHandler, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}

func InitDefault() {
	Init(os.Stdout, os.Stderr)
}
