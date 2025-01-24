package database

import (
	"fmt"
	"log"
	"os"

	"job-portal-api/internal/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

// ConnectDatabase initializes the database connection and runs migrations.
func ConnectDatabase() {
	dsn := os.Getenv("MYSQL_DSN") // Example: "user:password@tcp(localhost:3306)/job_portal?charset=utf8mb4&parseTime=True&loc=Local"
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Auto-migrate models
	err = DB.AutoMigrate(&models.User{}, &models.Profile{}, &models.Experience{}, &models.Company{}, &models.Job{}, &models.Application{})
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	fmt.Println("Database connection established and migrations applied!")
}
