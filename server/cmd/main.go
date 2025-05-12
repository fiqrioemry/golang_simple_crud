package main

import (
	"log"
	"os"
	"server/internal/config"
	"server/internal/middleware"
	"server/internal/routes"
	"server/internal/seeders"
	"server/internal/utils"

	"github.com/gin-gonic/gin"
)

func main() {
	utils.LoadEnv()
	config.InitRedis()
	config.InitMailer()
	config.InitDatabase()
	config.InitCloudinary()
	config.InitMidtrans()

	db := config.DB
	// ========== Seeder ==========
	seeders.ResetDatabase(db)

	// middleware config
	r := gin.Default()
	r.Use(
		middleware.Logger(),
		middleware.Recovery(),
		middleware.CORS(),
		middleware.RateLimiter(5, 10),
		middleware.LimitFileSize(12<<20),
		middleware.APIKeyGateway([]string{"/api/payments"}),
	)

	// ========== layer ==========
	authRepo := repositories.NewAuthRepository(db)
	authService := services.NewAuthService(authRepo)
	authHandler := handlers.NewAuthHandler(authService)
	// ========== Cron Job ==========

	// ========== Route Binding ==========
	routes.AuthRoutes(r, authHandler)

	// ========== Start Server ==========
	port := os.Getenv("PORT")
	log.Println("server running on port:", port)
	log.Fatal(r.Run(":" + port))
}
