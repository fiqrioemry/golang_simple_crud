package repositories

import (
	"server/internal/models"

	"gorm.io/gorm"
)

type FormRepository interface {
	Create(form *models.Form) error
	FindAllByUserID(userID string) ([]models.Form, error)
	FindByID(id string) (*models.Form, error)
	GetFormSetting(formID string) (*models.FormSetting, error)
	UpdateFormSetting(setting *models.FormSetting) error
	AddSection(section *models.FormSection) error
	GetQuestionsByFormID(formID string) ([]models.Question, error)
	DeleteSection(sectionID string) error
	UpdateSection(section *models.FormSection) error
	GetSectionsByFormID(formID string) ([]models.FormSection, error)
	AddQuestion(q *models.Question) error
	UpdateQuestion(q *models.Question) error
	DeleteQuestion(id string) error
}

type formRepository struct {
	db *gorm.DB
}

func NewFormRepository(db *gorm.DB) FormRepository {
	return &formRepository{db}
}

func (r *formRepository) Create(form *models.Form) error {
	return r.db.Create(form).Error
}

func (r *formRepository) FindAllByUserID(userID string) ([]models.Form, error) {
	var forms []models.Form
	err := r.db.Where("user_id = ?", userID).Order("created_at desc").Find(&forms).Error
	return forms, err
}

func (r *formRepository) FindByID(id string) (*models.Form, error) {
	var form models.Form
	err := r.db.First(&form, "id = ?", id).Error
	return &form, err
}

func (r *formRepository) GetFormSetting(formID string) (*models.FormSetting, error) {
	var setting models.FormSetting
	err := r.db.Where("form_id = ?", formID).First(&setting).Error
	return &setting, err
}

func (r *formRepository) UpdateFormSetting(setting *models.FormSetting) error {
	return r.db.Model(&models.FormSetting{}).
		Where("form_id = ?", setting.FormID).
		Updates(setting).Error
}

func (r *formRepository) AddSection(section *models.FormSection) error {
	return r.db.Create(section).Error
}

func (r *formRepository) GetSectionsByFormID(formID string) ([]models.FormSection, error) {
	var sections []models.FormSection
	err := r.db.Where("form_id = ?", formID).Order("`order` asc").Find(&sections).Error
	return sections, err
}

func (r *formRepository) UpdateSection(section *models.FormSection) error {
	return r.db.Model(&models.FormSection{}).
		Where("id = ?", section.ID).
		Updates(section).Error
}

func (r *formRepository) DeleteSection(sectionID string) error {
	return r.db.Delete(&models.FormSection{}, "id = ?", sectionID).Error
}

func (r *formRepository) GetQuestionsByFormID(formID string) ([]models.Question, error) {
	var questions []models.Question
	err := r.db.Preload("Options").
		Where("form_id = ?", formID).
		Order("`order` asc").
		Find(&questions).Error
	return questions, err
}

func (r *formRepository) AddQuestion(q *models.Question) error {
	return r.db.Create(q).Error
}

func (r *formRepository) UpdateQuestion(q *models.Question) error {
	return r.db.Model(&models.Question{}).Where("id = ?", q.ID).Updates(q).Error
}

func (r *formRepository) DeleteQuestion(id string) error {
	return r.db.Delete(&models.Question{}, "id = ?", id).Error
}
