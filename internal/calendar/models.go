package calendar

import (
	"google.golang.org/api/calendar/v3"
)

// Event is a wrapper around Google Calendar Event for any specific functionality
type Event = calendar.Event

// EventDateTime is a wrapper around Google Calendar EventDateTime
type EventDateTime = calendar.EventDateTime

// ExtendedProperties is a wrapper around Google Calendar ExtendedProperties
type ExtendedProperties = calendar.EventExtendedProperties
