package models

import (
	"time"
)

type Role string

const (
	Seeker 		Role = "seeker"
	Employer 	Role = "employer"
)

// User represents a user in the system (Employer or Job Seeker).
type User struct {
	ID       	uint           `gorm:"primaryKey"`
	Name     	string         `gorm:"size:100;not null"`
	Email    	string         `gorm:"unique;not null"`
	Password 	string         `gorm:"not null"`
	Role     	Role      	   `gorm:"not null;default:'seeker'"`
	CreatedAt 	time.Time
	UpdatedAt 	time.Time
}

// Profile represents a job seeker's profile.
type Profile struct {
	ID      	uint           `gorm:"primaryKey"`
	UserID   	uint           `gorm:"not null;unique"`
	Bio      	string         `gorm:"type:text"`
	Resume	 	string			`gorm:"type:text"`
	Skills   	[]string  `gorm:"type:json;not null" json:"skills"`
	CreatedAt 	time.Time
	UpdatedAt 	time.Time
}

// Experience represents a job seeker's work experience.
type Experience struct {
	ID        	uint           `gorm:"primaryKey"`
	ProfileID 	uint           `gorm:"not null"`
	Company   	string         `gorm:"size:100"`
	Title     	string         `gorm:"size:100"`
	StartDate 	time.Time
	EndDate   	*time.Time 
	CreatedAt 	time.Time
	UpdatedAt 	time.Time
}

// Company represents an employer's company.
type Company struct {
	ID          uint           `gorm:"primaryKey"`
	UserID      uint           `gorm:"not null;unique"` 
	Name        string         `gorm:"size:100;not null"`
	Description string         `gorm:"type:text"`
	Location	string		   `gorm:"type:text"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// Job represents a job posting.
// Job represents a job posting.
type Job struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	CompanyID   uint      `gorm:"not null" json:"company_id"` // Foreign key to Company
	Title       string    `gorm:"not null" json:"title"`
	Description string    `gorm:"type:text;not null" json:"description"`
	Location    string    `gorm:"size:100;not null" json:"location"`
	Type        string    `gorm:"size:20;not null" json:"type"`       // Enum: [fulltime, contract, freelance, internship]
	Skills      []string  `gorm:"type:json;not null" json:"skills"`   // Array of strings stored as JSON
	Experience  string    `gorm:"size:20;not null" json:"experience"` // Enum: [fresh graduate, 0-5 tahun, 5-10 tahun]
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
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
