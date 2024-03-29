package logging

import (
	"fmt"
)

// LogLevel specifies the used logging level
var LogLevel LogLevelType

// LogLevelType represents a type for the log levels
type LogLevelType int

const (
	// VerboseLevel logs debug, info, and error log messages
	VerboseLevel LogLevelType = iota
	// InfoLevel logs info and error log messages
	InfoLevel
	// QuietLevel logs error log messages
	QuietLevel
)

// PrintLog prints the log according to the log level that is set in the flags
func PrintLog(message string, printInLevel LogLevelType) {

	if LogLevel <= printInLevel {
		fmt.Println(message)
	}
}
