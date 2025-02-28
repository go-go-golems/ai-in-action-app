package sqlite

import (
	"time"

	"gorm.io/gorm"
)

// EventModel is the GORM model for events
type EventModel struct {
	gorm.Model
	Title       string
	Speaker     string
	Description string
	Date        time.Time
	IsUpcoming  bool
}

// TableName sets the table name for EventModel
func (EventModel) TableName() string {
	return "events"
}

// TimerModel is the GORM model for timers
type TimerModel struct {
	gorm.Model
	Duration      int64 // stored in nanoseconds
	RemainingTime int64 // stored in nanoseconds
	IsRunning     bool
	LastStartedAt time.Time
}

// TableName sets the table name for TimerModel
func (TimerModel) TableName() string {
	return "timers"
}

// NoteModel is the GORM model for notes
type NoteModel struct {
	gorm.Model
	Content    string
	PageNumber int
	TotalPages int
}

// TableName sets the table name for NoteModel
func (NoteModel) TableName() string {
	return "notes"
}

// QuestionModel is the GORM model for questions
type QuestionModel struct {
	gorm.Model
	Name        string
	Content     string
	SubmittedAt time.Time
	Answered    bool
}

// TableName sets the table name for QuestionModel
func (QuestionModel) TableName() string {
	return "questions"
}
