package calendar

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
	
	"github.com/dcotelessa/pike13sync/internal/config"
	"github.com/dcotelessa/pike13sync/internal/pike13"
	"github.com/dcotelessa/pike13sync/internal/util"
)

// Service handles interactions with Google Calendar
type Service struct {
	calendarService *calendar.Service
	config          *config.Config
}

// NewService creates a new calendar service
func NewService(config *config.Config) (*Service, error) {
	ctx := context.Background()
	
	// Set up Google Calendar service
	calendarService, err := setupGoogleCalendar(ctx, config.CredentialsPath)
	if err != nil {
		return nil, err
	}
	
	return &Service{
		calendarService: calendarService,
		config:          config,
	}, nil
}

// GetExistingEvents retrieves events from Google Calendar
func (s *Service) GetExistingEvents() ([]*calendar.Event, error) {
	// Get events from the past week and the next two weeks
	timeMin := time.Now().AddDate(0, 0, -7).Format(time.RFC3339)
	timeMax := time.Now().AddDate(0, 0, 14).Format(time.RFC3339)
	
	events, err := s.calendarService.Events.List(s.config.CalendarID).
		TimeMin(timeMin).
		TimeMax(timeMax).
		SingleEvents(true).
		OrderBy("startTime").
		Do()
	if err != nil {
		return nil, fmt.Errorf("error retrieving existing events: %v", err)
	}
	
	return events.Items, nil
}

// FormatEventData creates a Google Calendar event from Pike13 event data
func (s *Service) FormatEventData(pike13Event pike13.Pike13Event) *calendar.Event {
	// Format description
	description := ""
	if pike13Event.Description != "" {
		description += pike13Event.Description + "\n\n"
	}
	
	// Add status information
	status := "Active"
	if pike13Event.State != "active" {
		status = "Cancelled"
	}
	description += fmt.Sprintf("Status: %s\n", status)
	
	// Add capacity information
	capacityInfo := fmt.Sprintf("Spaces available: %d", pike13Event.CapacityRemaining)
	if pike13Event.Full {
		capacityInfo = "Class is FULL"
	}
	description += fmt.Sprintf("Capacity: %s\n", capacityInfo)
	
	// Add waitlist info
	waitlistInfo := "Waitlist is OPEN"
	if pike13Event.Waitlist.Full {
		waitlistInfo = "Waitlist is FULL"
	}
	description += fmt.Sprintf("Waitlist: %s\n", waitlistInfo)
	
	// Add staff members
	var staffNames []string
	for _, staff := range pike13Event.StaffMembers {
		staffNames = append(staffNames, staff.Name)
	}
	if len(staffNames) > 0 {
		description += "Instructor(s): "
		for i, name := range staffNames {
			if i > 0 {
				description += ", "
			}
			description += name
		}
		description += "\n"
	}
	
	// Determine color based on state
	colorId := "11" // Red for active
	if pike13Event.State != "active" {
		colorId = "8" // Gray for cancelled
	}
	
	// Convert ID to string for extended properties
	pike13IDStr := strconv.Itoa(pike13Event.ID)
	
	// Create the event object
	event := &calendar.Event{
		Summary:     pike13Event.Name,
		Description: description,
		Start: &calendar.EventDateTime{
			DateTime: pike13Event.StartAt,
			TimeZone: s.config.TimeZone,
		},
		End: &calendar.EventDateTime{
			DateTime: pike13Event.EndAt,
			TimeZone: s.config.TimeZone,
		},
		ColorId: colorId,
		Source: &calendar.EventSource{
			Title: "View on Pike13",
			Url:   pike13Event.URL,
		},
		ExtendedProperties: &calendar.EventExtendedProperties{
			Private: map[string]string{
				"pike13_id":   pike13IDStr,
				"pike13_sync": "true",
			},
		},
	}
	
	return event
}

// CreateEvent creates a new event in Google Calendar
func (s *Service) CreateEvent(event *calendar.Event) {
	if s.config.DryRun {
		fmt.Printf("Would CREATE: %s (%s to %s)\n", 
			event.Summary, 
			util.FormatDateTime(event.Start.DateTime), 
			util.FormatDateTime(event.End.DateTime))
		return
	}
	
	_, err := s.calendarService.Events.Insert(s.config.CalendarID, event).Do()
	if err != nil {
		log.Printf("Error creating event '%s': %v", event.Summary, err)
		return
	}
	log.Printf("Created event: %s", event.Summary)
}

