package utils

import (
	"fmt"
	"log"
	"math/rand"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

func LoadEnv() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, relying on system ENV")
	}
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func GenerateOTP(length int) string {
	digits := "0123456789"
	var sb strings.Builder

	for i := 0; i < length; i++ {
		sb.WriteByte(digits[rand.Intn(len(digits))])
	}

	return sb.String()
}

func MustGetUserID(c *gin.Context) string {
	userID, exists := c.Get("userID")
	if !exists {
		panic("userID not found in context")
	}
	idStr, ok := userID.(string)
	if !ok {
		panic("userID in context is not a string")
	}
	return idStr
}

func RandomUserAvatar(fullname string) string {
	return fmt.Sprintf("https://api.dicebear.com/6.x/initials/svg?seed=%s", fullname)
}
