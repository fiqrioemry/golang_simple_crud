package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `gorm:"type:char(36);primaryKey" json:"id"`
	Email     string    `gorm:"type:varchar(255);uniqueIndex;not null"json:"email"`
	Password  string    `gorm:"type:text;not null"json:"-"`
	Fullname  string    `gorm:"type:varchar(255);not null" json:"fullname"`
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
	LastResetAt        time.Time `gorm:"not null"`

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
	ID        uuid.UUID `gorm:"type:char(36);primaryKey"`
	UserID    uuid.UUID `gorm:"type:char(36);not null;index"`
	Title     string    `gorm:"type:varchar(255);not null"`
	Slug      string    `gorm:"type:varchar(100);uniqueIndex;not null"`
	Type      string    `gorm:"type:varchar(50);not null"`
	IsActive  bool      `gorm:"default:true"`
	Duration  *int      // in minutes, only for quiz
	CreatedAt time.Time

	User        User          `gorm:"foreignKey:UserID"`
	Questions   []Question    `gorm:"foreignKey:FormID"`
	MenuItems   []MenuItem    `gorm:"foreignKey:FormID"`
	Tickets     []QueueTicket `gorm:"foreignKey:FormID"`
	Submissions []Submission  `gorm:"foreignKey:FormID"`
}

type Question struct {
	ID         uuid.UUID `gorm:"type:char(36);primaryKey"`
	FormID     uuid.UUID `gorm:"type:char(36);not null;index"`
	Text       string    `gorm:"type:text;not null"`
	Type       string    `gorm:"type:varchar(50);not null"` // short_text, paragraph, multiple_choice, etc.
	IsRequired bool      `gorm:"default:false"`
	Score      *int      // only for quiz
	ImageURL   *string   `gorm:"type:varchar(255)"`

	Form    Form     `gorm:"foreignKey:FormID"`
	Options []Option `gorm:"foreignKey:QuestionID"`
}

type Option struct {
	ID         uint      `gorm:"primaryKey;autoIncrement"`
	QuestionID uuid.UUID `gorm:"type:char(36);not null;index"`
	Text       string    `gorm:"type:varchar(255);not null"`
	IsCorrect  *bool     `gorm:"default:null"` // true only for quiz, nil for survey

	Question Question `gorm:"foreignKey:QuestionID"`
}

type Submission struct {
	ID          uuid.UUID `gorm:"type:char(36);primaryKey"`
	FormID      uuid.UUID `gorm:"type:char(36);not null;index"`
	Email       string    `gorm:"type:varchar(100)"`
	Verified    bool      `gorm:"default:false"`
	Score       *float64  // only for quiz
	SubmittedAt time.Time

	Form    Form     `gorm:"foreignKey:FormID"`
	Answers []Answer `gorm:"foreignKey:SubmissionID"`
}

type Answer struct {
	ID           uuid.UUID `gorm:"type:char(36);primaryKey"`
	SubmissionID uuid.UUID `gorm:"type:char(36);not null;index"`
	QuestionID   uuid.UUID `gorm:"type:char(36);not null;index"`
	OptionID     *uint     // if answer is multiple_choice
	TextAnswer   *string   // if answer is short_text/paragraph

	Submission Submission `gorm:"foreignKey:SubmissionID"`
	Question   Question   `gorm:"foreignKey:QuestionID"`
}

type MenuItem struct {
	ID     uuid.UUID `gorm:"type:char(36);primaryKey"`
	FormID uuid.UUID `gorm:"type:char(36);not null;index"`
	Name   string    `gorm:"type:varchar(100);not null"`
	Price  float64   `gorm:"not null"`

	Form Form `gorm:"foreignKey:FormID"`
}

type QueueTicket struct {
	ID        uuid.UUID `gorm:"type:char(36);primaryKey"`
	FormID    uuid.UUID `gorm:"type:char(36);not null;index"`
	TicketNo  int       `gorm:"not null"`
	Name      string    `gorm:"type:varchar(100)"`
	CreatedAt time.Time

	OrderItems []OrderItem `gorm:"foreignKey:TicketID"`
	Form       Form        `gorm:"foreignKey:FormID"`
}

type OrderItem struct {
	ID       uint      `gorm:"primaryKey;autoIncrement"`
	TicketID uuid.UUID `gorm:"type:char(36);not null;index"`
	MenuID   uuid.UUID `gorm:"type:char(36);not null;index"`
	Quantity int       `gorm:"not null;default:1"`

	MenuItem    MenuItem    `gorm:"foreignKey:MenuID"`
	QueueTicket QueueTicket `gorm:"foreignKey:TicketID"`
}
