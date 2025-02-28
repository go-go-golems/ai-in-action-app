package sqlite

import (
	"github.com/go-go-golems/ai-in-action-app/internal/repository"
)

// RepositoryFactory creates and manages all SQLite repositories
type RepositoryFactory struct {
	dbManager          *DBManager
	eventRepository    *EventRepository
	timerRepository    *TimerRepository
	noteRepository     *NoteRepository
	questionRepository *QuestionRepository
}

// NewRepositoryFactory creates a new repository factory
func NewRepositoryFactory(dbPath string) (*RepositoryFactory, error) {
	// Create database manager
	dbManager, err := NewDBManager(dbPath)
	if err != nil {
		return nil, err
	}

	// Create repositories
	factory := &RepositoryFactory{
		dbManager:          dbManager,
		eventRepository:    NewEventRepository(dbManager),
		timerRepository:    NewTimerRepository(dbManager),
		noteRepository:     NewNoteRepository(dbManager),
		questionRepository: NewQuestionRepository(dbManager),
	}

	return factory, nil
}

// GetEventRepository returns the event repository
func (f *RepositoryFactory) GetEventRepository() repository.EventRepository {
	return f.eventRepository
}

// GetTimerRepository returns the timer repository
func (f *RepositoryFactory) GetTimerRepository() repository.TimerRepository {
	return f.timerRepository
}

// GetNoteRepository returns the note repository
func (f *RepositoryFactory) GetNoteRepository() repository.NoteRepository {
	return f.noteRepository
}

// GetQuestionRepository returns the question repository
func (f *RepositoryFactory) GetQuestionRepository() repository.QuestionRepository {
	return f.questionRepository
}

// Close closes all repositories and the database connection
func (f *RepositoryFactory) Close() error {
	return f.dbManager.Close()
}
