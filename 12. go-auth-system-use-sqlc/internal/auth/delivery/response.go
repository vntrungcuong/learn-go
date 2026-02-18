package delivery

import (
	"time"
)

// --- Response Data DTOs ---

type AuthData struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

type UserData struct {
	ID        int64     `json:"id"`
	Email     string    `json:"email"`
	FullName  string    `json:"fullname"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
