package dto

// AUTHENTICATION
type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=5"`
	Fullname string `json:"fullname" binding:"required,min=5"`
}

type AuthResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=5"`
}

type SendOTPRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type VerifyOTPRequest struct {
	Email string `json:"email" binding:"required,email"`
	OTP   string `json:"otp" binding:"required,len=6"`
}

type UserInfoResponse struct {
	UserID   string `json:"userId"`
	Email    string `json:"email"`
	Fullname string `json:"fullname"`
	Avatar   string `json:"avatar"`
}

// sSUBSCRIPTIONS
type CreateTierRequest struct {
	Name        string  `json:"name" binding:"required,min=3"`
	TokenLimit  int     `json:"tokenLimit" binding:"required,gte=1"`
	Price       float64 `json:"price" binding:"required,gte=0"`
	Duration    int     `json:"duration" binding:"required,gte=1"`
	Description string  `json:"description"`
}

type UpdateTierRequest struct {
	ID          uint    `json:"id" binding:"required"`
	Name        string  `json:"name" binding:"required,min=3"`
	TokenLimit  int     `json:"tokenLimit" binding:"required,gte=1"`
	Price       float64 `json:"price" binding:"required,gte=0"`
	Duration    int     `json:"duration" binding:"required,gte=1"`
	Description string  `json:"description"`
}

type TierResponse struct {
	ID          uint    `json:"id"`
	Name        string  `json:"name"`
	TokenLimit  int     `json:"tokenLimit"`
	Price       float64 `json:"price"`
	Duration    int     `json:"duration"`
	Description string  `json:"description"`
}

type UserSubscriptionResponse struct {
	UserID    string `json:"userId"`
	Email     string `json:"email"`
	Fullname  string `json:"fullname"`
	Avatar    string `json:"avatar"`
	TierName  string `json:"tierName"`
	IsActive  bool   `json:"isActive"`
	ExpiresAt string `json:"expiresAt"`
	Remaining int    `json:"remainingTokens"`
}

// PAYMENT
type CreatePaymentRequest struct {
	TierID      uint    `json:"tierId" binding:"required"`
	Method      string  `json:"method" binding:"required"`
	VoucherCode *string `json:"voucherCode"`
}

type CreatePaymentResponse struct {
	PaymentID string `json:"paymentId"`
	SnapToken string `json:"snapToken"`
	SnapURL   string `json:"snapUrl"`
}

type MidtransNotificationRequest struct {
	OrderID           string `json:"order_id"`
	TransactionStatus string `json:"transaction_status"`
	PaymentType       string `json:"payment_type"`
	FraudStatus       string `json:"fraud_status"`
}

type PaymentDetailResponse struct {
	ID       string  `json:"id"`
	UserID   string  `json:"userId"`
	TierID   uint    `json:"tierId"`
	Subtotal float64 `json:"subtotal"`
	Tax      float64 `json:"tax"`
	Total    float64 `json:"total"`
	Method   string  `json:"method"`
	Status   string  `json:"status"`
	PaidAt   string  `json:"paidAt"`
}

type PaymentResponse struct {
	ID            string  `json:"id"`
	UserID        string  `json:"userId"`
	UserEmail     string  `json:"userEmail"`
	Fullname      string  `json:"fullname"`
	TierID        uint    `json:"tierId"`
	TierName      string  `json:"tierName"`
	Subtotal      float64 `json:"subtotal"`
	Tax           float64 `json:"tax"`
	Total         float64 `json:"total"`
	PaymentMethod string  `json:"method"`
	Status        string  `json:"status"`
	PaidAt        string  `json:"paidAt"`
}

type PaymentListResponse struct {
	Payments []PaymentResponse `json:"payments"`
	Total    int64             `json:"total"`
	Page     int               `json:"page"`
	Limit    int               `json:"limit"`
}

// USER
type UpdateProfileRequest struct {
	Fullname string `json:"fullname" binding:"required,min=3"`
	Avatar   string `json:"avatar"`
}

type UserProfileResponse struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Fullname string `json:"fullname"`
	Avatar   string `json:"avatar"`
}

type MyFormResponse struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Type  string `json:"type"`
}

type MyFormDetailResponse struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Type        string `json:"type"`
	CreatedAt   string `json:"createdAt"`
}

// FORM
type CreateFormRequest struct {
	Title       string `json:"title" binding:"required,min=3"`
	Description string `json:"description"`
	Type        string `json:"type" binding:"required,oneof=quiz exam survey quisoner diagnose"`
	Duration    *int   `json:"duration"` // opsional untuk quiz dan exam
}

