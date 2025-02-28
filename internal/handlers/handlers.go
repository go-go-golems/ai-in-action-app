package handlers

import (
	"github.com/go-go-golems/ai-in-action-app/internal/repository"
	"github.com/labstack/echo/v4"
)

// RegisterHandlers registers all handlers with the Echo instance
func RegisterHandlers(e *echo.Echo, eventRepo repository.EventRepository, timerRepo repository.TimerRepository, noteRepo repository.NoteRepository, questionRepo repository.QuestionRepository) {
	// Register event handlers
	eventHandler := NewEventHandler(eventRepo)
	eventHandler.RegisterRoutes(e)

	// TODO: Register timer handlers
	// TODO: Register note handlers
	// TODO: Register question handlers
}
