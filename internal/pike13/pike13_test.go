package pike13_test

import (
	"net/http"
	"net/http/httptest"
	"os"
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
	
	// Save current environment and restore it after the test
	oldClientID := os.Getenv("PIKE13_CLIENT_ID")
	defer os.Setenv("PIKE13_CLIENT_ID", oldClientID)
	
	// Set environment variable for testing
	os.Setenv("PIKE13_CLIENT_ID", "test_client_id")
	
	// Create test config
	cfg := &config.Config{
		Pike13URL: ts.URL,
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
	
	// Test error handling - missing client ID
	os.Unsetenv("PIKE13_CLIENT_ID")
	_, err = client.FetchEvents(fromDate, toDate)
	if err == nil {
		t.Error("Expected error for missing client ID, got nil")
	}
	
	// Test error handling - invalid URL
	os.Setenv("PIKE13_CLIENT_ID", "test_client_id")
	cfg.Pike13URL = "http://invalid-url-that-does-not-exist"
	client = pike13.NewClient(cfg)
	
	_, err = client.FetchEvents(fromDate, toDate)
	if err == nil {
		t.Error("Expected error for invalid URL, got nil")
	}
}
