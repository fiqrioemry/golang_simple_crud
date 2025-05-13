type SubmissionService interface {
	SendSubmission(req *dto.SubmissionRequest) error
	GetFormSubmissions(formID string) ([]dto.SubmissionResponse, error)
	GetSubmissionResult(subID string) (*dto.SubmissionResultResponse, error)
}

type submissionService struct {
	repo     repositories.SubmissionRepository
	formRepo repositories.FormRepository
}

func NewSubmissionService(repo repositories.SubmissionRepository, formRepo repositories.FormRepository) SubmissionService {
	return &submissionService{repo, formRepo}
}

func (s *submissionService) SendSubmission(req *dto.SubmissionRequest) error {
	sub := &models.Submission{
		ID:           uuid.New(),
		FormID:       uuid.MustParse(req.FormID),
		Email:        req.Email,
		IPAddress:    req.IPAddress,
		UserAgent:    req.UserAgent,
		SessionToken: req.SessionToken,
		SubmittedAt:  time.Now(),
	}

	var answers []models.Answer
	for _, a := range req.Answers {
		ans := models.Answer{
			QuestionID: uuid.MustParse(a.QuestionID),
			OptionID:   a.OptionID,
			TextAnswer: a.TextAnswer,
		}
		answers = append(answers, ans)
	}

	return s.repo.Create(sub, answers)
}

func (s *submissionService) GetFormSubmissions(formID string) ([]dto.SubmissionResponse, error) {
	data, err := s.repo.GetByFormID(formID)
	if err != nil {
		return nil, err
	}
	var result []dto.SubmissionResponse
	for _, d := range data {
		result = append(result, dto.SubmissionResponse{
			ID:        d.ID.String(),
			FormID:    d.FormID.String(),
			Email:     d.Email,
			Score:     d.Score,
			Timestamp: d.SubmittedAt.Format("2006-01-02 15:04:05"),
		})
	}
	return result, nil
}

func (s *submissionService) GetSubmissionResult(subID string) (*dto.SubmissionResultResponse, error) {
	sub, err := s.repo.GetWithAnswers(subID)
	if err != nil {
		return nil, err
	}

	form, err := s.formRepo.FindByID(sub.FormID.String())
	if err != nil {
		return nil, err
	}

	var answers []dto.AnswerResponse
	for _, a := range sub.Answers {
		answerText := ""
		if a.TextAnswer != nil {
			answerText = *a.TextAnswer
		} else if a.OptionID != nil {
			answerText = fmt.Sprintf("Option #%d", *a.OptionID)
		}
		answers = append(answers, dto.AnswerResponse{
			Question: a.QuestionID.String(), // idealnya preload pertanyaan
			Answer:   answerText,
		})
	}

	return &dto.SubmissionResultResponse{
		FormTitle:  form.Title,
		TotalScore: sub.Score,
		Answers:    answers,
	}, nil
}
