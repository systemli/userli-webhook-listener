package main

import "time"

const (
	// EventTypeUserCreated is the event type for user creation
	EventTypeUserCreated = "user.created"
	// EventTypeUserDeleted is the event type for user deletion
	EventTypeUserDeleted = "user.deleted"
)

// UserEvent represents a user event in the system
// It contains the event type, timestamp, and user data
// The event type can be either "user.created" or "user.deleted"
// The timestamp indicates when the event occurred
// The user data contains the email of the user involved in the event
type UserEvent struct {
	Type      string    `json:"type"`
	Timestamp time.Time `json:"timestamp"`
	Data      struct {
		Email string `json:"email"`
	} `json:"data"`
}
