package logger

import (
	"log"
	"os"
	"sync"
)

var (
	infoLogger  = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	errorLogger = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
	warnLogger  = log.New(os.Stdout, "WARN: ", log.Ldate|log.Ltime|log.Lshortfile)
	debugLogger = log.New(os.Stdout, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile)
	logLevel    = "info" // Default log level
	mu          sync.Mutex
)

func InitLogger() {
	// You could read log level from config here
}

func SetLogLevel(level string) {
	mu.Lock()
	defer mu.Unlock()
	logLevel = level
}

func Info(format string, v ...interface{}) {
	mu.Lock()
	defer mu.Unlock()
	infoLogger.Printf(format, v...)
}

func Error(format string, v ...interface{}) {
	mu.Lock()
	defer mu.Unlock()
	errorLogger.Printf(format, v...)
}

func Warn(format string, v ...interface{}) {
	mu.Lock()
	defer mu.Unlock()
	warnLogger.Printf(format, v...)
}

func Debug(format string, v ...interface{}) {
	mu.Lock()
	defer mu.Unlock()
	if logLevel == "debug" {
		debugLogger.Printf(format, v...)
	}
}
