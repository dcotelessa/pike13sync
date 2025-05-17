package sync_test

import (
	"strconv"
	"testing"
	
	"github.com/dcotelessa/pike13sync/internal/calendar"
	"github.com/dcotelessa/pike13sync/internal/config"
	"github.com/dcotelessa/pike13sync/internal/pike13"
	"github.com/dcotelessa/pike13sync/internal/sync"
)

// Ensure MockCalendarService implements the same interface as the calendar.Service
// This interface must match the methods called by sync.SyncService
type CalendarServiceInterface interface {
	GetExistingEvents() ([]*calendar.Event, error)
	FormatEventData(pike13.Pike13Event) *calendar.Event
	CreateEvent(*calendar.Event)
	UpdateEvent(*calendar.Event, *calendar.Event) string
	DeleteEvent(*calendar.Event)
}

// MockCalendarService implements the calendar service interface for testing
type MockCalendarService struct {
	existingEvents []*calendar.Event
	createCalls    int
	updateCalls    int
	deleteCalls    int
	skipCalls      int
}

// Ensure the mock implements the interface
var _ CalendarServiceInterface = (*MockCalendarService)(nil)

// GetExistingEvents returns mock events
func (m *MockCalendarService) GetExistingEvents() ([]*calendar.Event, error) {
	return m.existingEvents, nil
}

	// FormatEventData converts a Pike13 event to a Google Calendar event
func (m *MockCalendarService) FormatEventData(event pike13.Pike13Event) *calendar.Event {
	// Create a simple event format for testing
	// Convert ID to string properly
	idStr := strconv.Itoa(event.ID)
	
	return &calendar.Event{
		Summary: event.Name,
		ExtendedProperties: &calendar.ExtendedProperties{
			Private: map[string]string{
				"pike13_id":   idStr,
				"pike13_sync": "true",
			},
		},
	}
}

// CreateEvent mocks event creation
func (m *MockCalendarService) CreateEvent(event *calendar.Event) {
	m.createCalls++
}

// UpdateEvent mocks event updates
func (m *MockCalendarService) UpdateEvent(existing *calendar.Event, new *calendar.Event) string {
	if existing.Summary == new.Summary {
		m.skipCalls++
		return "unchanged"
	}
	m.updateCalls++
	return "updated"
}

// DeleteEvent mocks event deletion
func (m *MockCalendarService) DeleteEvent(event *calendar.Event) {
	m.deleteCalls++
}

// TestSyncEvents tests the SyncEvents function
func TestSyncEvents(t *testing.T) {
	// Create test events
	mockEvents := []*calendar.Event{
		{
			Summary: "Existing Event 1",
			ExtendedProperties: &calendar.ExtendedProperties{
				Private: map[string]string{
					"pike13_id":   "123",
					"pike13_sync": "true",
				},
			},
		},
		{
			Summary: "Existing Event 2",
			ExtendedProperties: &calendar.ExtendedProperties{
				Private: map[string]string{
					"pike13_id":   "456",
					"pike13_sync": "true",
				},
			},
		},
	}

	// Create mock calendar service
	mockCalendar := &MockCalendarService{
		existingEvents: mockEvents,
	}
	
	// Create config
	cfg := &config.Config{}
	
	// Create sync service with mock calendar
	syncService := sync.NewSyncService(mockCalendar, cfg)
	
	// Test cases
	testCases := []struct {
		name           string
		pike13Events   []pike13.Pike13Event
		expectedStats  sync.SyncStats
		expectedCreate int
		expectedUpdate int
		expectedDelete int
		expectedSkip   int
	}{
		{
			name: "No Pike13 events - delete all existing events",
			pike13Events: []pike13.Pike13Event{},
			expectedStats: sync.SyncStats{
				Created: 0,
				Updated: 0,
				Deleted: 2,
				Skipped: 0,
			},
			expectedCreate: 0,
			expectedUpdate: 0,
			expectedDelete: 2,
			expectedSkip:   0,
		},
		{
			name: "New Pike13 events - create all",
			pike13Events: []pike13.Pike13Event{
				{ID: 789, Name: "New Event 1"},
				{ID: 101, Name: "New Event 2"},
			},
			expectedStats: sync.SyncStats{
				Created: 2,
				Updated: 0,
				Deleted: 2, // The existing events are still deleted
				Skipped: 0,
			},
			expectedCreate: 2,
			expectedUpdate: 0,
			expectedDelete: 2,
			expectedSkip:   0,
		},
		{
			name: "Mixed operations - create, update, delete",
			pike13Events: []pike13.Pike13Event{
				{ID: 123, Name: "Updated Event 1"}, // Update existing with ID 123
				{ID: 789, Name: "New Event 3"},     // Create new
			},
			expectedStats: sync.SyncStats{
				Created: 1,
				Updated: 1,
				Deleted: 1, // Event with ID 456 should be deleted
				Skipped: 0,
			},
			expectedCreate: 1,
			expectedUpdate: 1,
			expectedDelete: 1,
			expectedSkip:   0,
		},
	}
	
	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Reset mock counters
			mockCalendar.createCalls = 0
			mockCalendar.updateCalls = 0
			mockCalendar.deleteCalls = 0
			mockCalendar.skipCalls = 0
			
			// Reset mock existing events for each test
			mockCalendar.existingEvents = []*calendar.Event{
				{
					Summary: "Existing Event 1",
					ExtendedProperties: &calendar.ExtendedProperties{
						Private: map[string]string{
							"pike13_id":   "123",
							"pike13_sync": "true",
						},
					},
				},
				{
					Summary: "Existing Event 2",
					ExtendedProperties: &calendar.ExtendedProperties{
						Private: map[string]string{
							"pike13_id":   "456",
							"pike13_sync": "true",
						},
					},
				},
			}
			
			// Run sync
			stats := syncService.SyncEvents(tc.pike13Events)
			
			// Verify expected stats
			if stats.Created != tc.expectedStats.Created {
				t.Errorf("Expected %d creates, got %d", tc.expectedStats.Created, stats.Created)
			}
			if stats.Updated != tc.expectedStats.Updated {
				t.Errorf("Expected %d updates, got %d", tc.expectedStats.Updated, stats.Updated)
			}
			if stats.Deleted != tc.expectedStats.Deleted {
				t.Errorf("Expected %d deletes, got %d", tc.expectedStats.Deleted, stats.Deleted)
			}
			if stats.Skipped != tc.expectedStats.Skipped {
				t.Errorf("Expected %d skips, got %d", tc.expectedStats.Skipped, stats.Skipped)
			}
			
			// Verify mock calls
			if mockCalendar.createCalls != tc.expectedCreate {
				t.Errorf("Expected %d CreateEvent calls, got %d", tc.expectedCreate, mockCalendar.createCalls)
			}
			if mockCalendar.updateCalls != tc.expectedUpdate {
				t.Errorf("Expected %d UpdateEvent calls, got %d", tc.expectedUpdate, mockCalendar.updateCalls)
			}
			if mockCalendar.deleteCalls != tc.expectedDelete {
				t.Errorf("Expected %d DeleteEvent calls, got %d", tc.expectedDelete, mockCalendar.deleteCalls)
			}
			if mockCalendar.skipCalls != tc.expectedSkip {
				t.Errorf("Expected %d skipped updates, got %d", tc.expectedSkip, mockCalendar.skipCalls)
			}
		})
	}
}
