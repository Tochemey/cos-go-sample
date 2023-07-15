package log

import (
	"context"
	"os"
	"strings"
	"sync"

	"github.com/tochemey/gopack/log"
	"github.com/tochemey/gopack/log/zapl"
)

// defines the log levels map
var levelsMap = map[string]log.Level{
	"info":  log.InfoLevel,
	"debug": log.DebugLevel,
	"warn":  log.WarningLevel,
	"error": log.ErrorLevel,
	"fatal": log.FatalLevel,
	"panic": log.PanicLevel,
}
var (
	// globalLevel defines the log level used by all loggers
	globalLevel log.Level
	once        sync.Once
	logger      *zapl.Log
)

// init sets the global log level and the logger
func init() {
	// let us load the configuration
	config := loadConfig()
	once.Do(func() {
		// set the log level to debug level
		globalLevel = log.DebugLevel
		// lookup for the corresponding log level
		if level, ok := levelsMap[strings.ToLower(config.LogLevel)]; ok {
			globalLevel = level
		}
		// logger sets the logger to use in the application
		logger = zapl.New(globalLevel, os.Stderr)
	})
}

// Info logs to INFO level.
func Info(v ...any) {
	logger.Info(v...)
}

// Infof logs to INFO level
func Infof(format string, v ...any) {
	logger.Infof(format, v...)
}

// Warn logs to the WARNING level.
func Warn(v ...any) {
	logger.Warn(v...)
}

// Warnf logs to the WARNING level.
func Warnf(format string, v ...any) {
	logger.Warnf(format, v...)
}

// Error logs to the ERROR level.
func Error(v ...any) {
	logger.Error(v...)
}

// Errorf logs to the ERROR level.
func Errorf(format string, v ...any) {
	logger.Errorf(format, v...)
}

// Fatal logs to the FATAL level followed by a call to os.Exit(1).
func Fatal(v ...any) {
	logger.Fatal(v...)
}

// Fatalf logs to the FATAL level followed by a call to os.Exit(1).
func Fatalf(format string, v ...any) {
	logger.Fatalf(format, v...)
}

// Panic logs to the PANIC level followed by a call to panic().
func Panic(v ...any) {
	logger.Panic(v...)
}

// Panicf logs to the PANIC level followed by a call to panic().
func Panicf(format string, v ...any) {
	logger.Panicf(format, v...)
}

// WithContext returns the Logger associated with the ctx.
// This will set the traceid, requestid and spanid in case there are
// in the context
func WithContext(ctx context.Context) log.Logger {
	return logger.WithContext(ctx)
}
