package repositories

import (
	"server/internal/models"

	"gorm.io/gorm"
)

type SubmissionRepository interface {
	Create(sub *models.Submission, answers []models.Answer) error
	GetByFormID(formID string) ([]models.Submission, error)
	GetWithAnswers(subID string) (*models.Submission, error)
}

type submissionRepository struct {
	db *gorm.DB
}

func NewSubmissionRepository(db *gorm.DB) SubmissionRepository {
	return &submissionRepository{db}
}

func (r *submissionRepository) Create(sub *models.Submission, answers []models.Answer) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(sub).Error; err != nil {
			return err
		}
		for i := range answers {
			answers[i].SubmissionID = sub.ID
		}
		return tx.Create(&answers).Error
	})
}

func (r *submissionRepository) GetByFormID(formID string) ([]models.Submission, error) {
	var subs []models.Submission
	err := r.db.Where("form_id = ?", formID).Order("submitted_at desc").Find(&subs).Error
	return subs, err
}

func (r *submissionRepository) GetWithAnswers(subID string) (*models.Submission, error) {
	var sub models.Submission
	err := r.db.Preload("Answers").First(&sub, "id = ?", subID).Error
	return &sub, err
}
