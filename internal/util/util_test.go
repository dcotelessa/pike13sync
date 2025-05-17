package util_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"
	
	"github.com/dcotelessa/pike13sync/internal/util"
)

// TestLoadEnvFile tests loading environment variables from a file
func TestLoadEnvFile(t *testing.T) {
	// Create a temporary env file
	tmpDir := t.TempDir()
	envFile := filepath.Join(tmpDir, ".env.test")
	
	// Write test content to the file
	content := `
# This is a comment
KEY1=value1
KEY2="value2"
KEY3='value3'
EMPTY_KEY=
   SPACES_KEY   =   spaces value   
INVALID_LINE
`
	err := os.WriteFile(envFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to create test env file: %v", err)
	}
	
	// Clear existing environment variables to ensure test works
	os.Unsetenv("KEY1")
	os.Unsetenv("KEY2")
	os.Unsetenv("KEY3")
	os.Unsetenv("EMPTY_KEY")
	os.Unsetenv("SPACES_KEY")
	
	// Test loading the file
	err = util.LoadEnvFile(envFile)
	if err != nil {
		t.Fatalf("LoadEnvFile returned error: %v", err)
	}
	
	// Check that variables were loaded correctly
	testCases := []struct {
		key      string
		expected string
	}{
		{"KEY1", "value1"},
		{"KEY2", "value2"},
		{"KEY3", "value3"},
		{"EMPTY_KEY", ""},
		{"SPACES_KEY", "spaces value"},
	}
	
	for _, tc := range testCases {
		value := os.Getenv(tc.key)
		if value != tc.expected {
			t.Errorf("Expected %s=%s, got %s", tc.key, tc.expected, value)
		}
	}
	
	// Test error case with non-existent file
	err = util.LoadEnvFile(filepath.Join(tmpDir, "nonexistent.env"))
	if err == nil {
		t.Error("Expected error for non-existent file, got nil")
	}
}

// TestFormatDateTime tests the FormatDateTime function
func TestFormatDateTime(t *testing.T) {
	// Test with valid datetime
	testCases := []struct {
		input    string
		expected string
	}{
		{"2025-05-15T14:30:00Z", "May 15 at 2:30 PM"}, // UTC time
		{"2025-05-15T14:30:00-07:00", "May 15 at 2:30 PM"}, // Pacific time
		{"2025-01-01T00:00:00Z", "Jan 1 at 12:00 AM"}, 
		{"invalid-date", "invalid-date"}, // Should return input string for invalid input
	}
	
	for _, tc := range testCases {
		result := util.FormatDateTime(tc.input)
		if result != tc.expected {
			t.Errorf("FormatDateTime(%s) = %s, expected %s", tc.input, result, tc.expected)
		}
	}
}

// TestGetStartAndEndOfWeek tests the GetStartAndEndOfWeek function
func TestGetStartAndEndOfWeek(t *testing.T) {
	start, end := util.GetStartAndEndOfWeek()
	
	// Parse the returned dates
	startDate, err := time.Parse(time.RFC3339, start)
	if err != nil {
		t.Errorf("Failed to parse start date: %v", err)
	}
	
	endDate, err := time.Parse(time.RFC3339, end)
	if err != nil {
		t.Errorf("Failed to parse end date: %v", err)
	}
	
	// Check that the start date is a Sunday (0 in Go's time.Weekday)
	if startDate.Weekday() != time.Sunday {
		t.Errorf("Expected start date to be Sunday, got %s", startDate.Weekday())
	}
	
	// Check that the end date is 7 days after start date
	expectedEnd := startDate.AddDate(0, 0, 7)
	if !endDate.Equal(expectedEnd) {
		t.Errorf("Expected end date to be %v, got %v", expectedEnd, endDate)
	}
}
