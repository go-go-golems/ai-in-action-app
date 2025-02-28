package sqlite

import (
	"context"
	"fmt"
	"time"

	"github.com/go-go-golems/ai-in-action-app/internal/domain"
	"github.com/go-go-golems/ai-in-action-app/internal/repository"
	"gorm.io/gorm"
)

// NoteRepository implements the repository.NoteRepository interface using GORM
type NoteRepository struct {
	db *gorm.DB
}

// Ensure NoteRepository implements repository.NoteRepository
var _ repository.NoteRepository = &NoteRepository{}

// NewNoteRepository creates a new note repository
func NewNoteRepository(dbManager *DBManager) *NoteRepository {
	return &NoteRepository{
		db: dbManager.GetDB(),
	}
}

// GetNote returns a note for a specific page
func (r *NoteRepository) GetNote(ctx context.Context, pageNumber int) (domain.Note, error) {
	// Check if context is done
	if ctx.Err() != nil {
		return domain.Note{}, ctx.Err()
	}

	var model NoteModel
	result := r.db.WithContext(ctx).Where("page_number = ?", pageNumber).First(&model)

	// If no note exists for this page, create an empty one
	if result.Error == gorm.ErrRecordNotFound {
		// Get total pages
		var totalPages int
		if err := r.db.WithContext(ctx).Model(&NoteModel{}).Count(&totalPages).Error; err != nil {
			return domain.Note{}, fmt.Errorf("failed to count total pages: %w", err)
		}

		// Return an empty note with the correct page number
		return domain.Note{
			PageNumber: pageNumber,
			TotalPages: totalPages,
		}, nil
	} else if result.Error != nil {
		return domain.Note{}, fmt.Errorf("failed to get note: %w", result.Error)
	}

	return convertNoteModelToDomain(model), nil
}

// SaveNote saves a note
func (r *NoteRepository) SaveNote(ctx context.Context, note domain.Note) (bool, error) {
	// Check if context is done
	if ctx.Err() != nil {
		return false, ctx.Err()
	}

	var model NoteModel
	result := r.db.WithContext(ctx).Where("page_number = ?", note.PageNumber).First(&model)

	// If no note exists for this page, create a new one
	if result.Error == gorm.ErrRecordNotFound {
		model = NoteModel{
			Content:    note.Content,
			PageNumber: note.PageNumber,
			TotalPages: note.TotalPages,
		}

		if err := r.db.WithContext(ctx).Create(&model).Error; err != nil {
			return false, fmt.Errorf("failed to create note: %w", err)
		}

		return true, nil
	} else if result.Error != nil {
		return false, fmt.Errorf("failed to get note: %w", result.Error)
	}

	// Update existing note
	model.Content = note.Content
	model.TotalPages = note.TotalPages

	if err := r.db.WithContext(ctx).Save(&model).Error; err != nil {
		return false, fmt.Errorf("failed to update note: %w", err)
	}

	return true, nil
}

// Helper functions for conversion between domain and model

// convertNoteModelToDomain converts a NoteModel to a domain.Note
func convertNoteModelToDomain(model NoteModel) domain.Note {
	return domain.Note{
		ID:         model.Model.ID,
		Content:    model.Content,
		PageNumber: model.PageNumber,
		TotalPages: model.TotalPages,
	}
}

// convertDomainToNoteModel converts a domain.Note to a NoteModel
func convertDomainToNoteModel(note domain.Note) NoteModel {
	return NoteModel{
		Model: gorm.Model{
			ID:        note.ID,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		Content:    note.Content,
		PageNumber: note.PageNumber,
		TotalPages: note.TotalPages,
	}
}
