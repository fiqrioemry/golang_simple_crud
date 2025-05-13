package main

import (
	"log"
	"os"
	"server/internal/config"
	"server/internal/handlers"
	"server/internal/middleware"
	"server/internal/repositories"
	"server/internal/routes"
	"server/internal/services"
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
	// ================== AUTH ========================
	authRepo := repositories.NewAuthRepository(db)
	authService := services.NewAuthService(authRepo)
	authHandler := handlers.NewAuthHandler(authService)

	// =================== USER =======================
	userRepo := repositories.NewUserRepository(db)
	userService := services.NewUserService(userRepo)
	userHandler := handlers.NewUserHandler(userService)

	// =============== QUEUE (DIAGNOSIS) ==============
	queueRepo := repositories.NewQueueRepository(db)
	queueService := services.NewQueueService(queueRepo)
	queueHandler := handlers.NewQueueHandler(queueService)

	// ===================== FORM =====================
	formRepo := repositories.NewFormRepository(db)
	formService := services.NewFormService(formRepo)
	formHandler := handlers.NewFormHandler(formService)

	// ===================== PAYMENT ===================
	paymentRepo := repositories.NewPaymentRepository(db)
	paymentService := services.NewPaymentService(paymentRepo, subscriptionRepo, authRepo)
	paymentHandler := handlers.NewPaymentHandler(paymentService)

	// ===================== ANALYTICS =================
	analyticsRepo := repositories.NewAnalyticsRepository(db)
	analyticsService := services.NewAnalyticsService(analyticsRepo)
	analyticsHandler := handlers.NewAnalyticsHandler(analyticsService)

	// ===================== SUBMISSION ================
	submissionRepo := repositories.NewSubmissionRepository(db)
	submissionService := services.NewSubmissionService(submissionRepo)
	submissionHandler := handlers.NewSubmissionHandler(submissionService)

	// =================== ADMIN SUBSCRIPTION ===========
	subscriptionRepo := repositories.NewSubscriptionRepository(db)
	subscriptionService := services.NewSubscriptionService(subscriptionRepo)
	subscriptionHandler := handlers.NewSubscriptionHandler(subscriptionService)

	// ========== Route Binding ==========
	routes.AuthRoutes(r, authHandler)
	routes.UserRoutes(r, userHandler)
	routes.PaymentRoutes(r, paymentHandler)
	routes.FormRoutes(r, formHandler)
	routes.QueueRoutes(r, queueHandler)
	routes.AnalyticsRoutes(r, analyticsHandler)
	routes.SubmissionRoutes(r, submissionHandler)
	routes.SubscriptionRoutes(r, subscriptionHandler)

	// ========== Start Server ==========
	port := os.Getenv("PORT")
	log.Println("server running on port:", port)
	log.Fatal(r.Run(":" + port))
}
