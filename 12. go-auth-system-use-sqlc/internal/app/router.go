package app

import (
	"fmt"
	"net/http"

	_ "go-auth-system/api/docs" // Import generated docs
	"go-auth-system/internal/auth/delivery"
	customMiddleware "go-auth-system/internal/middleware"
	"go-auth-system/internal/util"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
)

// RouterConfig contains handlers required to map routes
type RouterConfig struct {
	Config      util.Config
	AuthHandler *delivery.AuthHandler
}

// InitRouter initializes and configures the entire routing system
func InitRouter(cfg RouterConfig) http.Handler {
	r := chi.NewRouter()

	// --- 1. GLOBAL MIDDLEWARES (MUST be defined BEFORE any routes) ---
	// Sequence matters: Trace -> RequestID -> RealIP -> Logger -> Recoverer
	r.Use(customMiddleware.TracingMiddleware) // Your custom tracing (RequestID & StartTime)
	r.Use(middleware.RequestID)               // Standard Chi RequestID
	r.Use(middleware.RealIP)                  // Get client real IP for logging/security
	r.Use(middleware.Logger)                  // Structured request logging
	r.Use(middleware.Recoverer)               // Prevent server crash on panic

	// --- 2. PUBLIC INFRASTRUCTURE ROUTES ---

	// Swagger Documentation
	r.Get("/swagger/*", httpSwagger.WrapHandler)

	// Health Check (Critical for Docker/K8s Liveness & Readiness probes)
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "OK")
	})

	// --- 3. API V1 ROUTES ---
	r.Route("/api/v1", func(r chi.Router) {
		// Unprotected Authentication Routes
		r.Group(func(r chi.Router) {
			r.Post("/auth/register", cfg.AuthHandler.Register)
			r.Post("/auth/login", cfg.AuthHandler.Login)
			r.Post("/auth/refresh-token", cfg.AuthHandler.RefreshToken)
			r.Post("/auth/forgot-password", cfg.AuthHandler.ForgotPassword)
			r.Post("/auth/reset-password", cfg.AuthHandler.ResetPassword)
		})

		// Protected Routes (Authentication Required)
		r.Group(func(r chi.Router) {
			// Middleware in Group MUST also be defined before the group's routes
			r.Use(customMiddleware.AuthMiddleware(cfg.Config))

			r.Post("/auth/logout", cfg.AuthHandler.Logout)
			r.Get("/auth/profile", cfg.AuthHandler.GetProfile)
			r.Put("/auth/profile", cfg.AuthHandler.UpdateProfile)
		})
	})

	return r
}
