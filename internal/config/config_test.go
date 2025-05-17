package config_test

import (
	"os"
	"path/filepath"
	"testing"
	
	"github.com/dcotelessa/pike13sync/internal/config"
)

// TestLoadConfig tests loading configuration from file and environment variables
func TestLoadConfig(t *testing.T) {
	// Create a temporary directory for test files
	tmpDir := t.TempDir()
	
	// Create a test config file
	configDir := filepath.Join(tmpDir, "config")
	err := os.MkdirAll(configDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create config directory: %v", err)
	}
	
	configPath := filepath.Join(configDir, "config.json")
	configContent := `{
		"pike13_url": "https://teststudio.pike13.com/api/v2/front/event_occurrences.json",
		"calendar_id": "test_calendar_id@group.calendar.google.com",
		"time_zone": "Europe/London",
		"credentials_path": "/test/path/credentials.json",
		"pike13_credentials_path": "/test/path/pike13_credentials.json",
		"log_path": "/test/path/logs/pike13sync.log",
		"dry_run": true
	}`
	
	err = os.WriteFile(configPath, []byte(configContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test config file: %v", err)
	}
	
	// Test with clean environment
	os.Unsetenv("DOCKER_ENV")
	os.Unsetenv("PIKE13_URL")
	os.Unsetenv("CALENDAR_ID")
	os.Unsetenv("TIME_ZONE")
	os.Unsetenv("GOOGLE_CREDENTIALS_FILE")
	os.Unsetenv("DRY_RUN")
	
	// Test loading without environment variables
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		t.Fatalf("LoadConfig returned error: %v", err)
	}
	
	// Verify config values from file
	if cfg.Pike13URL != "https://teststudio.pike13.com/api/v2/front/event_occurrences.json" {
		t.Errorf("Expected Pike13URL=%s, got %s", "https://teststudio.pike13.com/api/v2/front/event_occurrences.json", cfg.Pike13URL)
	}
	if cfg.CalendarID != "test_calendar_id@group.calendar.google.com" {
		t.Errorf("Expected CalendarID=%s, got %s", "test_calendar_id@group.calendar.google.com", cfg.CalendarID)
	}
	if cfg.TimeZone != "Europe/London" {
		t.Errorf("Expected TimeZone=%s, got %s", "Europe/London", cfg.TimeZone)
	}
	if cfg.DryRun != true {
		t.Errorf("Expected DryRun=true, got %v", cfg.DryRun)
	}
	
	// Test environment variable overrides
	os.Setenv("PIKE13_URL", "https://override.pike13.com/api/v2")
	os.Setenv("CALENDAR_ID", "override_calendar@group.calendar.google.com")
	os.Setenv("TIME_ZONE", "America/New_York")
	os.Setenv("GOOGLE_CREDENTIALS_FILE", "/override/path/credentials.json")
	
	// Important: Use "false" string for DRY_RUN to properly test boolean env var
	os.Setenv("DRY_RUN", "false")
	
	// Load config again
	cfg, err = config.LoadConfig(configPath)
	if err != nil {
		t.Fatalf("LoadConfig with env vars returned error: %v", err)
	}
	
	// Verify environment variables override file config
	if cfg.Pike13URL != "https://override.pike13.com/api/v2" {
		t.Errorf("Expected Pike13URL to be overridden to %s, got %s", "https://override.pike13.com/api/v2", cfg.Pike13URL)
	}
	if cfg.CalendarID != "override_calendar@group.calendar.google.com" {
		t.Errorf("Expected CalendarID to be overridden to %s, got %s", "override_calendar@group.calendar.google.com", cfg.CalendarID)
	}
	if cfg.TimeZone != "America/New_York" {
		t.Errorf("Expected TimeZone to be overridden to %s, got %s", "America/New_York", cfg.TimeZone)
	}
	if cfg.CredentialsPath != "/override/path/credentials.json" {
		t.Errorf("Expected CredentialsPath to be overridden to %s, got %s", "/override/path/credentials.json", cfg.CredentialsPath)
	}
	if cfg.DryRun != false {
		t.Errorf("Expected DryRun to be overridden to false, got %v", cfg.DryRun)
	}
	
	// Test Docker environment detection
	os.Setenv("DOCKER_ENV", "true")
	cfg, err = config.LoadConfig(configPath)
	if err != nil {
		t.Fatalf("LoadConfig with DOCKER_ENV returned error: %v", err)
	}
	
	if cfg.BaseDir != "/app" {
		t.Errorf("Expected BaseDir to be /app in Docker environment, got %s", cfg.BaseDir)
	}
	
	// Test with non-existent config file
	nonexistentPath := filepath.Join(tmpDir, "nonexistent.json")
	cfg, err = config.LoadConfig(nonexistentPath)
	// Should not return error, but use default values
	if err != nil {
		t.Fatalf("LoadConfig with nonexistent file returned error: %v", err)
	}
	
	// Verify default values are used
	if cfg.Pike13URL != "https://override.pike13.com/api/v2" { // from env var
		t.Errorf("Expected default Pike13URL from env var, got %s", cfg.Pike13URL)
	}
	if cfg.TimeZone != "America/New_York" { // from env var
		t.Errorf("Expected default TimeZone from env var, got %s", cfg.TimeZone)
	}
}
