package logger

import (
	"io"
	"os"
	"sync"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

// Package logger provides a simple logging utility using logrus and lumberjack for log rotation.
// It initializes different loggers for various log levels and provides functions to log messages at those levels.
var (
	once          sync.Once
	RequestLogger *logrus.Logger
	InfoLogger    *logrus.Logger
	WarnLogger    *logrus.Logger
	ErrorLogger   *logrus.Logger
	FatalLogger   *logrus.Logger
	PanicLogger   *logrus.Logger
	TraceLogger   *logrus.Logger
	DebugLogger   *logrus.Logger

	REQUEST_LOG_FILE = "logs/request.log"
	INFO_LOG_FILE    = "logs/info.log"
	WARN_LOG_FILE    = "logs/warn.log"
	ERROR_LOG_FILE   = "logs/error.log"
	FATAL_LOG_FILE   = "logs/fatal.log"
	PANIC_LOG_FILE   = "logs/panic.log"
	TRACE_LOG_FILE   = "logs/trace.log"
	DEBUG_LOG_FILE   = "logs/debug.log"
)

func InitLoggers() {
	once.Do(func() {
		// Using TextFormatter for log formatting
		// This allows for more human-readable logs
		formatter := &logrus.TextFormatter{
			TimestampFormat: "2006-01-02 15:04:05",
			FullTimestamp:   true,
		}

		// Using lumberjack for log rotation
		// This allows for log files to be rotated based on size and age
		// MaxSize is the maximum size in megabytes before the log file is rotated
		// MaxBackups is the maximum number of old log files to keep
		// MaxAge is the maximum number of days to retain old log files
		// Compress indicates whether to compress the rotated log files
		requestFile := &lumberjack.Logger{
			Filename:   REQUEST_LOG_FILE,
			MaxSize:    100,
			MaxBackups: 7,
			MaxAge:     7,
			Compress:   true,
		}

		infoFile := &lumberjack.Logger{
			Filename:   INFO_LOG_FILE,
			MaxSize:    50,
			MaxBackups: 5,
			MaxAge:     14,
			Compress:   true,
		}

		warnFile := &lumberjack.Logger{
			Filename:   WARN_LOG_FILE,
			MaxSize:    20,
			MaxBackups: 10,
			MaxAge:     30,
			Compress:   true,
		}

		errorFile := &lumberjack.Logger{
			Filename:   ERROR_LOG_FILE,
			MaxSize:    20,
			MaxBackups: 15,
			MaxAge:     90,
			Compress:   true,
		}

		fatalFile := &lumberjack.Logger{
			Filename:   FATAL_LOG_FILE,
			MaxSize:    10,
			MaxBackups: 10,
			MaxAge:     180,
			Compress:   true,
		}

		panicFile := &lumberjack.Logger{
			Filename:   PANIC_LOG_FILE,
			MaxSize:    10,
			MaxBackups: 10,
			MaxAge:     180,
			Compress:   true,
		}

		traceFile := &lumberjack.Logger{
			Filename:   TRACE_LOG_FILE,
			MaxSize:    30,
			MaxBackups: 3,
			MaxAge:     3,
			Compress:   true,
		}

		debugFile := &lumberjack.Logger{
			Filename:   DEBUG_LOG_FILE,
			MaxSize:    30,
			MaxBackups: 5,
			MaxAge:     7,
			Compress:   true,
		}

		// Configure each logger with the specified format and output
		// The loggers will write to both the console (stdout) and the specified log files
		RequestLogger = logrus.New()
		RequestLogger.SetOutput(io.MultiWriter(os.Stdout, requestFile))
		RequestLogger.SetFormatter(formatter)
		RequestLogger.SetLevel(logrus.InfoLevel)

		InfoLogger = logrus.New()
		InfoLogger.SetOutput(io.MultiWriter(os.Stdout, infoFile))
		InfoLogger.SetFormatter(formatter)
		InfoLogger.SetLevel(logrus.InfoLevel)

		WarnLogger = logrus.New()
		WarnLogger.SetOutput(io.MultiWriter(os.Stdout, warnFile))
		WarnLogger.SetFormatter(formatter)
		WarnLogger.SetLevel(logrus.WarnLevel)

		ErrorLogger = logrus.New()
		ErrorLogger.SetOutput(io.MultiWriter(os.Stdout, errorFile))
		ErrorLogger.SetFormatter(formatter)
		ErrorLogger.SetLevel(logrus.ErrorLevel)

		FatalLogger = logrus.New()
		FatalLogger.SetOutput(io.MultiWriter(os.Stdout, fatalFile))
		FatalLogger.SetFormatter(formatter)
		FatalLogger.SetLevel(logrus.FatalLevel)

		PanicLogger = logrus.New()
		PanicLogger.SetOutput(io.MultiWriter(os.Stdout, panicFile))
		PanicLogger.SetFormatter(formatter)
		PanicLogger.SetLevel(logrus.PanicLevel)

		TraceLogger = logrus.New()
		TraceLogger.SetOutput(io.MultiWriter(os.Stdout, traceFile))
		TraceLogger.SetFormatter(formatter)
		TraceLogger.SetLevel(logrus.TraceLevel)

		DebugLogger = logrus.New()
		DebugLogger.SetOutput(io.MultiWriter(os.Stdout, debugFile))
		DebugLogger.SetFormatter(formatter)
		DebugLogger.SetLevel(logrus.DebugLevel)
	})
}

// GetLogger returns a singleton instance of logrus.Logger
func GetLogger(level logrus.Level) *logrus.Logger {
	if RequestLogger == nil || InfoLogger == nil ||
		WarnLogger == nil || ErrorLogger == nil ||
		FatalLogger == nil || PanicLogger == nil ||
		TraceLogger == nil || DebugLogger == nil {
		// Initialize the loggers if they are not already initialized
		// This ensures that the loggers are only initialized once
		InitLoggers()
	}

	// Set the log level for the logger
	switch level {
	case logrus.InfoLevel:
		return InfoLogger
	case logrus.WarnLevel:
		return WarnLogger
	case logrus.ErrorLevel:
		return ErrorLogger
	case logrus.FatalLevel:
		return FatalLogger
	case logrus.PanicLevel:
		return PanicLogger
	case logrus.TraceLevel:
		return TraceLogger
	default:
		return DebugLogger
	}
}

// Log functions for different log levels
func Info(msg string, fields ...logrus.Fields) {
	logger := GetLogger(logrus.InfoLevel)
	if len(fields) > 0 {
		logger.WithFields(fields[0]).Info(msg)
	} else {
		logger.Info(msg)
	}
}

func Warn(msg string, fields ...logrus.Fields) {
	logger := GetLogger(logrus.WarnLevel)
	if len(fields) > 0 {
		logger.WithFields(fields[0]).Warn(msg)
	} else {
		logger.Warn(msg)
	}
}

func Error(msg string, fields ...logrus.Fields) {
	logger := GetLogger(logrus.ErrorLevel)
	if len(fields) > 0 {
		logger.WithFields(fields[0]).Error(msg)
	} else {
		logger.Error(msg)
	}
}

func Fatal(msg string, fields ...logrus.Fields) {
	logger := GetLogger(logrus.FatalLevel)
	if len(fields) > 0 {
		logger.WithFields(fields[0]).Fatal(msg)
	} else {
		logger.Fatal(msg)
	}
}

func Panic(msg string, fields ...logrus.Fields) {
	logger := GetLogger(logrus.PanicLevel)
	if len(fields) > 0 {
		logger.WithFields(fields[0]).Panic(msg)
	} else {
		logger.Panic(msg)
	}
}

func Trace(msg string, fields ...logrus.Fields) {
	logger := GetLogger(logrus.TraceLevel)
	if len(fields) > 0 {
		logger.WithFields(fields[0]).Trace(msg)
	} else {
		logger.Trace(msg)
	}
}

func Debug(msg string, fields ...logrus.Fields) {
	logger := GetLogger(logrus.DebugLevel)
	if len(fields) > 0 {
		logger.WithFields(fields[0]).Debug(msg)
	} else {
		logger.Debug(msg)
	}
}
