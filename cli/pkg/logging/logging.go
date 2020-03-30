package logging

import (
	"io"
	"log"
)

var (
	// Info provides a logger for info messages
	Info *log.Logger
	// Warning provides a logger for warnings
	Warning *log.Logger
	// Error provides a logger for error
	Error *log.Logger
)

// InitLoggers initializes the loggers
func InitLoggers(
	infoHandle io.Writer,
	warningHandle io.Writer,
	errorHandle io.Writer) {

	Info = log.New(infoHandle,
		"INFO: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Warning = log.New(warningHandle,
		"WARNING: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Error = log.New(errorHandle,
		"ERROR: ",
		log.Ldate|log.Ltime|log.Lshortfile)
}
