package models

import (
	"time"
)

type Role string

const (
	Seeker   Role = "seeker"
	Employer Role = "employer"
)

type User struct {
	ID           uint          `gorm:"primaryKey"`
	Name         string        `gorm:"size:100;not null"`
	Email        string        `gorm:"unique;not null"`
	Password     string        `gorm:"not null"`
	Role         Role          `gorm:"not null;default:'seeker'"`
	Profile      Profile       `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE"`
	Company      Company       `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE"`
	Applications []Application `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type Company struct {
	ID          uint   `gorm:"primaryKey"`
	UserID      uint   `gorm:"unique; not null"`
	Name        string `gorm:"size:100;not null"`
	Description string `gorm:"type:text"`
	Location    string `gorm:"type:text"`
	Jobs        []Job  `gorm:"foreignKey:CompanyID;constraint:OnUpdate:CASCADE"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Profile struct {
	ID         uint         `gorm:"primaryKey"`
	UserID     uint         `gorm:"unique; not null"`
	Bio        string       `gorm:"type:text"`
	Resume     string       `gorm:"type:text"`
	Skills     []string     `gorm:"type:json;not null"`
	Experience []Experience `gorm:"foreignKey:ProfileID;constraint:OnUpdate:CASCADE"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type Experience struct {
	ID        uint   `gorm:"primaryKey"`
	ProfileID uint   `gorm:"not null"`
	Company   string `gorm:"size:100"`
	Title     string `gorm:"size:100"`
	StartDate time.Time
	EndDate   *time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Application struct {
	ID        uint   `gorm:"primaryKey"`
	JobID     uint   `gorm:"not null"`
	UserID    uint   `gorm:"not null"`
	Status    string `gorm:"size:20;default:'Pending'"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Job       Job  `gorm:"foreignKey:JobID;constraint:OnUpdate:CASCADE"`
	User      User `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE"`
}

type Job struct {
	ID           uint          `gorm:"primaryKey"`
	CompanyID    uint          `gorm:"not null"`
	Title        string        `gorm:"not null"`
	Description  string        `gorm:"type:text;not null"`
	Location     string        `gorm:"size:100;not null"`
	Type         string        `gorm:"size:20;not null"`
	Experience   string        `gorm:"size:20;not null"`
	Skills       []string      `gorm:"type:json;not null"`
	Company      Company       `gorm:"foreignKey:CompanyID;constraint:OnUpdate:CASCADE"`
	Applications []Application `gorm:"foreignKey:JobID;constraint:OnUpdate:CASCADE"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
