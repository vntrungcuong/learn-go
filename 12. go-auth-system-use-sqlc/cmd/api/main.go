package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	// Import Bootstrap logic
	"go-auth-system/internal/app"

	// Import Auth Layers (DDD)
	authApp "go-auth-system/internal/auth/app"
	"go-auth-system/internal/auth/delivery"
	"go-auth-system/internal/auth/infra"

	// Import Shared/Infra
	db "go-auth-system/internal/db/sqlc"
	"go-auth-system/internal/util"
)

// @title           Go Auth System API
// @version         1.0
// @description     Authentication & User Management System.
// @host            localhost:8080
// @BasePath        /
// @securityDefinitions.basic  BasicAuth
func main() {
	// --- 1. SET UP CONFIGURATION ---
	// LoadConfig sẽ đọc file .env và ánh xạ vào struct
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatalf("❌ Cannot load config: %v", err)
	}

	// --- 2. SET UP INFRASTRUCTURE (DB & Redis) ---
	// Sử dụng pgxpool.Pool để tối ưu hóa connection pooling
	connPool, err := app.ConnectPostgres(config)
	if err != nil {
		log.Fatalf("❌ Database connection failed: %v", err)
	}
	defer connPool.Close()

	// Kết nối Redis cho session management
	redisClient := app.ConnectRedis(config)
	defer redisClient.Close()

	// --- 3. DEPENDENCY INJECTION (DI) ASSEMBLY ---
	// SQLC Store đóng vai trò là Unit of Work, nhận vào pgxpool
	store := db.NewStore(connPool)

	// Lắp ráp Repositories (Infrastructure Layer)
	authRepo := infra.NewAuthRepository(store)
	sessionRepo := infra.NewSessionRepository(redisClient)

	// Lắp ráp Service (Application Layer) thực thi nghiệp vụ
	authService := authApp.NewAuthService(authRepo, sessionRepo, config)

	// Lắp ráp Handler (Delivery Layer) xử lý HTTP request
	authHandler := delivery.NewAuthHandler(authService)

	// --- 4. INITIALIZE ROUTER ---
	// InitRouter cấu hình Chi Router và các Middleware
	router := app.InitRouter(app.RouterConfig{
		Config:      config,
		AuthHandler: authHandler,
	})

	// --- 5. START SERVER WITH GRACEFUL SHUTDOWN ---
	server := &http.Server{
		Addr:         fmt.Sprintf(":%v", config.APIServerPort),
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Chạy server trong Goroutine để thực hiện Non-blocking
	go func() {
		log.Printf("🚀 Server is running on port %v", config.APIServerPort)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("❌ Listen error: %v\n", err)
		}
	}()

	// Lắng nghe tín hiệu tắt hệ thống từ OS (SIGINT, SIGTERM)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("⚠️ Shutting down server...")

	// Graceful Shutdown: Đợi 5 giây để hoàn tất các request đang xử lý
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("❌ Server forced to shutdown: %v", err)
	}

	log.Println("✅ Server exited gracefully")
}
