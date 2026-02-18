// Package delivery handles HTTP requests for the authentication module.
package delivery

import (
	"fmt"
	"net/http"

	"go-auth-system/internal/auth/app"
	"go-auth-system/internal/pkg/rest"
	"go-auth-system/internal/util"
)

// AuthHandler coordinates authentication-related HTTP requests.
type AuthHandler struct {
	svc app.AuthService
}

// NewAuthHandler initializes a new AuthHandler with the provided AuthService.
func NewAuthHandler(svc app.AuthService) *AuthHandler {
	return &AuthHandler{svc: svc}
}

// Login handles the login endpoint.
// @Summary     Login user
// @Description Authenticate user and return access/refresh tokens
// @Tags        Authentication
// @Accept      json
// @Produce     json
// @Param       request body     LoginRequest true "Login credentials"
// @Success 200 {object} rest.APIResponse{data=AuthData}
// @Failure 401 {object} rest.APIResponse{data=any}
// @Router      /auth/login [post]
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if !rest.ParseAndValidate(w, r, &req) {
		return
	}

	accessToken, refreshToken, err := h.svc.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		rest.Error(w, r, http.StatusUnauthorized, err.Error(), nil)
		return
	}

	rest.Success(w, r, toAuthData(accessToken, refreshToken), "Login successful")
}

// Register handles new user account creation.
// @Summary     Register user
// @Description Creates a new user record in the system and hashes the password.
// @Tags        Authentication
// @Accept      json
// @Produce     json
// @Param       request body     RegisterRequest true "User registration details"
// @Success     201     {object} rest.APIResponse{data=UserData}
// @Failure     400     {object} rest.APIResponse{data=any}
// @Router      /auth/register [post]
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if !rest.ParseAndValidate(w, r, &req) {
		return
	}

	user, err := h.svc.Register(r.Context(), req.Email, req.Password, req.FullName)
	if err != nil {
		rest.Error(w, r, http.StatusBadRequest, err.Error(), nil)
		return
	}

	rest.Created(w, r, toUserData(user), "User registered successfully")
}

// GetProfile retrieves the authenticated user's profile.
// @Summary     Get user profile
// @Description Fetches profile information for the currently authenticated user.
// @Tags        User Profile
// @Security    BearerAuth
// @Produce     json
// @Success     200 {object} rest.APIResponse{data=UserData}
// @Failure     401 {object} rest.APIResponse{data=any}
// @Failure     404 {object} rest.APIResponse{data=any}
// @Router      /auth/profile [get]
func (h *AuthHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(util.UserIDKey).(int64)
	if !ok {
		rest.Error(w, r, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}

	user, err := h.svc.GetProfile(r.Context(), userID)
	if err != nil {
		rest.Error(w, r, http.StatusNotFound, err.Error(), nil)
		return
	}

	rest.Success(w, r, toUserData(user), "Get profile successfully")
}

// UpdateProfile modifies the authenticated user's information.
// @Summary     Update user profile
// @Description Updates the user's full name. Requires a valid Bearer token.
// @Tags        User Profile
// @Security    BearerAuth
// @Accept      json
// @Produce     json
// @Param       request body     UpdateProfileRequest true "Profile update data"
// @Success     200     {object} rest.APIResponse{data=UserData}
// @Failure     401     {object} rest.APIResponse{data=any}
// @Failure     500     {object} rest.APIResponse{data=any}
// @Router      /auth/profile [put]
func (h *AuthHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	var req UpdateProfileRequest
	if !rest.ParseAndValidate(w, r, &req) {
		return
	}

	userID, ok := r.Context().Value(util.UserIDKey).(int64)
	if !ok {
		rest.Error(w, r, http.StatusUnauthorized, "Unauthorized access", nil)
		return
	}

	user, err := h.svc.UpdateProfile(r.Context(), userID, req.FullName)
	if err != nil {
		rest.Error(w, r, http.StatusInternalServerError, "Failed to update user profile", nil)
		return
	}

	rest.Success(w, r, toUserData(user), "Profile updated successfully")
}

// Logout revokes the current user session.
// @Summary     Logout user
// @Description Revokes the user session by removing the refresh token from Redis.
// @Tags        Authentication
// @Security    BearerAuth
// @Produce     json
// @Success     200 {object} rest.APIResponse{data=any}
// @Failure     401 {object} rest.APIResponse{data=any}
// @Router      /auth/logout [post]
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value(util.UserIDKey).(int64)

	if err := h.svc.Logout(r.Context(), userID); err != nil {
		rest.Error(w, r, http.StatusInternalServerError, "Failed to logout", nil)
		return
	}

	rest.Success[any](w, r, nil, "Logged out successfully")
}

// ForgotPassword initiates the password recovery process.
// @Summary     Forgot password
// @Description Sends a password reset token. Always returns 200 for security.
// @Tags        Authentication
// @Accept      json
// @Produce     json
// @Param       request body     ForgotPasswordRequest true "User email"
// @Success     200     {object} rest.APIResponse{data=any}
// @Router      /auth/forgot-password [post]
func (h *AuthHandler) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	var req ForgotPasswordRequest
	if !rest.ParseAndValidate(w, r, &req) {
		return
	}

	_, err := h.svc.ForgotPassword(r.Context(), req.Email)
	if err != nil {
		// In production, use a structured logger (Zap/Logrus) here.
		fmt.Printf("Internal error in ForgotPassword: %v\n", err)
	}

	rest.Success[any](w, r, nil, "If email exists, a reset link has been sent")
}

// ResetPassword updates the user password using a recovery token.
// @Summary     Reset password
// @Description Updates the password using a valid reset token.
// @Tags        Authentication
// @Accept      json
// @Produce     json
// @Param       request body     ResetPasswordRequest true "Reset token and new password"
// @Success     200     {object} rest.APIResponse{data=any}
// @Failure     401     {object} rest.APIResponse{data=any}
// @Router      /auth/reset-password [post]
func (h *AuthHandler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	var req ResetPasswordRequest
	if !rest.ParseAndValidate(w, r, &req) {
		return
	}

	if err := h.svc.ResetPassword(r.Context(), req.Email, req.ResetToken, req.NewPassword); err != nil {
		rest.Error(w, r, http.StatusUnauthorized, err.Error(), nil)
		return
	}

	rest.Success[any](w, r, nil, "Password reset successfully")
}

// RefreshToken generates a new access token.
// @Summary     Refresh access token
// @Description Validates refresh token against Redis and returns a new access token.
// @Tags        Authentication
// @Accept      json
// @Produce     json
// @Param       request body     RefreshTokenRequest true "Refresh token details"
// @Success     200     {object} rest.APIResponse{data=AuthData}
// @Failure     401     {object} rest.APIResponse{data=any}
// @Router      /auth/refresh-token [post]
func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	var req RefreshTokenRequest
	if !rest.ParseAndValidate(w, r, &req) {
		return
	}

	newAccessToken, err := h.svc.RefreshToken(r.Context(), req.RefreshToken)
	if err != nil {
		rest.Error(w, r, http.StatusUnauthorized, err.Error(), nil)
		return
	}

	rest.Success(w, r, AuthData{AccessToken: newAccessToken}, "Token refreshed successfully")
}
