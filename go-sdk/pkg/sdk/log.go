package sdk

import (
	"log"
	"os"
)

// Logger interface used by the go sdk
type Logger interface {
	Debug(v ...interface{})
	Debugf(format string, v ...interface{})
	Info(v ...interface{})
	Infof(format string, v ...interface{})
	Warn(v ...interface{})
	Warnf(format string, v ...interface{})
	Error(v ...interface{})
	Errorf(format string, v ...interface{})
	Fatal(v ...interface{})
	Fatalf(format string, v ...interface{})
}

// DefaultLogger implementation of Logger using the go log package
type DefaultLogger struct {
	logger *log.Logger
}

// newDefaultLogger creates a new Default Logger
func newDefaultLogger() *DefaultLogger {
	return &DefaultLogger{logger: log.New(os.Stdout, "", 5)}
}

func (d DefaultLogger) Debug(v ...interface{}) {
	d.logger.Println(v...)
}

func (d DefaultLogger) Debugf(format string, v ...interface{}) {
	d.logger.Printf(format, v...)
}

func (d DefaultLogger) Info(v ...interface{}) {
	d.logger.Println(v...)
}

func (d DefaultLogger) Infof(format string, v ...interface{}) {
	d.logger.Printf(format, v...)
}

func (d DefaultLogger) Warn(v ...interface{}) {
	d.logger.Println(v...)
}

func (d DefaultLogger) Warnf(format string, v ...interface{}) {
	d.logger.Printf(format, v...)
}

func (d DefaultLogger) Error(v ...interface{}) {
	d.logger.Print(v...)
}

func (d DefaultLogger) Errorf(format string, v ...interface{}) {
	d.logger.Printf(format, v...)
}

func (d DefaultLogger) Fatal(v ...interface{}) {
	d.logger.Fatal(v...)
}

func (d DefaultLogger) Fatalf(format string, v ...interface{}) {
	d.logger.Fatalf(format, v...)
}
