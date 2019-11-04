package logging

import (
	"fmt"
	"strings"
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

// PrintLogStringLevel prints the log according to the log level that is set in the flags
func PrintLogStringLevel(message string, printInLevel string) {

	lev := GetLogLevel(printInLevel)
	if lev < 0 {
		PrintLog("Received unknown log level: "+printInLevel, InfoLevel)
	} else {
		PrintLog(message, lev)
	}
}

// GetLogLevel parses a string and returns the appropriate LogLevelType
func GetLogLevel(logLevel string) LogLevelType {

	if strings.ToLower(logLevel) == "info" {
		return InfoLevel
	} else if strings.ToLower(logLevel) == "debug" ||
		strings.ToLower(logLevel) == "verbose" ||
		strings.ToLower(logLevel) == "v" {
		return VerboseLevel
	} else if strings.ToLower(logLevel) == "error" ||
		strings.ToLower(logLevel) == "quiet" ||
		strings.ToLower(logLevel) == "q" {
		return QuietLevel
	}
	return -1
}
