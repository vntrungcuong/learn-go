package main

import (
	"go-auth-system/internal/delivery/http"
	"go-auth-system/internal/domain"
	"go-auth-system/internal/repository"
	"go-auth-system/internal/usecase"
	"go-auth-system/pkg/database"
	"go-auth-system/pkg/logger"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	// 1. Init Logger
	log := logger.NewLogger()

	// 2. Init Database
	db := database.NewPostgreDB("localhost", "postgres", "Aa@123456", "go-auth-system", "5432", "disable")

	// Auto migrate (Note: Only for dev. Production should use migration tools such as golang-migrate)
	db.AutoMigrate(&domain.User{})

	// 3. Setup layeres (Manual Dependency injection)
	// Repo
	userRepo := repository.NewUserRepository(db)

	// Usecase
	userUsecase := usecase.NewUserUsecase(userRepo, time.Second*2, "my_super_secret_jwt_key")

	// 4. Setup Router & Middleware
	r := gin.Default()

	// Middleware logging custom (Optional)
	r.Use(func(c *gin.Context) {
		log.Info("Request received",
			zap.String("path", c.Request.URL.Path),
			zap.String("method", c.Request.Method),
		)
		c.Next()
	})

	// 5. Register Handlers
	http.NewUserHandler(r, userUsecase)
	r.GET("/health", http.HealthCheck)

	// 6. Start Server
	log.Info("Starting server on port 3000")
	r.Run(":3000")
}
