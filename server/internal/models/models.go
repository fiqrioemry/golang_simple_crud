package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID        uuid.UUID `gorm:"type:char(36);primaryKey" json:"id"`
	Email     string    `gorm:"type:varchar(255);uniqueIndex;not null" json:"email"`
	Password  string    `gorm:"type:text;not null" json:"-"`
	Role      string    `gorm:"type:text;not null" json:"-"`
	Fullname  string    `gorm:"type:varchar(255);not null" json:"fullname"`
	Avatar    string    `gorm:"type:varchar(255);not null" json:"avatar"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`

	Forms []Form `gorm:"foreignKey:UserID"`
}

type Token struct {
	ID        uuid.UUID      `gorm:"type:char(36);primaryKey" json:"id"`
	UserID    uuid.UUID      `gorm:"type:char(36);index;not null" json:"userId"`
	Token     string         `gorm:"type:text;not null" json:"token"`
	ExpiredAt time.Time      `json:"expiredAt"`
	CreatedAt time.Time      `gorm:"autoCreateTime"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type UserSubscription struct {
	ID                 uint      `gorm:"primaryKey"`
	UserID             uuid.UUID `gorm:"type:uuid;not null;uniqueIndex"`
	SubscriptionTierID uint      `gorm:"not null"`
	StartedAt          time.Time `gorm:"not null"`
	ExpiresAt          time.Time `gorm:"not null"`
	IsActive           bool      `gorm:"default:true"`
	RemainingTokens    int       `gorm:"not null;default:0"`

	User             User             `gorm:"foreignKey:UserID"`
	SubscriptionTier SubscriptionTier `gorm:"foreignKey:SubscriptionTierID"`
}

type SubscriptionTier struct {
	ID          uint    `gorm:"primaryKey"`
	Name        string  `gorm:"type:varchar(50);uniqueIndex;not null"`
	TokenLimit  int     `gorm:"not null"`
	Price       float64 `gorm:"not null"`
	Duration    int     `gorm:"not null"`
	Description string  `gorm:"type:text"`
}

type Payment struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;index"`
	TierID    uint      `gorm:"not null"`
	Subtotal  float64   `gorm:"not null"`
	Tax       float64   `gorm:"not null"`
	Total     float64   `gorm:"not null"`
	Method    string    `gorm:"type:varchar(50);not null"`
	Status    string    `gorm:"type:varchar(20);default:'pending';check:status IN ('pending','paid','failed');not null"`
	PaidAt    *time.Time
	CreatedAt time.Time

	User User
	Tier SubscriptionTier `gorm:"foreignKey:TierID"`
}

type Form struct {
	ID          uuid.UUID `gorm:"type:char(36);primaryKey"`
	UserID      uuid.UUID `gorm:"type:char(36);not null;index"`
	Title       string    `gorm:"type:varchar(255);not null"`
	Description string    `gorm:"type:text"`
	Type        string    `gorm:"type:varchar(50);not null"` // quiz, exam, survey, quisoner, diagnose form
	IsActive    bool      `gorm:"default:true"`
	Duration    *int
	CreatedAt   time.Time

	Setting     FormSetting   `gorm:"foreignKey:FormID"`
	FormSection []FormSection `gorm:"foreignKey:FormID"`
	Questions   []Question    `gorm:"foreignKey:FormID"`
	Submissions []Submission  `gorm:"foreignKey:FormID"`
}
type FormSetting struct {
	ID                 uint      `gorm:"primaryKey"`
	FormID             uuid.UUID `gorm:"type:char(36);not null;uniqueIndex"`
	ShowResult         bool      `gorm:"default:true"`  // mengizinkan user untuk dapat melihat hasil secara langsung
	MultipleSubmission bool      `gorm:"default:false"` // pengaturan apakah user yang sama dapat melakukan submission berulang
	PassingGrade       *float64  `gorm:"default:0"`     // threshold untuk kelulusan khusus untuk quiz dan exam
	Grading            bool      `gorm:"default:false"` // mengaktifkan sistem grading (khusus untuk quiz dan exam)
	MaxSubmissions     *int      `gorm:"default:100"`   // jumlah responden yang dapat mengisi form
	StartAt            *time.Time
	EndAt              *time.Time
}

type FormSection struct {
	ID          uuid.UUID `gorm:"type:char(36);primaryKey"`
	FormID      uuid.UUID `gorm:"type:char(36);primaryKey"`
	Title       string    `gorm:"type:varchar(255);not null"`
	Description string    `gorm:"type:text"`
	Order       int
}

type Question struct {
	ID         uuid.UUID  `gorm:"type:char(36);primaryKey"`
	FormID     uuid.UUID  `gorm:"type:char(36);not null;index"`
	SectionID  *uuid.UUID `gorm:"type:char(36);not null;index"`
	Text       string     `gorm:"type:text;not null"`
	Type       string     `gorm:"type:varchar(50);not null"`
	IsRequired bool       `gorm:"default:false"`
	Order      int
	Score      *int
	ImageURL   *string `gorm:"type:varchar(255)"`

	Options []Option `gorm:"foreignKey:QuestionID"`
}

type Option struct {
	ID         uint      `gorm:"primaryKey;autoIncrement"`
	QuestionID uuid.UUID `gorm:"type:char(36);not null;index"`
	Text       string    `gorm:"type:varchar(255);not null"`
	ImageURL   *string   `gorm:"type:varchar(255)"`
	IsCorrect  *bool     `gorm:"default:null"`
}

type Submission struct {
	ID          uuid.UUID `gorm:"type:char(36);primaryKey"`
	FormID      uuid.UUID `gorm:"type:char(36);not null;index"`
	Email       string    `gorm:"type:varchar(100)"`
	Score       *float64
	SubmittedAt time.Time

	// untuk identifikasi user agar tidak dapat melakukan spam dari form (jika formsetting diaktifkan)
	IPAddress    *string `gorm:"type:varchar(45)"`
	UserAgent    *string `gorm:"type:text"`
	SessionToken *string `gorm:"type:char(36);index"`

	Answers []Answer `gorm:"foreignKey:SubmissionID"`
}

type Answer struct {
	ID           uuid.UUID `gorm:"type:char(36);primaryKey"`
	SubmissionID uuid.UUID `gorm:"type:char(36);not null;index"`
	QuestionID   uuid.UUID `gorm:"type:char(36);not null;index"`
	OptionID     *uint
	TextAnswer   *string
}

// opsional untuk fitur form diagnosa
type Queue struct {
	ID         uuid.UUID `gorm:"type:char(36);primaryKey"`
	ResponseID uuid.UUID `gorm:"type:char(36);not null;index"`
	QueuNumber int
	Status     string `gorm:"type:varchar(20);default:'pending';check:status IN ('waiting','progress','done');not null"`
}
