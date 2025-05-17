package pike13_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	
	"github.com/dcotelessa/pike13sync/internal/config"
	"github.com/dcotelessa/pike13sync/internal/pike13"
)

// TestFetchEvents tests the FetchEvents function
func TestFetchEvents(t *testing.T) {
	// Create a test server that responds with a sample Pike13 response
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check for required query parameters
		if !r.URL.Query().Has("from") || !r.URL.Query().Has("to") {
			t.Error("Missing required date query parameters")
			http.Error(w, "Missing date parameters", http.StatusBadRequest)
			return
		}
		
		// Check for client_id if it's in the request
		if clientID := r.URL.Query().Get("client_id"); clientID != "" {
			// Get expected client ID based on test case
			expectedClientID := r.Header.Get("X-Expected-Client-ID")
			if expectedClientID != "" && clientID != expectedClientID {
				t.Errorf("Expected client_id=%s, got %s", expectedClientID, clientID)
			}
		}
		
		// Return sample Pike13 response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"event_occurrences": [
				{
					"id": 123456,
					"event_id": 789,
					"name": "Test Class",
					"description": "Test description",
					"start_at": "2025-05-15T14:00:00Z",
					"end_at": "2025-05-15T15:00:00Z",
					"url": "https://teststudio.pike13.com/e/123456",
					"state": "active",
					"full": false,
					"capacity_remaining": 5,
					"staff_members": [
						{
							"id": 101,
							"name": "Test Instructor"
						}
					],
					"waitlist": {
						"full": false
					}
				}
			]
		}`))
	}))
	defer ts.Close()
	
	// Create a temporary directory for test credentials
	tmpDir := t.TempDir()
	credsDir := filepath.Join(tmpDir, "credentials")
	err := os.MkdirAll(credsDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create credentials directory: %v", err)
	}
	
	// Create a test Pike13 credentials file
	pike13CredsPath := filepath.Join(credsDir, "pike13_credentials.json")
	pike13CredsContent := `{
		"client_id": "test_client_id"
	}`
	
	err = os.WriteFile(pike13CredsPath, []byte(pike13CredsContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test Pike13 credentials file: %v", err)
	}
	
	// Clear any existing environment variables
	os.Unsetenv("PIKE13_CLIENT_ID")
	
	// Create test config
	cfg := &config.Config{
		Pike13URL:             ts.URL,
		Pike13CredentialsPath: pike13CredsPath,
	}
	
	// Create Pike13 client
	client := pike13.NewClient(cfg)
	
	// Set expected header for the test server
	client.SetTestHeader("X-Expected-Client-ID", "test_client_id")
	
	// Test FetchEvents
	fromDate := "2025-05-01T00:00:00Z"
	toDate := "2025-05-31T00:00:00Z"
	response, err := client.FetchEvents(fromDate, toDate)
	
	if err != nil {
		t.Fatalf("FetchEvents returned error: %v", err)
	}
	
	// Verify response
	if len(response.EventOccurrences) != 1 {
		t.Fatalf("Expected 1 event, got %d", len(response.EventOccurrences))
	}
	
	event := response.EventOccurrences[0]
	if event.ID != 123456 {
		t.Errorf("Expected event ID=123456, got %d", event.ID)
	}
	if event.Name != "Test Class" {
		t.Errorf("Expected event Name='Test Class', got %s", event.Name)
	}
	if event.StartAt != "2025-05-15T14:00:00Z" {
		t.Errorf("Expected StartAt='2025-05-15T14:00:00Z', got %s", event.StartAt)
	}
	if event.State != "active" {
		t.Errorf("Expected State='active', got %s", event.State)
	}
	if event.CapacityRemaining != 5 {
		t.Errorf("Expected CapacityRemaining=5, got %d", event.CapacityRemaining)
	}
	if len(event.StaffMembers) != 1 || event.StaffMembers[0].Name != "Test Instructor" {
		t.Errorf("Staff member data incorrect")
	}
	
	// Test with environment variable client ID
	os.Setenv("PIKE13_CLIENT_ID", "env_client_id")
	client.SetTestHeader("X-Expected-Client-ID", "env_client_id")
	
	// Create a new client with an invalid credentials path to force using env var
	cfg.Pike13CredentialsPath = filepath.Join(tmpDir, "nonexistent.json")
	client = pike13.NewClient(cfg)
	client.SetTestHeader("X-Expected-Client-ID", "env_client_id")
	
	// Test FetchEvents with env var client ID
	response, err = client.FetchEvents(fromDate, toDate)
	if err != nil {
		t.Fatalf("FetchEvents with env var returned error: %v", err)
	}
	
	// Verify response
	if len(response.EventOccurrences) != 1 {
		t.Fatalf("Expected 1 event with env var, got %d", len(response.EventOccurrences))
	}
	
	// Test error handling - invalid URL
	cfg.Pike13URL = "http://invalid-url-that-does-not-exist"
	client = pike13.NewClient(cfg)
	
	_, err = client.FetchEvents(fromDate, toDate)
	if err == nil {
		t.Error("Expected error for invalid URL, got nil")
	}
}

// TestLoadCredentials tests the loadCredentials function
func TestLoadCredentials(t *testing.T) {
	// Create a temporary directory
	tmpDir := t.TempDir()
	credsDir := filepath.Join(tmpDir, "credentials")
	err := os.MkdirAll(credsDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create credentials directory: %v", err)
	}
	
	// Clear any existing environment variables
	os.Unsetenv("PIKE13_CLIENT_ID")
	
	// Create a test Pike13 credentials file
	pike13CredsPath := filepath.Join(credsDir, "pike13_credentials.json")
	pike13CredsContent := `{
		"client_id": "file_client_id"
	}`
	
	err = os.WriteFile(pike13CredsPath, []byte(pike13CredsContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test Pike13 credentials file: %v", err)
	}
	
	// Test with file
	cfg := &config.Config{
		Pike13CredentialsPath: pike13CredsPath,
	}
	
	client := pike13.NewClient(cfg)
	client.SetTestHeader("X-Expected-Client-ID", "file_client_id")
	
	// Create a test server that validates the client ID
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		clientID := r.URL.Query().Get("client_id")
		expectedClientID := r.Header.Get("X-Expected-Client-ID")
		if expectedClientID != "" && clientID != expectedClientID {
			t.Errorf("Expected client_id=%s, got %s", expectedClientID, clientID)
			http.Error(w, "Invalid client ID", http.StatusBadRequest)
			return
		}
		
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"event_occurrences": []}`))
	}))
	defer ts.Close()
	
	cfg.Pike13URL = ts.URL
	
	// Test with client ID from file
	_, err = client.FetchEvents("2025-05-01T00:00:00Z", "2025-05-31T00:00:00Z")
	if err != nil {
		t.Fatalf("FetchEvents with file creds returned error: %v", err)
	}
	
	// Test with environment variable
	os.Setenv("PIKE13_CLIENT_ID", "env_client_id")
	
	// Create a new test server for the env var test
	tsEnv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		clientID := r.URL.Query().Get("client_id")
		expectedClientID := r.Header.Get("X-Expected-Client-ID")
		if expectedClientID != "" && clientID != expectedClientID {
			t.Errorf("Expected client_id=%s, got %s", expectedClientID, clientID)
			http.Error(w, "Invalid client ID", http.StatusBadRequest)
			return
		}
		
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"event_occurrences": []}`))
	}))
	defer tsEnv.Close()
	
	// Create a new client with an invalid credentials path to force using env var
	cfg.Pike13URL = tsEnv.URL
	cfg.Pike13CredentialsPath = filepath.Join(tmpDir, "nonexistent.json")
	client = pike13.NewClient(cfg)
	client.SetTestHeader("X-Expected-Client-ID", "env_client_id")
	
	// Test with client ID from environment variable
	_, err = client.FetchEvents("2025-05-01T00:00:00Z", "2025-05-31T00:00:00Z")
	if err != nil {
		t.Fatalf("FetchEvents with env creds returned error: %v", err)
	}
}
