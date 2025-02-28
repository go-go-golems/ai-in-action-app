package sqlite

import (
	"context"
	"fmt"
	"time"

	"github.com/go-go-golems/ai-in-action-app/internal/domain"
	"github.com/go-go-golems/ai-in-action-app/internal/repository"
	"gorm.io/gorm"
)

// EventRepository implements the repository.EventRepository interface using GORM
type EventRepository struct {
	db *gorm.DB
}

// Ensure EventRepository implements repository.EventRepository
var _ repository.EventRepository = &EventRepository{}

// NewEventRepository creates a new event repository
func NewEventRepository(dbManager *DBManager) *EventRepository {
	return &EventRepository{
		db: dbManager.GetDB(),
	}
}

// GetUpcomingEvents returns all upcoming events
func (r *EventRepository) GetUpcomingEvents(ctx context.Context) ([]domain.Event, error) {
	// Check if context is done
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	var models []EventModel
	if err := r.db.WithContext(ctx).Where("is_upcoming = ?", true).Find(&models).Error; err != nil {
		return nil, fmt.Errorf("failed to get upcoming events: %w", err)
	}

	// Convert models to domain entities
	events := make([]domain.Event, len(models))
	for i, model := range models {
		events[i] = convertEventModelToDomain(model)
	}

	return events, nil
}

// GetPastEvents returns all past events
func (r *EventRepository) GetPastEvents(ctx context.Context) ([]domain.Event, error) {
	// Check if context is done
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	var models []EventModel
	if err := r.db.WithContext(ctx).Where("is_upcoming = ?", false).Find(&models).Error; err != nil {
		return nil, fmt.Errorf("failed to get past events: %w", err)
	}

	// Convert models to domain entities
	events := make([]domain.Event, len(models))
	for i, model := range models {
		events[i] = convertEventModelToDomain(model)
	}

	return events, nil
}

// AddEvent adds a new event and returns it with an ID
func (r *EventRepository) AddEvent(ctx context.Context, event domain.Event) (domain.Event, error) {
	// Check if context is done
	if ctx.Err() != nil {
		return domain.Event{}, ctx.Err()
	}

	// Convert domain entity to model
	model := convertDomainToEventModel(event)

	// Save to database
	if err := r.db.WithContext(ctx).Create(&model).Error; err != nil {
		return domain.Event{}, fmt.Errorf("failed to add event: %w", err)
	}

	// Return the event with the new ID
	return convertEventModelToDomain(model), nil
}

// UpdateEvent updates an existing event
func (r *EventRepository) UpdateEvent(ctx context.Context, event domain.Event) (bool, error) {
	// Check if context is done
	if ctx.Err() != nil {
		return false, ctx.Err()
	}

	// Convert domain entity to model
	model := convertDomainToEventModel(event)

	// Update in database
	result := r.db.WithContext(ctx).Save(&model)
	if result.Error != nil {
		return false, fmt.Errorf("failed to update event: %w", result.Error)
	}

	// Check if any rows were affected
	return result.RowsAffected > 0, nil
}

// Helper functions for conversion between domain and model

// convertEventModelToDomain converts an EventModel to a domain.Event
func convertEventModelToDomain(model EventModel) domain.Event {
	return domain.Event{
		ID:          model.Model.ID,
		Title:       model.Title,
		Speaker:     model.Speaker,
		Description: model.Description,
		Date:        model.Date,
		IsUpcoming:  model.IsUpcoming,
	}
}

// convertDomainToEventModel converts a domain.Event to an EventModel
func convertDomainToEventModel(event domain.Event) EventModel {
	return EventModel{
		Model: gorm.Model{
			ID:        event.ID,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		Title:       event.Title,
		Speaker:     event.Speaker,
		Description: event.Description,
		Date:        event.Date,
		IsUpcoming:  event.IsUpcoming,
	}
}
