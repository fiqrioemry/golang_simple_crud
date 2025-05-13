package services

import (
	"server/internal/dto"
	"server/internal/models"
	"server/internal/repositories"
	"time"

	"github.com/google/uuid"
)

type FormService interface {
	CreateForm(userID string, req *dto.CreateFormRequest) error
	GetAllForms(userID string) ([]dto.FormResponse, error)
	GetFormDetail(formID string) (*dto.FormDetailResponse, error)
	AddFormSection(req *dto.AddSectionRequest) (*dto.SectionResponse, error)
	UpdateFormSettings(formID string, req *dto.UpdateFormSettingRequest) error
	GetFormSettings(formID string) (*dto.FormSettingResponse, error)
	GetFormQuestion(formID string) ([]dto.QuestionResponse, error)
	DeleteFormSections(sectionID string) error
	UpdateFormSections(req *dto.UpdateSectionRequest) error
	GetFormSections(formID string) ([]dto.SectionResponse, error)
	DeleteQuestion(id string) error

	AddFormQuestion(req *dto.AddQuestionRequest) error
	UpdateQuestion(req *dto.UpdateQuestionRequest) error
}

type formService struct {
	repo repositories.FormRepository
}

func NewFormService(repo repositories.FormRepository) FormService {
	return &formService{repo}
}

func (s *formService) CreateForm(userID string, req *dto.CreateFormRequest) error {
	form := &models.Form{
		ID:          uuid.New(),
		UserID:      uuid.MustParse(userID),
		Title:       req.Title,
		Description: req.Description,
		Type:        req.Type,
		IsActive:    true,
		Duration:    req.Duration,
	}
	return s.repo.Create(form)
}

