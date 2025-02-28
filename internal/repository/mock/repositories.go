package mock

import (
	"context"
	"sync"
	"time"

	"github.com/wesen/ai-in-action-app/internal/domain"
	"github.com/wesen/ai-in-action-app/internal/repository"
)

// MockEventRepository implements the EventRepository interface with in-memory storage
type MockEventRepository struct {
	events []domain.Event
	mu     sync.RWMutex
	nextID uint
}

var _ repository.EventRepository = &MockEventRepository{}

// NewMockEventRepository creates a new mock event repository with sample data
func NewMockEventRepository() *MockEventRepository {
	repo := &MockEventRepository{
		events: make([]domain.Event, 0),
		nextID: 1,
	}

	// Add sample upcoming events
	repo.events = append(repo.events, domain.Event{
		ID:          repo.nextID,
		Title:       "Generative AI for Scientific Discovery",
		Speaker:     "Dr. Alex Chen",
		Description: "Exploring how generative AI models can accelerate scientific discovery in various domains.",
		Date:        time.Date(2025, 3, 6, 18, 0, 0, 0, time.Local),
		IsUpcoming:  true,
	})
	repo.nextID++

	repo.events = append(repo.events, domain.Event{
		ID:          repo.nextID,
		Title:       "Ethical Considerations in AI Development",
		Speaker:     "Prof. Maya Johnson",
		Description: "Discussing the ethical frameworks necessary for responsible AI development.",
		Date:        time.Date(2025, 3, 13, 18, 0, 0, 0, time.Local),
		IsUpcoming:  true,
	})
	repo.nextID++

	// Add sample past events
	repo.events = append(repo.events, domain.Event{
		ID:          repo.nextID,
		Title:       "Multimodal Learning in AI",
		Speaker:     "Sam Rodriguez",
		Description: "How combining different data modalities can enhance AI model capabilities.",
		Date:        time.Date(2025, 2, 20, 18, 0, 0, 0, time.Local),
		IsUpcoming:  false,
	})
	repo.nextID++

	repo.events = append(repo.events, domain.Event{
		ID:          repo.nextID,
		Title:       "Reinforcement Learning from Human Feedback",
		Speaker:     "Dr. Jamie Park",
		Description: "Deep dive into how RLHF is transforming the alignment of AI systems.",
		Date:        time.Date(2025, 2, 13, 18, 0, 0, 0, time.Local),
		IsUpcoming:  false,
	})
	repo.nextID++

	return repo
}

// GetUpcomingEvents returns all upcoming events
func (m *MockEventRepository) GetUpcomingEvents(ctx context.Context) ([]domain.Event, error) {
	// Check if context is done
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	upcomingEvents := make([]domain.Event, 0)
	for _, event := range m.events {
		if event.IsUpcoming {
			upcomingEvents = append(upcomingEvents, event)
		}
	}
	return upcomingEvents, nil
}

// GetPastEvents returns all past events
func (m *MockEventRepository) GetPastEvents(ctx context.Context) ([]domain.Event, error) {
	// Check if context is done
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	pastEvents := make([]domain.Event, 0)
	for _, event := range m.events {
		if !event.IsUpcoming {
			pastEvents = append(pastEvents, event)
		}
	}
	return pastEvents, nil
}

