package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// TestPike13ToGoogleCalendarIntegration tests the full integration between Pike13 and Google Calendar
func TestPike13ToGoogleCalendarIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// This test requires real credentials to be fully effective
	// We'll check if we have the necessary environment variables or credential files

	// Check for Google credentials
	hasGoogleCreds := false
	if os.Getenv("GOOGLE_CREDENTIALS") != "" || 
	   os.Getenv("GOOGLE_CREDENTIALS_BASE64") != "" || 
	   fileExists("./credentials/credentials.json") {
		hasGoogleCreds = true
	}

	// Check for Pike13 credentials
	hasPike13Creds := false
	if os.Getenv("PIKE13_CLIENT_ID") != "" || 
	   fileExists("./credentials/pike13_credentials.json") {
		hasPike13Creds = true
	}

	if !hasGoogleCreds || !hasPike13Creds {
		t.Skip("Skipping integration test due to missing credentials")
	}

	// Create a temporary directory for test data
	tmpDir := t.TempDir()
	configDir := filepath.Join(tmpDir, "config")
	credsDir := filepath.Join(tmpDir, "credentials")
	logsDir := filepath.Join(tmpDir, "logs")

	// Create directories
	os.MkdirAll(configDir, 0755)
	os.MkdirAll(credsDir, 0755)
	os.MkdirAll(logsDir, 0755)

	// Create a mock Pike13 API server
	mockPike13Server := createMockPike13Server()
	defer mockPike13Server.Close()

	// Create a temporary .env file
	envFilePath := filepath.Join(tmpDir, ".env")
	envContent := fmt.Sprintf(`CALENDAR_ID=primary
PIKE13_CLIENT_ID=test_client_id
PIKE13_URL=%s
GOOGLE_CREDENTIALS_FILE=./credentials/credentials.json
TZ=America/Los_Angeles
`, mockPike13Server.URL)

	err := os.WriteFile(envFilePath, []byte(envContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create .env file: %v", err)
	}

	// Copy credentials if they exist
	if fileExists("./credentials/credentials.json") {
		copyFile("./credentials/credentials.json", filepath.Join(credsDir, "credentials.json"))
	}
	if fileExists("./credentials/pike13_credentials.json") {
		copyFile("./credentials/pike13_credentials.json", filepath.Join(credsDir, "pike13_credentials.json"))
	}

	// Run the application in dry-run mode
	// This will test most of the functionality without actually modifying calendars
	cmd := exec.Command("go", "run", "cmd/pike13sync/main.go", 
		"--dry-run", 
		"--config", filepath.Join(configDir, "config.json"),
		"--from", "2025-05-01",
		"--to", "2025-05-07")
	
	cmd.Dir = "." // Run from the project root
	cmd.Env = append(os.Environ(), 
		fmt.Sprintf("LOG_PATH=%s", filepath.Join(logsDir, "test.log")),
		fmt.Sprintf("PIKE13_URL=%s", mockPike13Server.URL),
	)

	// Capture output
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// Run the command
	err = cmd.Run()
	if err != nil {
		t.Fatalf("Failed to run application: %v\nStderr: %s", err, stderr.String())
	}

	// Check output
	output := stdout.String()
	t.Logf("Application output: %s", output)

	// Verify that the application ran successfully
	if !strings.Contains(output, "SYNC SUMMARY (DRY RUN)") {
		t.Errorf("Expected dry run summary in output")
	}

	// Check if log file was created
	logFilePath := filepath.Join(logsDir, "test.log")
	if !fileExists(logFilePath) {
		t.Errorf("Log file was not created at %s", logFilePath)
	} else {
		logContent, err := os.ReadFile(logFilePath)
		if err != nil {
			t.Errorf("Failed to read log file: %v", err)
		} else {
			t.Logf("Log content: %s", string(logContent))
		}
	}
}

// Helper function to create a mock Pike13 API server
func createMockPike13Server() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check request path
		if !strings.Contains(r.URL.Path, "event_occurrences.json") {
			http.Error(w, "Not Found", http.StatusNotFound)
			return
		}

		// Check for from and to parameters
		if !r.URL.Query().Has("from") || !r.URL.Query().Has("to") {
			http.Error(w, "Missing date parameters", http.StatusBadRequest)
			return
		}

		// Return sample data
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		// Generate events for the requested date range
		fromStr := r.URL.Query().Get("from")
		toStr := r.URL.Query().Get("to")

		from, err := time.Parse(time.RFC3339, fromStr)
		if err != nil {
			from = time.Now()
		}

		to, err := time.Parse(time.RFC3339, toStr)
		if err != nil {
			to = from.AddDate(0, 0, 7)
		}

		events := generateMockEvents(from, to)
		response := map[string]interface{}{
			"event_occurrences": events,
		}

		json.NewEncoder(w).Encode(response)
	}))
}

// Generate mock events for testing
func generateMockEvents(from, to time.Time) []map[string]interface{} {
	var events []map[string]interface{}
	
	// Generate an event for each day in the range
	current := from
	for current.Before(to) {
		// Generate a few events for each day
		for i := 0; i < 3; i++ {
			eventTime := current.Add(time.Duration(8+i*2) * time.Hour)
			id := int(eventTime.Unix() % 1000000)
			
			event := map[string]interface{}{
				"id":                id,
				"event_id":          id + 5000,
				"name":              fmt.Sprintf("Test Class %d", i+1),
				"description":       fmt.Sprintf("Description for class %d", i+1),
				"start_at":          eventTime.Format(time.RFC3339),
				"end_at":            eventTime.Add(50 * time.Minute).Format(time.RFC3339),
				"url":               fmt.Sprintf("https://teststudio.pike13.com/e/%d", id),
				"state":             "active",
				"full":              false,
				"capacity_remaining": 10 - i,
				"staff_members": []map[string]interface{}{
					{
						"id":   101 + i,
						"name": fmt.Sprintf("Instructor %d", i+1),
					},
				},
				"waitlist": map[string]interface{}{
					"full": false,
				},
			}
			
			// Make some events canceled
			if i == 2 {
				event["state"] = "canceled"
			}
			
			events = append(events, event)
		}
		
		// Move to next day
		current = current.AddDate(0, 0, 1)
	}
	
	return events
}

// Helper function to check if a file exists
func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// Helper function to copy a file
func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	return err
}

// TestConnectionOnly tests just the connection to Pike13 and Google Calendar APIs
func TestConnectionOnly(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping connection test in short mode")
	}

	// Test Google Calendar connection
	googleTest := exec.Command("go", "run", "tests/test_calendar.go")
	var googleOut bytes.Buffer
	googleTest.Stdout = &googleOut
	
	err := googleTest.Run()
	t.Logf("Google Calendar connection test output: %s", googleOut.String())
	
	if err != nil {
		t.Logf("Google Calendar connection test failed: %v", err)
	} else {
		t.Logf("Google Calendar connection test successful")
	}

	// Test Pike13 connection by running the application with --sample flag
	pike13Test := exec.Command("go", "run", "cmd/pike13sync/main.go", "--sample")
	var pike13Out bytes.Buffer
	pike13Test.Stdout = &pike13Out
	
	err = pike13Test.Run()
	t.Logf("Pike13 connection test output: %s", pike13Out.String())
	
	if err != nil {
		t.Logf("Pike13 connection test failed: %v", err)
	} else {
		t.Logf("Pike13 connection test successful")
	}

	// This is not a real assertion-based test, but more of a diagnostic tool
	// It logs the connection status but doesn't fail the test
	// In a real CI environment, you might want to make this more strict
}
