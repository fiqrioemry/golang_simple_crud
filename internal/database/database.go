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
	dialect := mysql.Open(dsn)
	var err error
	DB, err = gorm.Open(dialect, &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	err = DB.AutoMigrate(
		&models.User{},
		&models.Company{},
		&models.Job{},
		&models.Application{},
		&models.Profile{},
	)
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	fmt.Println("Database connection established and migrations completed")
}
