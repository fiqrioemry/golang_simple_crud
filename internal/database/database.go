package database

import (
	"fmt"
	"golang_project/internal/models"
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	dsn := os.Getenv("MYSQL_DSN")
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Auto-migrate models in the correct order
	err = DB.AutoMigrate(
		&models.User{},        // Create users table first
		&models.Company{},     // Create companies table
		&models.Job{},         // Create jobs table
		&models.Profile{},     // Create profiles table after users
		&models.Application{}, // Create applications table after jobs
	)
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	fmt.Println("Database connection established and migrations completed")
}