// UpdateEvent updates an existing event in Google Calendar
func (s *Service) UpdateEvent(existingEvent *calendar.Event, newEventData *calendar.Event) string {
	// Check if update is needed by comparing key fields
	needsUpdate := false
	
	if existingEvent.Summary != newEventData.Summary {
		needsUpdate = true
	}
	if existingEvent.Description != newEventData.Description {
		needsUpdate = true
	}
	if existingEvent.Start.DateTime != newEventData.Start.DateTime {
		needsUpdate = true
	}
	if existingEvent.End.DateTime != newEventData.End.DateTime {
		needsUpdate = true
	}
	if existingEvent.ColorId != newEventData.ColorId {
		needsUpdate = true
	}
	
	// Only update if changes detected
	if needsUpdate {
		if s.config.DryRun {
			fmt.Printf("Would UPDATE: %s (%s to %s)\n", 
				newEventData.Summary, 
				util.FormatDateTime(newEventData.Start.DateTime), 
				util.FormatDateTime(newEventData.End.DateTime))
			return "updated"
		}
		
		// Preserve the Google Calendar event ID
		newEventData.Id = existingEvent.Id
		
		_, err := s.calendarService.Events.Update(s.config.CalendarID, existingEvent.Id, newEventData).Do()
		if err != nil {
			log.Printf("Error updating event '%s': %v", newEventData.Summary, err)
			return "error"
		}
		log.Printf("Updated event: %s", newEventData.Summary)
		return "updated"
	} else {
		if s.config.DryRun {
			fmt.Printf("Would SKIP (no changes): %s (%s to %s)\n", 
				existingEvent.Summary, 
				util.FormatDateTime(existingEvent.Start.DateTime), 
				util.FormatDateTime(existingEvent.End.DateTime))
		} else {
			log.Printf("No changes needed for event: %s", existingEvent.Summary)
		}
		return "unchanged"
	}
}

// DeleteEvent deletes an event from Google Calendar
func (s *Service) DeleteEvent(event *calendar.Event) {
	if s.config.DryRun {
		fmt.Printf("Would DELETE: %s (%s to %s)\n", 
			event.Summary, 
			util.FormatDateTime(event.Start.DateTime), 
			util.FormatDateTime(event.End.DateTime))
		return
	}
	
	err := s.calendarService.Events.Delete(s.config.CalendarID, event.Id).Do()
	if err != nil {
		log.Printf("Error deleting event '%s': %v", event.Summary, err)
		return
	}
	log.Printf("Deleted event: %s", event.Summary)
}

// setupGoogleCalendar creates a Google Calendar service
func setupGoogleCalendar(ctx context.Context, credentialsPath string) (*calendar.Service, error) {
	var credBytes []byte
	var err error
	
	log.Printf("Attempting to read credentials from: %s", credentialsPath)
	
	// Try to read the file
	credBytes, err = os.ReadFile(credentialsPath)
	
	// If file read fails, try to adapt Docker path to local path
	if err != nil && strings.HasPrefix(credentialsPath, "/app/") {
		localPath := "." + strings.TrimPrefix(credentialsPath, "/app")
		log.Printf("Docker path failed, trying local equivalent: %s", localPath)
		credBytes, err = os.ReadFile(localPath)
	}
	
	// If both fail, check for environment variables
	if err != nil {
		// Check for credentials in environment variables
		if credContent := os.Getenv("GOOGLE_CREDENTIALS"); credContent != "" {
			log.Println("Using Google credentials from GOOGLE_CREDENTIALS environment variable")
			credBytes = []byte(credContent)
			err = nil
		} else if encodedCreds := os.Getenv("GOOGLE_CREDENTIALS_BASE64"); encodedCreds != "" {
			log.Println("Using Google credentials from GOOGLE_CREDENTIALS_BASE64 environment variable")
			credBytes, err = base64.StdEncoding.DecodeString(encodedCreds)
		}
	}
	
	// If still having issues, return error
	if err != nil {
		return nil, fmt.Errorf("unable to read credentials: %v", err)
	}
	
	// Use the credentials to authenticate
	config, err := google.JWTConfigFromJSON(credBytes, calendar.CalendarScope)
	if err != nil {
		return nil, fmt.Errorf("unable to parse credentials: %v", err)
	}
	
	client := config.Client(ctx)
	srv, err := calendar.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, fmt.Errorf("unable to create calendar service: %v", err)
	}
	
	return srv, nil
}
