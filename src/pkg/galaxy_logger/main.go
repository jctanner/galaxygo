package galaxy_logger

import (
	"fmt"
	"time"
)

type Logger struct {
	level string
}

func (l *Logger) Info(message string) {
	l.printLog("INFO", message)
}

func (l *Logger) Error(message string) {
	l.printLog("ERROR", message)
}

func (l *Logger) Debug(message string) {
	l.printLog("DEBUG", message)
}

func (l *Logger) printLog(level string, message string) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	logMessage := fmt.Sprintf("[%s] %s - %s", timestamp, level, message)
	fmt.Println(logMessage)
}
