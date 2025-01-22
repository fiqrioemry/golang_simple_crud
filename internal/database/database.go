// /internal/database/database.go
package database

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
    err := godotenv.Load() 
    if err != nil {
        log.Fatal("Error loading .env file")
    }

    dbUser := os.Getenv("DB_USER")
    dbPassword := os.Getenv("DB_PASSWORD")
    dbHost := os.Getenv("DB_HOST")
    dbPort := os.Getenv("DB_PORT")
    dbName := os.Getenv("DB_NAME")

    dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", 
                       dbUser, dbPassword, dbHost, dbPort, dbName)

    var err2 error
    DB, err2 = gorm.Open(mysql.Open(dsn), &gorm.Config{})
    if err2 != nil {
        log.Fatal("Failed to connect to the database:", err2)
    }
    log.Println("Connected to the database successfully")
}
