package handlers

import (
	"ecommerce-api/internal/models"
	"ecommerce-api/internal/services"

	"github.com/gofiber/fiber/v2"
)

// POST /api/v1/auth/register
func RegisterHandler(c *fiber.Ctx) error {
	var input models.RegisterInput

	// 1. Parse JSON body
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// 2. Validate
	if input.Email == "" || input.Password == "" || input.FullName == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Missing required fields",
		})
	}

	// 3. Call service
	user, err := services.RegisterUser(c.Context(), input)
	if err != nil {
		if err.Error() == "email already exists" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Email already exists",
			})
		}
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// 4. Operation successful, response 201 Created
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "User registered successfully",
		"user": fiber.Map{
			"id":        user.ID,
			"email":     user.Email,
			"full_name": user.FullName,
		},
	})
}

// POST /api/v1/auth/login
func LoginHandler(c *fiber.Ctx) error {
	var input models.LoginInput

	// 1. Parse JSON body
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// 2. Validate
	if input.Email == "" || input.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Missing required fields",
		})
	}

	// 3. Call service
	token, err := services.LoginUser(c.Context(), input)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid email or password"})
	}

	// 4. Operation successful, response 200 OK
	return c.JSON(fiber.Map{
		"token": token,
		"type":  "Bearer",
	})
}

// POST /api/v1/auth/refresh-token
func RefreshTokenHanlder(c *fiber.Ctx) error {
	var input models.RefreshTokenInput

	// 1. Validate
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid input",
		})
	}

	// 2. Execute refresh token
	tokens, err := services.RefreshToken(c.Context(), input.RefreshToken)
	if err != nil {
		c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(tokens)
}

// POST /api/v1/auth/logout
func LogoutHanlder(c *fiber.Ctx) error {
	var input models.RefreshTokenInput

	// 1. Validate
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid input",
		})
	}

	// 2. Execute logout
	if err := services.LogoutUser(c.Context(), input.RefreshToken); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Logout failed",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Logged out successfully",
	})
}

// POST /api/v1/auth/reset-password
func ResetPasswordHanlder(c *fiber.Ctx) error {
	var input models.ResetPasswordInput

	// 1. Validate
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid input",
		})
	}

	// 2. Execute reset password
	if err := services.ResetPassword(c.Context(), input); err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Reset password successfully",
	})
}
