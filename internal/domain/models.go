package domain

import (
	"time"
)

// Event represents a talk or workshop in the AI in Action group
type Event struct {
	ID          uint
	Title       string
	Speaker     string
	Description string
	Date        time.Time
	IsUpcoming  bool
}

// Timer represents a countdown timer for talks
type Timer struct {
	ID            uint
	Duration      time.Duration
	RemainingTime time.Duration
	IsRunning     bool
	LastStartedAt time.Time
}

// Note represents speaker notes for a talk
type Note struct {
	ID         uint
	Content    string
	PageNumber int
	TotalPages int
}

// Question represents a question submitted by an attendee
type Question struct {
	ID          uint
	Name        string
	Content     string
	SubmittedAt time.Time
	Answered    bool
} 