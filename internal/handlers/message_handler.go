package handlers

import (
	"strconv"

	"kinetic-backend/internal/services"
	"kinetic-backend/internal/utils"

	"github.com/gin-gonic/gin"
)

type MessageHandler struct {
	messageService *services.MessageService
}

func NewMessageHandler(messageService *services.MessageService) *MessageHandler {
	return &MessageHandler{messageService: messageService}
}

// @Summary Get channel messages
// @Description Returns messages from a channel with pagination
// @Tags messages
// @Produce json
// @Param id path uint true "Channel ID"
// @Param limit query int false "Limit (default 50)"
// @Param offset query int false "Offset"
// @Success 200 {array} models.Message
// @Security BearerAuth
// @Router /channels/{id}/messages [get]
func (h *MessageHandler) GetMessages(c *gin.Context) {
	channelID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "Invalid channel ID")
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	messages, err := h.messageService.GetChannelMessages(uint(channelID), limit, offset)
	if err != nil {
		utils.InternalServerError(c, "Failed to fetch messages")
		return
	}

	utils.SuccessResponse(c, messages)
}

// @Summary Create a new message
// @Description Creates a new message in a channel
// @Tags messages
// @Accept json
// @Produce json
// @Param message body object{channel_id=uint,content=string} true "Message data"
// @Success 200 {object} models.Message
// @Security BearerAuth
// @Router /channels/{id}/messages [post]
func (h *MessageHandler) CreateMessage(c *gin.Context) {
	userID := c.GetUint("userID")

	var input struct {
		ChannelID uint   `json:"channel_id" binding:"required"`
		Content   string `json:"content" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	message, err := h.messageService.CreateMessage(services.CreateMessageInput{
		ChannelID: input.ChannelID,
		AuthorID:  userID,
		Content:   input.Content,
	})
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.SuccessResponse(c, message)
}

// @Summary Delete a message
// @Description Deletes a message
// @Tags messages
// @Produce json
// @Param id path uint true "Message ID"
// @Success 200 {string} string "Message deleted"
// @Security BearerAuth
// @Router /channels/{id}/messages/{messageId} [delete]
func (h *MessageHandler) DeleteMessage(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "Invalid message ID")
		return
	}

	if err := h.messageService.DeleteMessage(uint(id)); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.SuccessMessage(c, "Message deleted")
}
