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
