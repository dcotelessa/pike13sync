package tests

import (
	"os"
	"testing"
)

// TestCalendar tests the calendar package
// Stub file to maintain package consistency
// Actual calendar testing is in test_calendar.go run as a separate binary
func TestCalendar(t *testing.T) {
	// Skip in short mode
	if testing.Short() {
		t.Skip("Skipping calendar test in short mode")
	}
	
	// Check if Google credentials are available
	if os.Getenv("GOOGLE_CREDENTIALS") == "" && 
	   os.Getenv("GOOGLE_CREDENTIALS_BASE64") == "" && 
	   !fileExists("./credentials/credentials.json") {
		t.Skip("Skipping calendar test due to missing credentials")
	}
	
	// This is just a stub - the real test is run by TestConnectionOnly
	// which executes test_calendar.go as a separate process
	t.Log("Calendar tests are run via test_calendar.go")
}
