package pike13

// Pike13Response represents the response from Pike13 API
type Pike13Response struct {
	EventOccurrences []Pike13Event `json:"event_occurrences"`
}

// Pike13Event represents an event from Pike13
type Pike13Event struct {
	ID                int           `json:"id"`
	EventID           int           `json:"event_id"`
	Name              string        `json:"name"`
	Description       string        `json:"description"`
	StartAt           string        `json:"start_at"`
	EndAt             string        `json:"end_at"`
	URL               string        `json:"url"`
	State             string        `json:"state"`
	Full              bool          `json:"full"`
	CapacityRemaining int           `json:"capacity_remaining"`
	StaffMembers      []StaffMember `json:"staff_members"`
	Waitlist          Waitlist      `json:"waitlist"`
}

// StaffMember represents a staff member assigned to an event
type StaffMember struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// Waitlist represents waitlist information for an event
type Waitlist struct {
	Full bool `json:"full"`
}

// Pike13Credentials represents credentials for Pike13 API
type Pike13Credentials struct {
	ClientID string `json:"client_id"`
}
