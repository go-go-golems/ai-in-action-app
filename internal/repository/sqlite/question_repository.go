package sqlite

import (
	"context"
	"fmt"
	"time"

	"github.com/go-go-golems/ai-in-action-app/internal/domain"
	"github.com/go-go-golems/ai-in-action-app/internal/repository"
	"gorm.io/gorm"
)

// QuestionRepository implements the repository.QuestionRepository interface using GORM
type QuestionRepository struct {
	db *gorm.DB
}

// Ensure QuestionRepository implements repository.QuestionRepository
var _ repository.QuestionRepository = &QuestionRepository{}

// NewQuestionRepository creates a new question repository
func NewQuestionRepository(dbManager *DBManager) *QuestionRepository {
	return &QuestionRepository{
		db: dbManager.GetDB(),
	}
}

// GetQuestions returns all questions
func (r *QuestionRepository) GetQuestions(ctx context.Context) ([]domain.Question, error) {
	// Check if context is done
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	var models []QuestionModel
	if err := r.db.WithContext(ctx).Order("submitted_at desc").Find(&models).Error; err != nil {
		return nil, fmt.Errorf("failed to get questions: %w", err)
	}

	// Convert models to domain entities
	questions := make([]domain.Question, len(models))
	for i, model := range models {
		questions[i] = convertQuestionModelToDomain(model)
	}

	return questions, nil
}

// AddQuestion adds a new question
func (r *QuestionRepository) AddQuestion(ctx context.Context, question domain.Question) (domain.Question, error) {
	// Check if context is done
	if ctx.Err() != nil {
		return domain.Question{}, ctx.Err()
	}

	// Set submission time if not already set
	if question.SubmittedAt.IsZero() {
		question.SubmittedAt = time.Now()
	}
	
	// Ensure question is not marked as answered when first added
	question.Answered = false

	// Convert domain entity to model
	model := convertDomainToQuestionModel(question)

	// Save to database
	if err := r.db.WithContext(ctx).Create(&model).Error; err != nil {
		return domain.Question{}, fmt.Errorf("failed to add question: %w", err)
	}

	// Return the question with the new ID
	return convertQuestionModelToDomain(model), nil
}

// MarkAsAnswered marks a question as answered
func (r *QuestionRepository) MarkAsAnswered(ctx context.Context, id uint) (bool, error) {
	// Check if context is done
	if ctx.Err() != nil {
		return false, ctx.Err()
	}

	// Update the question
	result := r.db.WithContext(ctx).Model(&QuestionModel{}).Where("id = ?", id).Update("answered", true)
	if result.Error != nil {
		return false, fmt.Errorf("failed to mark question as answered: %w", result.Error)
	}

	// Check if any rows were affected
	return result.RowsAffected > 0, nil
}

// Helper functions for conversion between domain and model

// convertQuestionModelToDomain converts a QuestionModel to a domain.Question
func convertQuestionModelToDomain(model QuestionModel) domain.Question {
	return domain.Question{
		ID:          model.Model.ID,
		Name:        model.Name,
		Content:     model.Content,
		SubmittedAt: model.SubmittedAt,
		Answered:    model.Answered,
	}
}

// convertDomainToQuestionModel converts a domain.Question to a QuestionModel
func convertDomainToQuestionModel(question domain.Question) QuestionModel {
	return QuestionModel{
		Model: gorm.Model{
			ID:        question.ID,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		Name:        question.Name,
		Content:     question.Content,
		SubmittedAt: question.SubmittedAt,
		Answered:    question.Answered,
	}
} 