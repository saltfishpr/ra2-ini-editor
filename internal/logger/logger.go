package logger

import (
	"os"

	"github.com/charmbracelet/log"
	"github.com/wailsapp/wails/v2/pkg/logger"
)

type Logger struct{}

var _ logger.Logger = (*Logger)(nil)

// Print implements [logger.Logger].
func (l *Logger) Print(message string) {
	log.Print("PRINT | " + message)
}

// Trace implements [logger.Logger].
func (l *Logger) Trace(message string) {
	log.Print("TRACE | " + message)
}

// Debug implements [logger.Logger].
func (l *Logger) Debug(message string) {
	log.Print("DEBUG | " + message)
}

// Info implements [logger.Logger].
func (l *Logger) Info(message string) {
	log.Print("INFO  | " + message)
}

// Warning implements [logger.Logger].
func (l *Logger) Warning(message string) {
	log.Print("WARN  | " + message)
}

// Error implements [logger.Logger].
func (l *Logger) Error(message string) {
	log.Print("ERROR | " + message)
}

// Fatal implements [logger.Logger].
func (l *Logger) Fatal(message string) {
	log.Print("FATAL | " + message)
	os.Exit(1)
}
