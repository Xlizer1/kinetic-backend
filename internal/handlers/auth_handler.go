package handlers

import (
	"kinetic-backend/internal/services"
	"kinetic-backend/internal/utils"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService *services.AuthService
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// @Summary Register a new user
// @Description Creates a new user account with email, username, and password
// @Tags auth
// @Accept json
// @Produce json
// @Param user body services.RegisterInput true "User registration data"
// @Success 200 {object} services.AuthResponse
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var input services.RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	result, err := h.authService.Register(input)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.SuccessResponse(c, result)
}

// @Summary Login user
// @Description Authenticates user and returns JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body services.LoginInput true "Login credentials"
// @Success 200 {object} services.AuthResponse
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var input services.LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	result, err := h.authService.Login(input)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.SuccessResponse(c, result)
}

// @Summary Forgot password
// @Description Generates password reset token and returns reset link
// @Tags auth
// @Accept json
// @Produce json
// @Param email body object{email=string} true "Email address"
// @Success 200 {object} map[string]string
// @Router /auth/forgot-password [post]
func (h *AuthHandler) ForgotPassword(c *gin.Context) {
	var input services.ForgotPasswordInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	resetLink, err := h.authService.ForgotPassword(input.Email)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.SuccessResponse(c, gin.H{
		"message":    "Password reset link generated",
		"reset_link": resetLink,
	})
}

// @Summary Reset password
// @Description Resets user password using reset token
// @Tags auth
// @Accept json
// @Produce json
// @Param reset body services.ResetPasswordInput true "Reset token and new password"
// @Success 200 {string} string "Password reset successful"
// @Router /auth/reset-password [post]
func (h *AuthHandler) ResetPassword(c *gin.Context) {
	var input services.ResetPasswordInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	if err := h.authService.ResetPassword(input.Token, input.NewPassword); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.SuccessMessage(c, "Password reset successful")
}

// @Summary Verify email
// @Description Verifies user email address (placeholder for future implementation)
// @Tags auth
// @Accept json
// @Produce json
// @Param token body object{token=string} true "Verification token"
// @Success 200 {string} string "Email verification placeholder"
// @Router /auth/verify-email [post]
func (h *AuthHandler) VerifyEmail(c *gin.Context) {
	utils.SuccessMessage(c, "Email verification placeholder - implementation pending")
}

// @Summary Refresh token
// @Description Refreshes JWT access token using refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param token body services.RefreshTokenInput true "Refresh token"
// @Success 200 {object} services.AuthResponse
// @Router /auth/refresh-token [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var input services.RefreshTokenInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	result, err := h.authService.RefreshToken(input)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.SuccessResponse(c, result)
}
