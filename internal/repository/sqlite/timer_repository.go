package sqlite

import (
	"context"
	"fmt"
	"time"

	"github.com/go-go-golems/ai-in-action-app/internal/domain"
	"github.com/go-go-golems/ai-in-action-app/internal/repository"
	"gorm.io/gorm"
)

// TimerRepository implements the repository.TimerRepository interface using GORM
type TimerRepository struct {
	db *gorm.DB
}

// Ensure TimerRepository implements repository.TimerRepository
var _ repository.TimerRepository = &TimerRepository{}

// NewTimerRepository creates a new timer repository
func NewTimerRepository(dbManager *DBManager) *TimerRepository {
	return &TimerRepository{
		db: dbManager.GetDB(),
	}
}

// GetTimer returns the current timer
func (r *TimerRepository) GetTimer(ctx context.Context) (domain.Timer, error) {
	// Check if context is done
	if ctx.Err() != nil {
		return domain.Timer{}, ctx.Err()
	}

	var model TimerModel
	result := r.db.WithContext(ctx).First(&model)

	// If no timer exists, create a default one
	if result.Error == gorm.ErrRecordNotFound {
		defaultTimer := domain.Timer{
			Duration:      15 * time.Minute,
			RemainingTime: 15 * time.Minute,
			IsRunning:     false,
		}

		// Save the default timer
		newTimer, err := r.createTimer(ctx, defaultTimer)
		if err != nil {
			return domain.Timer{}, fmt.Errorf("failed to create default timer: %w", err)
		}

		return newTimer, nil
	} else if result.Error != nil {
		return domain.Timer{}, fmt.Errorf("failed to get timer: %w", result.Error)
	}

	// If timer is running, calculate the remaining time
	timer := convertTimerModelToDomain(model)
	if timer.IsRunning {
		elapsed := time.Since(timer.LastStartedAt)
		if elapsed < timer.RemainingTime {
			timer.RemainingTime = timer.RemainingTime - elapsed

			// Update the timer with the new remaining time and last started time
			model.RemainingTime = timer.RemainingTime.Nanoseconds()
			model.LastStartedAt = time.Now()

			if err := r.db.WithContext(ctx).Save(&model).Error; err != nil {
				return domain.Timer{}, fmt.Errorf("failed to update timer: %w", err)
			}
		} else {
			// Timer has expired
			timer.RemainingTime = 0
			timer.IsRunning = false

			// Update the timer
			model.RemainingTime = 0
			model.IsRunning = false

			if err := r.db.WithContext(ctx).Save(&model).Error; err != nil {
				return domain.Timer{}, fmt.Errorf("failed to update timer: %w", err)
			}
		}
	}

	return timer, nil
}

// UpdateTimer updates the timer state
func (r *TimerRepository) UpdateTimer(ctx context.Context, timer domain.Timer) (bool, error) {
	// Check if context is done
	if ctx.Err() != nil {
		return false, ctx.Err()
	}

	// Convert domain entity to model
	model := convertDomainToTimerModel(timer)

	// Update in database
	result := r.db.WithContext(ctx).Save(&model)
	if result.Error != nil {
		return false, fmt.Errorf("failed to update timer: %w", result.Error)
	}

	// Check if any rows were affected
	return result.RowsAffected > 0, nil
}

// ResetTimer resets the timer with a new duration
func (r *TimerRepository) ResetTimer(ctx context.Context, duration time.Duration) (domain.Timer, error) {
	// Check if context is done
	if ctx.Err() != nil {
		return domain.Timer{}, ctx.Err()
	}

	var model TimerModel
	result := r.db.WithContext(ctx).First(&model)

	// If no timer exists, create a new one
	if result.Error == gorm.ErrRecordNotFound {
		newTimer := domain.Timer{
			Duration:      duration,
			RemainingTime: duration,
			IsRunning:     false,
		}

		return r.createTimer(ctx, newTimer)
	} else if result.Error != nil {
		return domain.Timer{}, fmt.Errorf("failed to get timer: %w", result.Error)
	}

	// Reset the timer
	model.Duration = duration.Nanoseconds()
	model.RemainingTime = duration.Nanoseconds()
	model.IsRunning = false

	// Save the updated timer
	if err := r.db.WithContext(ctx).Save(&model).Error; err != nil {
		return domain.Timer{}, fmt.Errorf("failed to reset timer: %w", err)
	}

	return convertTimerModelToDomain(model), nil
}

// createTimer creates a new timer
func (r *TimerRepository) createTimer(ctx context.Context, timer domain.Timer) (domain.Timer, error) {
	model := TimerModel{
		Duration:      timer.Duration.Nanoseconds(),
		RemainingTime: timer.RemainingTime.Nanoseconds(),
		IsRunning:     timer.IsRunning,
		LastStartedAt: time.Now(),
	}

	if err := r.db.WithContext(ctx).Create(&model).Error; err != nil {
		return domain.Timer{}, fmt.Errorf("failed to create timer: %w", err)
	}

	return convertTimerModelToDomain(model), nil
}

// Helper functions for conversion between domain and model

// convertTimerModelToDomain converts a TimerModel to a domain.Timer
func convertTimerModelToDomain(model TimerModel) domain.Timer {
	return domain.Timer{
		ID:            model.Model.ID,
		Duration:      time.Duration(model.Duration),
		RemainingTime: time.Duration(model.RemainingTime),
		IsRunning:     model.IsRunning,
		LastStartedAt: model.LastStartedAt,
	}
}

// convertDomainToTimerModel converts a domain.Timer to a TimerModel
func convertDomainToTimerModel(timer domain.Timer) TimerModel {
	return TimerModel{
		Model: gorm.Model{
			ID:        timer.ID,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		Duration:      timer.Duration.Nanoseconds(),
		RemainingTime: timer.RemainingTime.Nanoseconds(),
		IsRunning:     timer.IsRunning,
		LastStartedAt: timer.LastStartedAt,
	}
}
