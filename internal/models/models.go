package models

import (
	"time"
)

// User represents a user in the system (Employer or Job Seeker).
type User struct {
	ID       uint           `gorm:"primaryKey"`
	Name     string         `gorm:"size:100;not null"`
	Email    string         `gorm:"unique;not null"`
	Password string         `gorm:"not null"` // Store hashed passwords
	Role     string         `gorm:"size:20;not null"` // Employer or Job Seeker
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Profile represents a job seeker's profile.
type Profile struct {
	ID       uint           `gorm:"primaryKey"`
	UserID   uint           `gorm:"not null;unique"` // Linked to User
	Bio      string         `gorm:"type:text"`
	Resume	 string			`gorm:"type:text"`
	Skills   string         `gorm:"type:text"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Experience represents a job seeker's work experience.
type Experience struct {
	ID        uint           `gorm:"primaryKey"`
	ProfileID uint           `gorm:"not null"`
	Company   string         `gorm:"size:100"`
	Title     string         `gorm:"size:100"`
	StartDate time.Time
	EndDate   *time.Time // Nullable
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Company represents an employer's company.
type Company struct {
	ID          uint           `gorm:"primaryKey"`
	UserID      uint           `gorm:"not null;unique"` // Linked to User
	Name        string         `gorm:"size:100;not null"`
	Description string         `gorm:"type:text"`
	Location	string		   `gorm:"type:text"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// Job represents a job posting.
type Job struct {
	ID          uint           `gorm:"primaryKey"`
	CompanyID   uint           `gorm:"not null"`
	Title       string         `gorm:"size:100;not null"`
	Description string         `gorm:"type:text"`
	Location    string         `gorm:"size:100"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// Application represents a job application by a job seeker.
type Application struct {
	ID        uint           `gorm:"primaryKey"`
	JobID     uint           `gorm:"not null"`
	UserID    uint           `gorm:"not null"`
	Status    string         `gorm:"size:20;default:'Pending'"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
