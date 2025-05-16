package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

// Config holds application configuration
type Config struct {
	Pike13URL             string `json:"pike13_url"`
	CalendarID            string `json:"calendar_id"`
	TimeZone              string `json:"time_zone"`
	CredentialsPath       string `json:"credentials_path"`
	Pike13CredentialsPath string `json:"pike13_credentials_path"`
	LogPath               string `json:"log_path"`
	DryRun                bool   `json:"dry_run"`
	BaseDir               string `json:"-"` // Not serialized
}

// LoadConfig loads configuration from file and environment variables
func LoadConfig(configPath string) (*Config, error) {
	config := &Config{
		TimeZone:              "America/Los_Angeles",
		DryRun:                false,
	}
	
	// Determine base directory
	if os.Getenv("DOCKER_ENV") == "true" {
		config.BaseDir = "/app"
	} else {
		var err error
		config.BaseDir, err = os.Getwd()
		if err != nil {
			config.BaseDir = "."
		}
	}
	
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
	
	config.Pike13CredentialsPath = filepath.Join(credentialsDir, "pike13_credentials.json")
	config.LogPath = filepath.Join(logsDir, "pike13sync.log")
	
	// Create directories if they don't exist
	os.MkdirAll(configDir, 0755)
	os.MkdirAll(credentialsDir, 0755)
	os.MkdirAll(logsDir, 0755)
	
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
	} else {
		log.Printf("Warning: Could not read config file, using default values: %v", err)
	}
	
	// Override from environment variables
	if os.Getenv("PIKE13_URL") != "" {
		config.Pike13URL = os.Getenv("PIKE13_URL")
	}
	if os.Getenv("CALENDAR_ID") != "" {
		config.CalendarID = os.Getenv("CALENDAR_ID")
	}
	if os.Getenv("TIME_ZONE") != "" {
		config.TimeZone = os.Getenv("TIME_ZONE")
	}
	if os.Getenv("GOOGLE_CREDENTIALS_FILE") != "" {
		config.CredentialsPath = os.Getenv("GOOGLE_CREDENTIALS_FILE")
	}
	if os.Getenv("DRY_RUN") == "true" {
		config.DryRun = true
	}
	
	return config, nil
}
