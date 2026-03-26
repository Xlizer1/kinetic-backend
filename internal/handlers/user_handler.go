package handlers

import (
	"kinetic-backend/internal/services"
	"kinetic-backend/internal/utils"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService *services.UserService
}

func NewUserHandler(userService *services.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// @Summary Get current user
// @Description Returns the authenticated user's profile
// @Tags users
// @Produce json
// @Success 200 {object} models.User
// @Security BearerAuth
// @Router /users/@me [get]
func (h *UserHandler) GetMe(c *gin.Context) {
	userID := c.GetUint("userID")

	user, err := h.userService.GetUserByID(userID)
	if err != nil {
		utils.NotFound(c, "User not found")
		return
	}

	utils.SuccessResponse(c, user)
}

// @Summary Update current user
// @Description Updates the authenticated user's profile
// @Tags users
// @Accept json
// @Produce json
// @Param user body services.UpdateUserInput true "User update data"
// @Success 200 {object} models.User
// @Security BearerAuth
// @Router /users/@me [patch]
func (h *UserHandler) UpdateMe(c *gin.Context) {
	userID := c.GetUint("userID")

	var input services.UpdateUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	user, err := h.userService.UpdateUser(userID, input)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.SuccessResponse(c, user)
}

// @Summary Update user settings
// @Description Updates the authenticated user's settings
// @Tags users
// @Accept json
// @Produce json
// @Param settings body services.UpdateSettingsInput true "User settings"
// @Success 200 {object} models.User
// @Security BearerAuth
// @Router /users/@me/settings [patch]
func (h *UserHandler) UpdateSettings(c *gin.Context) {
	userID := c.GetUint("userID")

	var input services.UpdateSettingsInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	user, err := h.userService.UpdateUserSettings(userID, input)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.SuccessResponse(c, user)
}
