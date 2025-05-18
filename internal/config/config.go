package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"runtime"
)

// Config holds application configuration
type Config struct {
	Pike13URL             string `json:"pike13_url"`
	CalendarID            string `json:"calendar_id"`
	TimeZone              string `json:"time_zone"`
	CredentialsPath       string `json:"credentials_path"`
	LogPath               string `json:"log_path"`
	DryRun                bool   `json:"dry_run"`
	BaseDir               string `json:"-"` // Not serialized
}

// LoadConfig loads configuration from file and environment variables
func LoadConfig(configPath string) (*Config, error) {
	config := &Config{
		TimeZone: "America/Los_Angeles",
		DryRun:   false,
	}
	
	// Determine base directory
	config.BaseDir = determineBaseDir()
	
	// Set derived paths
	configDir := filepath.Join(config.BaseDir, "config")
	credentialsDir := filepath.Join(config.BaseDir, "credentials")
	logsDir := filepath.Join(config.BaseDir, "logs")
	
	// Set default values with correct paths
	config.Pike13URL = "https://herosjourneyfitness.pike13.com/api/v2/front/event_occurrences.json"
	
	// Calendar ID from environment variable or default
	config.CalendarID = os.Getenv("CALENDAR_ID")
	if config.CalendarID == "" {
		config.CalendarID = "primary"
	}
	
	// Credentials path from environment variable or default
	config.CredentialsPath = os.Getenv("GOOGLE_CREDENTIALS_FILE")
	if config.CredentialsPath == "" {
		config.CredentialsPath = filepath.Join(credentialsDir, "credentials.json")
	}
	
	// Log path from environment variable or default
	config.LogPath = os.Getenv("LOG_PATH")
	if config.LogPath == "" {
		config.LogPath = filepath.Join(logsDir, "pike13sync.log")
	}
	
	// Create directories if they don't exist
	os.MkdirAll(configDir, 0755)
	os.MkdirAll(credentialsDir, 0755)
	os.MkdirAll(logsDir, 0755)
	
	// First, try to load from environment variables (highest priority)
	loadConfigFromEnv(config)
	
	// Then, try to load from config file (only if values not already set from env)
	// Use provided config path or default
	confPath := configPath
	if confPath == "" {
		confPath = filepath.Join(configDir, "config.json")
	}
	
	// Try to load from file if exists
	data, err := os.ReadFile(confPath)
	if err == nil {
		err = json.Unmarshal(data, config)
		if err != nil {
			return config, fmt.Errorf("error parsing config file: %v", err)
		}
		log.Printf("Loaded configuration from file: %s", confPath)
	} else {
		log.Printf("Warning: Could not read config file, using default values: %v", err)
	}
	
	// Override with environment variables again to ensure they have highest priority
	loadConfigFromEnv(config)
	
	return config, nil
}

// loadConfigFromEnv loads configuration values from environment variables
func loadConfigFromEnv(config *Config) {
	// Pike13 URL from environment variable
	if pike13URL := os.Getenv("PIKE13_URL"); pike13URL != "" {
		config.Pike13URL = pike13URL
	}
	
	// Calendar ID from environment variable
	if calendarID := os.Getenv("CALENDAR_ID"); calendarID != "" {
		config.CalendarID = calendarID
	}
	
	// Time zone from environment variable
	if timeZone := os.Getenv("TIME_ZONE"); timeZone != "" {
		config.TimeZone = timeZone
	}
	
	// Credentials path from environment variable
	if credPath := os.Getenv("GOOGLE_CREDENTIALS_FILE"); credPath != "" {
		config.CredentialsPath = credPath
	}
	
	// Log path from environment variable
	if logPath := os.Getenv("LOG_PATH"); logPath != "" {
		config.LogPath = logPath
	}
	
	// Dry run from environment variable
	if dryRunEnv := os.Getenv("DRY_RUN"); dryRunEnv != "" {
		// Consider "true", "1", "yes", "y" as true values (case insensitive)
		dryRunValue := strings.ToLower(dryRunEnv)
		config.DryRun = dryRunValue == "true" || dryRunValue == "1" || 
		               dryRunValue == "yes" || dryRunValue == "y"
	}
}

// determineBaseDir finds the project root directory
func determineBaseDir() string {
	// If DOCKER_ENV is set, use /app
	if os.Getenv("DOCKER_ENV") == "true" {
		return "/app"
	}
	
	// If TEST_MODE is set, use a specific test directory
	if testDir := os.Getenv("TEST_BASE_DIR"); testDir != "" {
		return testDir
	}
	
	// Try to find the project root using the filename of the caller
	_, filename, _, ok := runtime.Caller(2)
	if ok {
		// Walk up until we find the .git directory or go.mod file
		dir := filepath.Dir(filename)
		for dir != "/" && dir != "." {
			// Check for project root indicators
			if fileExists(filepath.Join(dir, ".git")) || 
			   fileExists(filepath.Join(dir, "go.mod")) {
				return dir
			}
			// Move up one directory
			dir = filepath.Dir(dir)
		}
	}
	
	// Fall back to current working directory
	dir, err := os.Getwd()
	if err != nil {
		return "."
	}
	return dir
}

// fileExists checks if a file or directory exists
func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}