type FormResponse struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Type        string `json:"type"`
	IsActive    bool   `json:"isActive"`
	Duration    *int   `json:"duration"`
	CreatedAt   string `json:"createdAt"`
}

type FormDetailResponse struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Type        string `json:"type"`
	IsActive    bool   `json:"isActive"`
	Duration    *int   `json:"duration"`
	CreatedAt   string `json:"createdAt"`
}

type FormSettingResponse struct {
	FormID             string   `json:"formId"`
	ShowResult         bool     `json:"showResult"`
	MultipleSubmission bool     `json:"multipleSubmission"`
	PassingGrade       *float64 `json:"passingGrade"`
	Grading            bool     `json:"grading"`
	MaxSubmissions     *int     `json:"maxSubmissions"`
	StartAt            *string  `json:"startAt"`
	EndAt              *string  `json:"endAt"`
}

type UpdateFormSettingRequest struct {
	ShowResult         bool     `json:"showResult"`
	MultipleSubmission bool     `json:"multipleSubmission"`
	PassingGrade       *float64 `json:"passingGrade"`
	Grading            bool     `json:"grading"`
	MaxSubmissions     *int     `json:"maxSubmissions"`
	StartAt            *string  `json:"startAt"` // ISO 8601 format
	EndAt              *string  `json:"endAt"`   // ISO 8601 format
}

type AddSectionRequest struct {
	FormID      string `json:"formId" binding:"required"`
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	Order       int    `json:"order"`
}

type SectionResponse struct {
	ID          string `json:"id"`
	FormID      string `json:"formId"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Order       int    `json:"order"`
}

type UpdateSectionRequest struct {
	ID          string `json:"id" binding:"required"`
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	Order       int    `json:"order"`
}

type QuestionResponse struct {
	ID         string   `json:"id"`
	Text       string   `json:"text"`
	Type       string   `json:"type"`
	IsRequired bool     `json:"isRequired"`
	Order      int      `json:"order"`
	Score      *int     `json:"score"`
	ImageURL   *string  `json:"imageUrl"`
	Options    []Option `json:"options"`
}

type Option struct {
	ID        uint    `json:"id"`
	Text      string  `json:"text"`
	ImageURL  *string `json:"imageUrl"`
	IsCorrect *bool   `json:"isCorrect"`
}

type AddQuestionRequest struct {
	FormID     string  `json:"formId" binding:"required"`
	SectionID  string  `json:"sectionId" binding:"required"`
	Text       string  `json:"text" binding:"required"`
	Type       string  `json:"type" binding:"required"` // text, radio, checkbox, etc.
	IsRequired bool    `json:"isRequired"`
	Order      int     `json:"order"`
	Score      *int    `json:"score"`
	ImageURL   *string `json:"imageUrl"`
}

type UpdateQuestionRequest struct {
	ID         string  `json:"id" binding:"required"`
	Text       string  `json:"text" binding:"required"`
	Type       string  `json:"type" binding:"required"`
	IsRequired bool    `json:"isRequired"`
	Order      int     `json:"order"`
	Score      *int    `json:"score"`
	ImageURL   *string `json:"imageUrl"`
}

// SUBMISSIONS
type AnswerRequest struct {
	QuestionID string  `json:"questionId" binding:"required"`
	OptionID   *uint   `json:"optionId,omitempty"`
	TextAnswer *string `json:"textAnswer,omitempty"`
}

type SubmissionRequest struct {
	FormID       string          `json:"formId" binding:"required"`
	Email        string          `json:"email"`
	IPAddress    *string         `json:"ipAddress"`
	UserAgent    *string         `json:"userAgent"`
	SessionToken *string         `json:"sessionToken"`
	Answers      []AnswerRequest `json:"answers" binding:"required,min=1"`
}

type SubmissionResponse struct {
	ID        string   `json:"id"`
	FormID    string   `json:"formId"`
	Email     string   `json:"email"`
	Score     *float64 `json:"score"`
	Timestamp string   `json:"submittedAt"`
}

type SubmissionResultResponse struct {
	FormTitle  string           `json:"formTitle"`
	TotalScore *float64         `json:"totalScore,omitempty"`
	Answers    []AnswerResponse `json:"answers"`
}

type AnswerResponse struct {
	Question string `json:"question"`
	Answer   string `json:"answer"`
	Correct  *bool  `json:"correct,omitempty"` // jika quiz atau exam
}
