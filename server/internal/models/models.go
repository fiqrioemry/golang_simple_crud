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
	UserID             uuid.UUID `gorm:"type:uuid;not null;uniqueIndex"` // hanya 1 aktif per user
	SubscriptionTierID uint      `gorm:"not null"`
	StartedAt          time.Time `gorm:"not null"`
	ExpiresAt          time.Time `gorm:"not null"`
	IsActive           bool      `gorm:"default:true"`

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
