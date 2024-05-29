package log

import (
	"errors"
	"fmt"
	"log"
	"os"
)

type LogLevel int

const (
	DebugLevel LogLevel = iota
	InfoLevel
	WarnLevel
	ErrorLevel
)

var (
	currentLogLevel = ErrorLevel
	logger          = log.New(os.Stdout, "", 0)
)

type AikidoFormatter struct{}

func (f *AikidoFormatter) Format(level LogLevel, message string) string {
	var levelStr string
	switch level {
	case DebugLevel:
		levelStr = "DEBUG"
	case InfoLevel:
		levelStr = "INFO"
	case WarnLevel:
		levelStr = "WARN"
	case ErrorLevel:
		levelStr = "ERROR"
	default:
		return "invalid log level"
	}

	logMessage := fmt.Sprintf("[AIKIDO][%s][GO] %s\n", levelStr, message)
	return logMessage
}

func logMessage(level LogLevel, args ...interface{}) {
	if level >= currentLogLevel {
		formatter := &AikidoFormatter{}
		message := fmt.Sprint(args...)
		formattedMessage := formatter.Format(level, message)
		logger.Print(formattedMessage)
	}
}

func Debug(args ...interface{}) {
	logMessage(DebugLevel, args...)
}

func Info(args ...interface{}) {
	logMessage(InfoLevel, args...)
}

func Warn(args ...interface{}) {
	logMessage(WarnLevel, args...)
}

func Error(args ...interface{}) {
	logMessage(ErrorLevel, args...)
}

func SetLogLevel(level string) error {
	switch level {
	case "DEBUG":
		currentLogLevel = DebugLevel
	case "INFO":
		currentLogLevel = InfoLevel
	case "WARN":
		currentLogLevel = WarnLevel
	case "ERROR":
		currentLogLevel = ErrorLevel
	default:
		return errors.New("invalid log level")
	}
	return nil
}

func init() {
	SetLogLevel("ERROR")
}
