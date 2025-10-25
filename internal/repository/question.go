package repository

import (
	"survey-system/internal/model"

	"gorm.io/gorm"
)

// QuestionRepository defines the interface for question data operations
type QuestionRepository interface {
	Create(question *model.Question) error
	Update(question *model.Question) error
	Delete(id uint) error
	FindByID(id uint) (*model.Question, error)
	FindBySurveyID(surveyID uint) ([]model.Question, error)
	BatchUpdateOrder(questions []model.Question) error
}

// questionRepository implements QuestionRepository interface
type questionRepository struct {
	db *gorm.DB
}

// NewQuestionRepository creates a new question repository instance
func NewQuestionRepository(db *gorm.DB) QuestionRepository {
	return &questionRepository{db: db}
}

// Create creates a new question
func (r *questionRepository) Create(question *model.Question) error {
	return r.db.Create(question).Error
}

// Update updates an existing question
func (r *questionRepository) Update(question *model.Question) error {
	return r.db.Save(question).Error
}

// Delete deletes a question by ID
func (r *questionRepository) Delete(id uint) error {
	return r.db.Delete(&model.Question{}, id).Error
}

// FindByID finds a question by ID
func (r *questionRepository) FindByID(id uint) (*model.Question, error) {
	var question model.Question
	err := r.db.First(&question, id).Error
	if err != nil {
		return nil, err
	}
	return &question, nil
}

// FindBySurveyID finds all questions for a survey, ordered by the order field
func (r *questionRepository) FindBySurveyID(surveyID uint) ([]model.Question, error) {
	var questions []model.Question
	err := r.db.Where("survey_id = ?", surveyID).
		Order("\"order\" ASC").
		Find(&questions).Error
	if err != nil {
		return nil, err
	}
	return questions, nil
}

// BatchUpdateOrder updates the order field for multiple questions in a transaction
func (r *questionRepository) BatchUpdateOrder(questions []model.Question) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		for _, question := range questions {
			if err := tx.Model(&model.Question{}).
				Where("id = ?", question.ID).
				Update("order", question.Order).Error; err != nil {
				return err
			}
		}
		return nil
	})
}
