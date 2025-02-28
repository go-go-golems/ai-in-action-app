package repository

import (
	"time"

	"github.com/wesen/ai-in-action-app/internal/domain"
)

// EventRepository defines the interface for event data operations
type EventRepository interface {
	GetUpcomingEvents() []domain.Event
	GetPastEvents() []domain.Event
	AddEvent(event domain.Event) domain.Event
	UpdateEvent(event domain.Event) bool
}

// TimerRepository defines the interface for timer data operations
type TimerRepository interface {
	GetTimer() domain.Timer
	UpdateTimer(timer domain.Timer) bool
	ResetTimer(duration time.Duration) domain.Timer
}

// NoteRepository defines the interface for note data operations
type NoteRepository interface {
	GetNote(pageNumber int) domain.Note
	SaveNote(note domain.Note) bool
}

// QuestionRepository defines the interface for question data operations
type QuestionRepository interface {
	GetQuestions() []domain.Question
	AddQuestion(question domain.Question) domain.Question
	MarkAsAnswered(id uint) bool
} 