func (s *formService) GetAllForms(userID string) ([]dto.FormResponse, error) {
	forms, err := s.repo.FindAllByUserID(userID)
	if err != nil {
		return nil, err
	}
	var result []dto.FormResponse
	for _, f := range forms {
		result = append(result, dto.FormResponse{
			ID:          f.ID.String(),
			Title:       f.Title,
			Description: f.Description,
			Type:        f.Type,
			IsActive:    f.IsActive,
			Duration:    f.Duration,
			CreatedAt:   f.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}
	return result, nil
}

func (s *formService) GetFormDetail(formID string) (*dto.FormDetailResponse, error) {
	form, err := s.repo.FindByID(formID)
	if err != nil {
		return nil, err
	}
	return &dto.FormDetailResponse{
		ID:          form.ID.String(),
		Title:       form.Title,
		Description: form.Description,
		Type:        form.Type,
		IsActive:    form.IsActive,
		Duration:    form.Duration,
		CreatedAt:   form.CreatedAt.Format("2006-01-02 15:04:05"),
	}, nil
}

func (s *formService) GetFormSettings(formID string) (*dto.FormSettingResponse, error) {
	setting, err := s.repo.GetFormSetting(formID)
	if err != nil {
		return nil, err
	}

	var start, end *string
	if setting.StartAt != nil {
		s := setting.StartAt.Format(time.RFC3339)
		start = &s
	}
	if setting.EndAt != nil {
		e := setting.EndAt.Format(time.RFC3339)
		end = &e
	}

	return &dto.FormSettingResponse{
		FormID:             setting.FormID.String(),
		ShowResult:         setting.ShowResult,
		MultipleSubmission: setting.MultipleSubmission,
		PassingGrade:       setting.PassingGrade,
		Grading:            setting.Grading,
		MaxSubmissions:     setting.MaxSubmissions,
		StartAt:            start,
		EndAt:              end,
	}, nil
}

func (s *formService) UpdateFormSettings(formID string, req *dto.UpdateFormSettingRequest) error {
	startAt, endAt := parseTimePointer(req.StartAt), parseTimePointer(req.EndAt)

	setting := &models.FormSetting{
		FormID:             uuid.MustParse(formID),
		ShowResult:         req.ShowResult,
		MultipleSubmission: req.MultipleSubmission,
		PassingGrade:       req.PassingGrade,
		Grading:            req.Grading,
		MaxSubmissions:     req.MaxSubmissions,
		StartAt:            startAt,
		EndAt:              endAt,
	}
	return s.repo.UpdateFormSetting(setting)
}

func (s *formService) AddFormSection(req *dto.AddSectionRequest) (*dto.SectionResponse, error) {
	section := &models.FormSection{
		ID:          uuid.New(),
		FormID:      uuid.MustParse(req.FormID),
		Title:       req.Title,
		Description: req.Description,
		Order:       req.Order,
	}
	if err := s.repo.AddSection(section); err != nil {
		return nil, err
	}
	return &dto.SectionResponse{
		ID:          section.ID.String(),
		FormID:      section.FormID.String(),
		Title:       section.Title,
		Description: section.Description,
		Order:       section.Order,
	}, nil
}

func parseTimePointer(input *string) *time.Time {
	if input == nil {
		return nil
	}
	t, err := time.Parse(time.RFC3339, *input)
	if err != nil {
		return nil
	}
	return &t
}

func (s *formService) GetFormSections(formID string) ([]dto.SectionResponse, error) {
	sections, err := s.repo.GetSectionsByFormID(formID)
	if err != nil {
		return nil, err
	}
	var result []dto.SectionResponse
	for _, s := range sections {
		result = append(result, dto.SectionResponse{
			ID:          s.ID.String(),
			FormID:      s.FormID.String(),
			Title:       s.Title,
			Description: s.Description,
			Order:       s.Order,
		})
	}
	return result, nil
}

func (s *formService) UpdateFormSections(req *dto.UpdateSectionRequest) error {
	section := &models.FormSection{
		ID:          uuid.MustParse(req.ID),
		Title:       req.Title,
		Description: req.Description,
		Order:       req.Order,
	}
	return s.repo.UpdateSection(section)
}

func (s *formService) DeleteFormSections(sectionID string) error {
	return s.repo.DeleteSection(sectionID)
}

func (s *formService) GetFormQuestion(formID string) ([]dto.QuestionResponse, error) {
	questions, err := s.repo.GetQuestionsByFormID(formID)
	if err != nil {
		return nil, err
	}

	var result []dto.QuestionResponse
	for _, q := range questions {
		var opts []dto.Option
		for _, o := range q.Options {
			opts = append(opts, dto.Option{
				ID:        o.ID,
				Text:      o.Text,
				ImageURL:  o.ImageURL,
				IsCorrect: o.IsCorrect,
			})
		}
		result = append(result, dto.QuestionResponse{
			ID:         q.ID.String(),
			Text:       q.Text,
			Type:       q.Type,
			IsRequired: q.IsRequired,
			Order:      q.Order,
			Score:      q.Score,
			ImageURL:   q.ImageURL,
			Options:    opts,
		})
	}

	return result, nil
}

func (s *formService) AddFormQuestion(req *dto.AddQuestionRequest) error {

	sectionID := uuid.MustParse(req.SectionID)
	q := &models.Question{
		ID:         uuid.New(),
		FormID:     uuid.MustParse(req.FormID),
		SectionID:  &sectionID,
		Text:       req.Text,
		Type:       req.Type,
		IsRequired: req.IsRequired,
		Order:      req.Order,
		Score:      req.Score,
		ImageURL:   req.ImageURL,
	}
	return s.repo.AddQuestion(q)
}

func (s *formService) UpdateQuestion(req *dto.UpdateQuestionRequest) error {
	q := &models.Question{
		ID:         uuid.MustParse(req.ID),
		Text:       req.Text,
		Type:       req.Type,
		IsRequired: req.IsRequired,
		Order:      req.Order,
		Score:      req.Score,
		ImageURL:   req.ImageURL,
	}
	return s.repo.UpdateQuestion(q)
}

func (s *formService) DeleteQuestion(id string) error {
	return s.repo.DeleteQuestion(id)
}
