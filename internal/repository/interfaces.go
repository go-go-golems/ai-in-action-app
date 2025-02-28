package repository

import (
	"context"
	"time"

	"github.com/wesen/ai-in-action-app/internal/domain"
)

// EventRepository defines the interface for event data operations
type EventRepository interface {
	GetUpcomingEvents(ctx context.Context) ([]domain.Event, error)
	GetPastEvents(ctx context.Context) ([]domain.Event, error)
	AddEvent(ctx context.Context, event domain.Event) (domain.Event, error)
	UpdateEvent(ctx context.Context, event domain.Event) (bool, error)
}

// TimerRepository defines the interface for timer data operations
type TimerRepository interface {
	GetTimer(ctx context.Context) (domain.Timer, error)
	UpdateTimer(ctx context.Context, timer domain.Timer) (bool, error)
	ResetTimer(ctx context.Context, duration time.Duration) (domain.Timer, error)
}

// NoteRepository defines the interface for note data operations
type NoteRepository interface {
	GetNote(ctx context.Context, pageNumber int) (domain.Note, error)
	SaveNote(ctx context.Context, note domain.Note) (bool, error)
}

// QuestionRepository defines the interface for question data operations
type QuestionRepository interface {
	GetQuestions(ctx context.Context) ([]domain.Question, error)
	AddQuestion(ctx context.Context, question domain.Question) (domain.Question, error)
	MarkAsAnswered(ctx context.Context, id uint) (bool, error)
}
