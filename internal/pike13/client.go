package pike13

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/dcotelessa/pike13sync/internal/config"
	"github.com/dcotelessa/pike13sync/internal/util"
)

// Client handles interactions with Pike13 API
type Client struct {
	config     *config.Config
	testHeader map[string]string // For testing purposes only
}

// NewClient creates a new Pike13 client
func NewClient(config *config.Config) *Client {
	return &Client{
		config:     config,
		testHeader: make(map[string]string),
	}
}

// SetTestHeader sets a header for testing purposes
// This is only used in tests and is not part of the normal API
func (c *Client) SetTestHeader(key, value string) {
	if c.testHeader == nil {
		c.testHeader = make(map[string]string)
	}
	c.testHeader[key] = value
}

// FetchEvents retrieves events from Pike13 API
func (c *Client) FetchEvents(fromDate, toDate string) (Pike13Response, error) {
	var response Pike13Response
	
	// Load Pike13 credentials
	pike13Creds, err := c.loadCredentials()
	if err != nil {
		log.Printf("Warning: Could not load Pike13 credentials: %v", err)
		// Continue without credentials
	}
	
	// Construct the URL with query parameters
	url := fmt.Sprintf("%s?from=%s&to=%s", c.config.Pike13URL, fromDate, toDate)
	
	// Add client_id if available
	if pike13Creds.ClientID != "" {
		url = fmt.Sprintf("%s&client_id=%s", url, pike13Creds.ClientID)
	}
	
	log.Printf("Requesting URL: %s", url)
	
	// Make the HTTP request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return response, fmt.Errorf("error creating HTTP request: %v", err)
	}
	
	// Add headers that might help with authentication
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
	req.Header.Add("Accept", "application/json")
	
	// Add test headers if any (for testing only)
	for k, v := range c.testHeader {
		req.Header.Add(k, v)
	}
	
	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return response, fmt.Errorf("error making HTTP request: %v", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		// Try to read body for more info
		body, _ := io.ReadAll(resp.Body)
		return response, fmt.Errorf("API returned non-OK status: %d - %s", resp.StatusCode, string(body))
	}
	
	// Read and parse the response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return response, fmt.Errorf("error reading response body: %v", err)
	}
	
	// Save a copy of the raw JSON for debugging
	if c.config.LogPath != "" {
		logDir := filepath.Dir(c.config.LogPath)
		err = os.WriteFile(filepath.Join(logDir, "pike13_response.json"), body, 0644)
		if err != nil {
			log.Printf("Warning: Could not save raw API response: %v", err)
		}
	}
	
	err = json.Unmarshal(body, &response)
	if err != nil {
		return response, fmt.Errorf("error parsing JSON response: %v", err)
	}
	
	return response, nil
}

// loadCredentials loads Pike13 API credentials
func (c *Client) loadCredentials() (Pike13Credentials, error) {
	var creds Pike13Credentials
	
	// Check for environment variable first
	envClientID := os.Getenv("PIKE13_CLIENT_ID")
	if envClientID != "" {
		creds.ClientID = envClientID
		return creds, nil
	}
	
	// Fall back to file if environment variable not set
	data, err := os.ReadFile(c.config.Pike13CredentialsPath)
	if err != nil {
		return creds, fmt.Errorf("error reading Pike13 credentials file: %v", err)
	}
	
	err = json.Unmarshal(data, &creds)
	if err != nil {
		return creds, fmt.Errorf("error parsing Pike13 credentials file: %v", err)
	}
	
	return creds, nil
}

// DisplaySampleEvents prints sample events to the console
func (c *Client) DisplaySampleEvents(events Pike13Response) {
	fmt.Println("\n=== SAMPLE PIKE13 EVENTS ===")
	for i, event := range events.EventOccurrences {
		if i >= 5 { // Only show first 5 events
			fmt.Println("... (more events available)")
			break
		}
		
		fmt.Printf("\nEvent %d:\n", i+1)
		fmt.Printf("  ID: %d\n", event.ID)
		fmt.Printf("  Name: %s\n", event.Name)
		fmt.Printf("  Start: %s\n", util.FormatDateTime(event.StartAt))
		fmt.Printf("  End: %s\n", util.FormatDateTime(event.EndAt))
		fmt.Printf("  State: %s\n", event.State)
		fmt.Printf("  Full: %v\n", event.Full)
		fmt.Printf("  Capacity Remaining: %d\n", event.CapacityRemaining)
		
		if len(event.StaffMembers) > 0 {
			fmt.Printf("  Staff: ")
			for i, staff := range event.StaffMembers {
				if i > 0 {
					fmt.Printf(", ")
				}
				fmt.Printf("%s", staff.Name)
			}
			fmt.Println()
		}
		
		if event.Description != "" {
			desc := event.Description
			if len(desc) > 100 {
				desc = desc[:97] + "..."
			}
			fmt.Printf("  Description: %s\n", desc)
		}
	}
	fmt.Println("\nSample mode only - no calendar operations performed")
}
