package models

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// [Database]
type User struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Password  string    `json:"-"` // Không trả về password
	Provider  string    `json:"provider"`
	FullName  string    `json:"full_name"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Product struct {
	ID          int     `json:"id"`
	CategoryID  int     `json:"category_id"`
	Name        string  `json:"name"`
	Price       float64 `json:"price"`
	Stock       int     `json:"stock"`
	Description string  `json:"description"`
}

// [End] Database

// [Login]
// Struct cho Login request
type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Struct cho Generate JWT Token cho Login response
type Claims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// [End] Login

// [Register]
// Struct cho Register request
type RegisterInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	FullName string `json:"full_name"`
}

// [End] Register
