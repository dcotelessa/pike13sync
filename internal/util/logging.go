package util

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

// SetupLogging configures logging to both file and console
func SetupLogging() (*os.File, error) {
	// Get log directory from environment or use default
	logPath := os.Getenv("LOG_PATH")
	if logPath == "" {
		// Determine if running in Docker or local
		baseDir := "."
		if os.Getenv("DOCKER_ENV") == "true" {
			baseDir = "/app"
		}
		logPath = filepath.Join(baseDir, "logs", "pike13sync.log")
	}
	
	// Create logs directory if it doesn't exist
	logDir := filepath.Dir(logPath)
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, fmt.Errorf("error creating log directory: %v", err)
	}
	
	// Open log file
	logFile, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, fmt.Errorf("error opening log file: %v", err)
	}
	
	// Set up multi-writer to log to both file and console
	multiWriter := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(multiWriter)
	
	// Log startup message
	log.Println("Starting Pike13 to Google Calendar sync")
	
	return logFile, nil
}
