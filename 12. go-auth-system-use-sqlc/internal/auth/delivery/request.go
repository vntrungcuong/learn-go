package delivery

// --- Request DTOs ---

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
	FullName string `json:"fullname" validate:"required"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type ResetPasswordRequest struct {
	Email       string `json:"email" validate:"required,email"`
	ResetToken  string `json:"reset_token" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=6"`
}

type UpdateProfileRequest struct {
	FullName string `json:"fullname" validate:"required"`
}
