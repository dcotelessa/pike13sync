package main

import (
	"bufio"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"golang.org/x/oauth2/google"
	calendar "google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

// loadEnv loads environment variables from a .env file
func loadEnv(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("error opening .env file: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		// Skip comments and empty lines
		if strings.HasPrefix(line, "#") || len(strings.TrimSpace(line)) == 0 {
			continue
		}

		// Split on first equals sign
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue // Skip malformed lines
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		// Remove quotes if present
		value = strings.Trim(value, `"'`)

		// Set environment variable
		os.Setenv(key, value)
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading .env file: %v", err)
	}

	return nil
}

func main() {
	// Load environment variables from .env file if it exists
	err := loadEnv(".env")
	if err != nil {
		fmt.Printf("Warning: Could not load .env file: %v\n", err)
		fmt.Println("Continuing with existing environment variables...")
	} else {
		fmt.Println("Loaded environment variables from .env file")
	}

	// Parse command-line flags
	checkWritePermission := false
	if len(os.Args) > 1 && os.Args[1] == "--check-write" {
		checkWritePermission = true
		fmt.Println("Running with --check-write flag: Will test write permissions")
	}

	// Set up logging
	log.SetOutput(os.Stdout)
	log.Println("Starting Google Calendar API test")

	// Get calendar ID from environment variable or use a default test value
	calendarID := os.Getenv("CALENDAR_ID")
	if calendarID == "" {
		calendarID = "primary" // Default to primary calendar
		log.Println("CALENDAR_ID not set, using 'primary'")
	} else {
		log.Printf("Using calendar ID: %s", calendarID)
	}

	// Set up Google Calendar service with flexible credential loading
	ctx := context.Background()
	srv, credBytes, err := setupGoogleCalendar(ctx)
	if err != nil {
		log.Fatalf("Error setting up Google Calendar: %v", err)
	}
	log.Println("Successfully authenticated with Google Calendar API")

	// Print service account email if available
	serviceEmail := extractServiceEmail(credBytes)
	if serviceEmail != "" {
		fmt.Printf("\nService Account Email: %s\n", serviceEmail)
		fmt.Println("Make sure this email has 'Make changes to events' permission on your calendar")
	}

	// List available calendars
	fmt.Println("\n=== AVAILABLE CALENDARS ===")
	calList, err := srv.CalendarList.List().Do()
	if err != nil {
		log.Fatalf("Error retrieving calendar list: %v", err)
	}

	if len(calList.Items) == 0 {
		fmt.Println("No calendars found accessible to this service account.")
	} else {
		for i, cal := range calList.Items {
			fmt.Printf("%d. %s (ID: %s)\n", i+1, cal.Summary, cal.Id)
		}
	}

	// Try to access the specified calendar
	fmt.Printf("\n=== ACCESSING CALENDAR: %s ===\n", calendarID)
	calendarObj, err := srv.Calendars.Get(calendarID).Do()
	if err != nil {
		fmt.Printf("Error accessing specified calendar: %v\n", err)
		fmt.Printf("Make sure the calendar ID is correct and the service account has access.\n")
		os.Exit(1)
	}
	fmt.Printf("Successfully accessed calendar: %s\n", calendarObj.Summary)

	// Get upcoming events as a test
	fmt.Println("\n=== UPCOMING EVENTS ===")
	t := time.Now().Format(time.RFC3339)
	events, err := srv.Events.List(calendarID).
		ShowDeleted(false).
		SingleEvents(true).
		TimeMin(t).
		MaxResults(10).
		OrderBy("startTime").
		Do()
	if err != nil {
		log.Fatalf("Error retrieving upcoming events: %v", err)
	}

	if len(events.Items) == 0 {
		fmt.Println("No upcoming events found.")
	} else {
		for i, item := range events.Items {
			date := item.Start.DateTime
			if date == "" {
				date = item.Start.Date
			}
			fmt.Printf("%d. %s (%s)\n", i+1, item.Summary, formatDateTime(date))
		}
	}

	// If write permission check requested, try creating and deleting a test event
	if checkWritePermission {
		fmt.Println("\n=== TESTING WRITE PERMISSIONS ===")
		fmt.Println("Creating a test event...")
		
		testEvent := &calendar.Event{
			Summary: "Permission Test - Will be Deleted",
			Description: "This is a test event created to verify service account permissions.",
			Start: &calendar.EventDateTime{
				DateTime: time.Now().Add(1 * time.Hour).Format(time.RFC3339),
				TimeZone: "America/Los_Angeles",
			},
			End: &calendar.EventDateTime{
				DateTime: time.Now().Add(2 * time.Hour).Format(time.RFC3339),
				TimeZone: "America/Los_Angeles",
			},
			ColorId: "7", // Light blue
		}
		
		createdEvent, err := srv.Events.Insert(calendarID, testEvent).Do()
		if err != nil {
			fmt.Printf("❌ ERROR - Cannot create event: %v\n", err)
			fmt.Println("The service account does not have write permission to this calendar.")
			fmt.Printf("Go to Google Calendar, share %s with %s, and give it 'Make changes to events' permission.\n", 
				calendarObj.Summary, serviceEmail)
			os.Exit(1)
		}
		
		fmt.Printf("✓ Successfully created test event: %s\n", createdEvent.HtmlLink)
		
		// Delete the test event
		fmt.Println("Cleaning up - deleting test event...")
		err = srv.Events.Delete(calendarID, createdEvent.Id).Do()
		if err != nil {
			fmt.Printf("⚠️ Warning - Could not delete test event: %v\n", err)
			fmt.Println("You may need to manually delete the test event")
		} else {
			fmt.Println("✓ Successfully deleted test event")
		}
		
		fmt.Println("\n✓ ALL TESTS PASSED - Service account has proper write access to the calendar!")
	} else {
		fmt.Println("\nRead access test completed successfully.")
		fmt.Println("To verify write permissions, run with: go run tests/test_calendar.go --check-write")
	}
}

// Use the same flexible credential loading as in main.go
func setupGoogleCalendar(ctx context.Context) (*calendar.Service, []byte, error) {
	var credBytes []byte
	var err error
	
	// Try different methods to get credentials in order of preference
	credContent := os.Getenv("GOOGLE_CREDENTIALS")
	encodedCreds := os.Getenv("GOOGLE_CREDENTIALS_BASE64")
	envPath := os.Getenv("GOOGLE_CREDENTIALS_FILE")
	
	// Method 1: Direct JSON content in environment variable
	if credContent != "" {
		log.Println("Using Google credentials from GOOGLE_CREDENTIALS environment variable")
		credBytes = []byte(credContent)
	} else if encodedCreds != "" {
		// Method 2: Base64 encoded JSON in environment variable (useful for multi-line JSON)
		log.Println("Using Google credentials from GOOGLE_CREDENTIALS_BASE64 environment variable")
		credBytes, err = base64.StdEncoding.DecodeString(encodedCreds)
		if err != nil {
			return nil, nil, fmt.Errorf("unable to decode base64 credentials: %v", err)
		}
	} else if envPath != "" {
		// Method 3: Load from file path specified in environment variable
		log.Printf("Reading Google credentials from environment-specified file: %s", envPath)
		credBytes, err = os.ReadFile(envPath)
		if err != nil {
			// If Docker path fails, try local path
			if strings.HasPrefix(envPath, "/app/") {
				localPath := "." + strings.TrimPrefix(envPath, "/app")
				log.Printf("Docker path failed, trying local path: %s", localPath)
				credBytes, err = os.ReadFile(localPath)
				if err != nil {
					return nil, nil, fmt.Errorf("unable to read credentials file: %v", err)
				}
			} else {
				return nil, nil, fmt.Errorf("unable to read credentials file: %v", err)
			}
		}
	} else {
		// Method 4: Fall back to default path
		defaultPath := "./credentials/credentials.json"
		log.Printf("Reading Google credentials from default file path: %s", defaultPath)
		credBytes, err = os.ReadFile(defaultPath)
		if err != nil {
			return nil, nil, fmt.Errorf("unable to read credentials file: %v", err)
		}
	}
	
	// Use the credentials to authenticate
	config, err := google.JWTConfigFromJSON(credBytes, calendar.CalendarScope)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to parse credentials: %v", err)
	}
	
	client := config.Client(ctx)
	srv, err := calendar.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, nil, fmt.Errorf("unable to create calendar service: %v", err)
	}
	
	return srv, credBytes, nil
}

// Helper function to format date/time for display
func formatDateTime(dateTime string) string {
	t, err := time.Parse(time.RFC3339, dateTime)
	if err != nil {
		return dateTime
	}
	return t.Format("Mon, Jan 2, 2006 at 3:04 PM")
}

// Extract service account email from credentials JSON
func extractServiceEmail(credBytes []byte) string {
	var creds map[string]interface{}
	if err := json.Unmarshal(credBytes, &creds); err != nil {
		return ""
	}
	
	if email, ok := creds["client_email"].(string); ok {
		return email
	}
	return ""
}
