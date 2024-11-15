package utils

import (
	"log"
	"os"
)

type Logger struct {
	Info  *log.Logger
	Error *log.Logger
}

// NewLogger initializes and returns a new Logger instance.
func NewLogger() *Logger {
	return &Logger{
		Info:  log.New(os.Stdout, "INFO: ", log.LstdFlags),
		Error: log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile),
	}
}
