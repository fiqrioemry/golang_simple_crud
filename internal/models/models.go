package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

type Role string

const (
	Seeker   Role = "seeker"
	Employer Role = "employer"
)


type User struct {
	ID        uint      `gorm:"primaryKey"`
	Name      string    `gorm:"size:100;not null"`
	Email     string    `gorm:"unique;not null"`
	Password  string    `gorm:"not null"`
	Role      Role      `gorm:"not null;default:'seeker'"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Profile struct {
	ID        uint      `gorm:"primaryKey"`
	UserID    uint      `gorm:"not null:unique"` 
	Bio       string    `gorm:"type:text"`
	Resume    string    `gorm:"type:text"`
	Skills    Skills    `gorm:"type:json;default:null" json:"skills"`
	CreatedAt time.Time
	UpdatedAt time.Time
}


type Skills []string


func (s Skills) Value() (driver.Value, error) {
	if len(s) == 0 {
		return nil, nil // Return NULL if skills empty
	}
	return json.Marshal(s)
}

func (s *Skills) Scan(value interface{}) error {
	if value == nil {
		*s = []string{} 
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("failed to unmarshal JSON value into Skills")
	}
	return json.Unmarshal(bytes, s)
}


type Experience struct {
	ID        uint       `gorm:"primaryKey"`
	ProfileID uint       `gorm:"not null:unique"`
	Company   string     `gorm:"size:100"`
	Title     string     `gorm:"size:100"`
	StartDate time.Time
	EndDate   *time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}

// for company
type Company struct {
	ID          uint      `gorm:"primaryKey"`
	UserID      uint      `gorm:"not null;unique"`
	Name        string    `gorm:"size:100;not null"`
	Description string    `gorm:"type:text"`
	Location    string    `gorm:"type:text"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}



type Job struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	CompanyID   uint      `gorm:"not null"`
	Company     Company   `gorm:"foreignKey:CompanyID;references:ID" json:"company"`
	Title       string    `gorm:"not null" json:"title"`
	Description string    `gorm:"type:text;not null" json:"description"`
	Location    string    `gorm:"size:100;not null" json:"location"`
	Type        string    `gorm:"size:20;not null" json:"type"`
	Skills      Skills    `gorm:"type:json;default:null" json:"skills"`
	Experience  string    `gorm:"size:20;not null" json:"experience"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}


type JobResponse struct {
	ID          uint      `json:"id"`
	CompanyID   uint      `json:"company_id"` 
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Location    string    `json:"location"`
	Type        string    `json:"type"`
	Skills      Skills    `json:"skills"`
	Experience  string    `json:"experience"`
	CompanyName string    `json:"company_name"` 
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Application struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	JobID     uint      `gorm:"not null" json:"job_id"`
	UserID    uint      `gorm:"not null" json:"user_id"`
	Status    string    `gorm:"size:20;default:'Pending'" json:"status"`
	User      User      `gorm:"foreignKey:UserID;references:ID" json:"user"`
	Profile   Profile   `gorm:"foreignKey:UserID;references:UserID" json:"profile"`
	Job       Job       `gorm:"foreignKey:JobID;references:ID" json:"job"`  
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}


type ApplicationResponse struct {
	ID        uint      `json:"id"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// User details
	User struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	} `json:"user"`

	// Profile details
	Profile struct {
		Bio    string   `json:"bio"`
		Resume string   `json:"resume"`
		Skills []string `json:"skills"`
	} `json:"profile"`

	// Company details
	Company struct {
		ID   uint   `json:"company_id"`
		Name string `json:"company_name"`
	} `json:"company"`
}