// AddEvent adds a new event and returns it with an ID
func (m *MockEventRepository) AddEvent(ctx context.Context, event domain.Event) (domain.Event, error) {
	// Check if context is done
	if ctx.Err() != nil {
		return domain.Event{}, ctx.Err()
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	event.ID = m.nextID
	m.nextID++
	m.events = append(m.events, event)
	return event, nil
}

// UpdateEvent updates an existing event
func (m *MockEventRepository) UpdateEvent(ctx context.Context, event domain.Event) (bool, error) {
	// Check if context is done
	if ctx.Err() != nil {
		return false, ctx.Err()
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	for i, e := range m.events {
		if e.ID == event.ID {
			m.events[i] = event
			return true, nil
		}
	}
	return false, nil
}

// MockTimerRepository implements the TimerRepository interface with in-memory storage
type MockTimerRepository struct {
	timer domain.Timer
	mu    sync.RWMutex
}

var _ repository.TimerRepository = &MockTimerRepository{}

// NewMockTimerRepository creates a new mock timer repository
func NewMockTimerRepository() *MockTimerRepository {
	return &MockTimerRepository{
		timer: domain.Timer{
			ID:            1,
			Duration:      15 * time.Minute,
			RemainingTime: 15 * time.Minute,
			IsRunning:     false,
		},
	}
}

// GetTimer returns the current timer
func (m *MockTimerRepository) GetTimer(ctx context.Context) (domain.Timer, error) {
	// Check if context is done
	if ctx.Err() != nil {
		return domain.Timer{}, ctx.Err()
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	// If timer is running, calculate the remaining time
	if m.timer.IsRunning {
		elapsed := time.Since(m.timer.LastStartedAt)
		if elapsed < m.timer.RemainingTime {
			m.timer.RemainingTime = m.timer.RemainingTime - elapsed
			m.timer.LastStartedAt = time.Now()
		} else {
			// Timer has expired
			m.timer.RemainingTime = 0
			m.timer.IsRunning = false
		}
	}

	return m.timer, nil
}

// UpdateTimer updates the timer state
func (m *MockTimerRepository) UpdateTimer(ctx context.Context, timer domain.Timer) (bool, error) {
	// Check if context is done
	if ctx.Err() != nil {
		return false, ctx.Err()
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	m.timer = timer
	return true, nil
}

// ResetTimer resets the timer with a new duration
func (m *MockTimerRepository) ResetTimer(ctx context.Context, duration time.Duration) (domain.Timer, error) {
	// Check if context is done
	if ctx.Err() != nil {
		return domain.Timer{}, ctx.Err()
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	m.timer.Duration = duration
	m.timer.RemainingTime = duration
	m.timer.IsRunning = false
	return m.timer, nil
}

// MockNoteRepository implements the NoteRepository interface with in-memory storage
type MockNoteRepository struct {
	notes map[int]domain.Note
	mu    sync.RWMutex
}

var _ repository.NoteRepository = &MockNoteRepository{}

// NewMockNoteRepository creates a new mock note repository
func NewMockNoteRepository() *MockNoteRepository {
	return &MockNoteRepository{
		notes: map[int]domain.Note{
			1: {
				ID:         1,
				Content:    "Introduction to the talk",
				PageNumber: 1,
				TotalPages: 10,
			},
			2: {
				ID:         2,
				Content:    "Key concepts and definitions",
				PageNumber: 2,
				TotalPages: 10,
			},
			// Add more sample notes as needed
		},
	}
}

// GetNote returns a note for a specific page
func (m *MockNoteRepository) GetNote(ctx context.Context, pageNumber int) (domain.Note, error) {
	// Check if context is done
	if ctx.Err() != nil {
		return domain.Note{}, ctx.Err()
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	if note, exists := m.notes[pageNumber]; exists {
		return note, nil
	}

	// Return an empty note with the correct page number if not found
	return domain.Note{
		PageNumber: pageNumber,
		TotalPages: len(m.notes),
	}, nil
}

// SaveNote saves a note
func (m *MockNoteRepository) SaveNote(ctx context.Context, note domain.Note) (bool, error) {
	// Check if context is done
	if ctx.Err() != nil {
		return false, ctx.Err()
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	m.notes[note.PageNumber] = note
	return true, nil
}

// MockQuestionRepository implements the QuestionRepository interface with in-memory storage
type MockQuestionRepository struct {
	questions []domain.Question
	mu        sync.RWMutex
	nextID    uint
}

var _ repository.QuestionRepository = &MockQuestionRepository{}

// NewMockQuestionRepository creates a new mock question repository
func NewMockQuestionRepository() *MockQuestionRepository {
	return &MockQuestionRepository{
		questions: make([]domain.Question, 0),
		nextID:    1,
	}
}

// GetQuestions returns all questions
func (m *MockQuestionRepository) GetQuestions(ctx context.Context) ([]domain.Question, error) {
	// Check if context is done
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	return m.questions, nil
}

// AddQuestion adds a new question
func (m *MockQuestionRepository) AddQuestion(ctx context.Context, question domain.Question) (domain.Question, error) {
	// Check if context is done
	if ctx.Err() != nil {
		return domain.Question{}, ctx.Err()
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	question.ID = m.nextID
	question.SubmittedAt = time.Now()
	question.Answered = false
	m.nextID++
	m.questions = append(m.questions, question)
	return question, nil
}

// MarkAsAnswered marks a question as answered
func (m *MockQuestionRepository) MarkAsAnswered(ctx context.Context, id uint) (bool, error) {
	// Check if context is done
	if ctx.Err() != nil {
		return false, ctx.Err()
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	for i, q := range m.questions {
		if q.ID == id {
			m.questions[i].Answered = true
			return true, nil
		}
	}
	return false, nil
}
