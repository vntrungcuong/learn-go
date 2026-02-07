package main

import (
	"ecommerce-api/internal/config"
	"ecommerce-api/internal/handlers"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	// 1. Connect to DB
	config.ConnectDB()

	// 2. Close DB connection when app exits
	defer config.ClostDB()

	// 3. Khởi tạo Fiber App
	app := fiber.New(fiber.Config{
		AppName: "Go Ecommerce API",
	})

	// 4. Middleware
	app.Use(logger.New()) // Log request ra console
	app.Use(cors.New())   // Cho phép Frontend gọi API

	// 5. Setup Routes
	api := app.Group("/api/v1")

	// Test Route
	api.Get("/ping", func(c *fiber.Ctx) error {
		return c.SendString("Pong! Server is running 🚀")
	})

	// Auth Routes
	auth := api.Group("/auth")
	auth.Post("/register", handlers.RegisterHandler)
	auth.Post("/reset-password", handlers.ResetPasswordHanlder)
	auth.Post("/login", handlers.LoginHandler)
	auth.Post("/refresh-token", handlers.RefreshTokenHanlder)
	auth.Post("/logout", handlers.LogoutHanlder)
	/*
		*** Authentication flow: ***
		S1: Login: return accessToken (15m) and refreshToken (7 days)
		S2: User use accessToken to use other API
		S3: When accessToken expired, use API will response status 401 (Unauthorired), Client call API refresh-token to get new accessToken
		Logic refresh-token API: Server check DB, if OK will generate new accessToken and delete old accessToken
		S4: Logout: Delete refreshToken to DB to revoke token
	*/

	// 6. Start Server
	log.Fatal(app.Listen(":3000"))
}
