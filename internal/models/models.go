package models

import (
	"time"

	"gorm.io/gorm"
)

type Role string

type User struct {
	gorm.Model
	Email    string `gorm:"unique;not null"`
	Password string `gorm:"not null"`
	Role     Role   `gorm:"not null;default:'seeker'"`
}

// Seeker ini sebagai pencari kerja
type Seeker struct {
	gorm.Model
	UserID       uint          `gorm:"unique;not null"`
	Name         string        `gorm:"size:100;not null"`
	Avatar       string        `gorm:"type:text"`
	Gender       string        `gorm:"size:20"`
	Bio          string        `gorm:"type:text"`
	Resume       string        `gorm:"type:text"`
	Birthday     *time.Time    `gorm:"type:date"`
	Skills       []string      `gorm:"serializer:json"`
	Experience   []Experience  `gorm:"foreignKey:SeekerID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Applications []Application `gorm:"foreignKey:SeekerID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	User         User          `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

type Experience struct {
	gorm.Model
	SeekerID  uint      `gorm:"not null"`
	Company   string    `gorm:"size:100"`
	Title     string    `gorm:"size:100"`
	StartDate time.Time `gorm:"type:date"`
	EndDate   *time.Time
}

// Employer ini sebagai company
type Employer struct {
	gorm.Model
	UserID      uint   `gorm:"unique;not null"`
	Name        string `gorm:"size:100;not null"`
	Avatar      string `gorm:"type:text"`
	Picture     string `gorm:"type:text"`
	Description string `gorm:"type:text"`
	Location    string `gorm:"type:text"`
	Jobs        []Job  `gorm:"foreignKey:EmployerID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	User        User   `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

type Job struct {
	gorm.Model
	EmployerID   uint          `gorm:"not null"`
	Title        string        `gorm:"size:100;not null"`
	Description  string        `gorm:"type:text;not null"`
	Salary       float64       `gorm:"not null"`
	Location     string        `gorm:"size:100;not null"`
	Type         string        `gorm:"size:20;not null"`
	Experience   string        `gorm:"size:20;not null"`
	Skills       []string      `gorm:"serializer:json"`
	IsActive     bool          `gorm:"default:true"`
	Employer     Employer      `gorm:"foreignKey:EmployerID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Applications []Application `gorm:"foreignKey:JobID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

type Application struct {
	gorm.Model
	JobID    uint   `gorm:"not null"`
	SeekerID uint   `gorm:"not null"`
	Status   string `gorm:"size:20;default:'Pending'"`
	Job      Job    `gorm:"foreignKey:JobID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Seeker   Seeker `gorm:"foreignKey:SeekerID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

type EmployerResponse struct {
	ID          uint      `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	UserID      uint      `json:"user_id"`
	Name        string    `json:"name"`
	Avatar      string    `json:"avatar"`
	Picture     string    `json:"picture"`
	Description string    `json:"description"`
	Location    string    `json:"location"`
}
