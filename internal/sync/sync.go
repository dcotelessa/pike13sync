package sync

import (
	"log"
	"strconv"

	"github.com/dcotelessa/pike13sync/internal/calendar"
	"github.com/dcotelessa/pike13sync/internal/config"
	"github.com/dcotelessa/pike13sync/internal/pike13"
)

// SyncStats holds statistics about sync operations
type SyncStats struct {
	Created int
	Updated int
	Deleted int
	Skipped int
}

// SyncService handles synchronization between Pike13 and Google Calendar
type SyncService struct {
	calendarService *calendar.Service
	config          *config.Config
}

// NewSyncService creates a new sync service
func NewSyncService(calendarService *calendar.Service, config *config.Config) *SyncService {
	return &SyncService{
		calendarService: calendarService,
		config:          config,
	}
}

// SyncEvents synchronizes Pike13 events with Google Calendar
func (s *SyncService) SyncEvents(pike13Events []pike13.Pike13Event) SyncStats {
	stats := SyncStats{}
	
	// Get existing events from Google Calendar
	existingEvents, err := s.calendarService.GetExistingEvents()
	if err != nil {
		log.Printf("Error retrieving existing events: %v", err)
		return stats
	}
	
	// Create a map of existing events by Pike13 event ID
	existingEventMap := make(map[string]*calendar.Event)
	for _, event := range existingEvents {
		// Check if this is a Pike13 event by looking for a custom property
		if event.ExtendedProperties != nil && 
		   event.ExtendedProperties.Private != nil && 
		   event.ExtendedProperties.Private["pike13_id"] != "" {
			pike13ID := event.ExtendedProperties.Private["pike13_id"]
			existingEventMap[pike13ID] = event
		}
	}
	
	// Process Pike13 events
	for _, pike13Event := range pike13Events {
		// Format event data
		eventData := s.calendarService.FormatEventData(pike13Event)
		
		// Convert ID to string for lookup
		pike13IDStr := strconv.Itoa(pike13Event.ID)
		
		// Check if event already exists
		if existingEvent, exists := existingEventMap[pike13IDStr]; exists {
			// Update existing event if needed
			updateStatus := s.calendarService.UpdateEvent(existingEvent, eventData)
			if updateStatus == "updated" {
				stats.Updated++
			} else if updateStatus == "unchanged" {
				stats.Skipped++
			}
			// Remove from map to track what's been processed
			delete(existingEventMap, pike13IDStr)
		} else {
			// Create new event
			s.calendarService.CreateEvent(eventData)
			stats.Created++
		}
	}
	
	// Any events still in the map need to be deleted (they're no longer in Pike13)
	for _, eventToDelete := range existingEventMap {
		s.calendarService.DeleteEvent(eventToDelete)
		stats.Deleted++
	}
	
	return stats
}
