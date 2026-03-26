package handlers

import (
	"strconv"

	"kinetic-backend/internal/services"
	"kinetic-backend/internal/utils"

	"github.com/gin-gonic/gin"
)

type ChannelHandler struct {
	channelService *services.ChannelService
}

func NewChannelHandler(channelService *services.ChannelService) *ChannelHandler {
	return &ChannelHandler{channelService: channelService}
}

// @Summary Get server channels
// @Description Returns all channels in a server
// @Tags channels
// @Produce json
// @Param id path uint true "Server ID"
// @Success 200 {array} models.Channel
// @Security BearerAuth
// @Router /servers/{id}/channels [get]
func (h *ChannelHandler) GetChannels(c *gin.Context) {
	serverID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "Invalid server ID")
		return
	}

	channels, err := h.channelService.GetServerChannels(uint(serverID))
	if err != nil {
		utils.InternalServerError(c, "Failed to fetch channels")
		return
	}

	utils.SuccessResponse(c, channels)
}

// @Summary Create a new channel
// @Description Creates a new channel in a server
// @Tags channels
// @Accept json
// @Produce json
// @Param channel body services.CreateChannelInput true "Channel data"
// @Success 200 {object} models.Channel
// @Security BearerAuth
// @Router /channels [post]
func (h *ChannelHandler) CreateChannel(c *gin.Context) {
	var input services.CreateChannelInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	channel, err := h.channelService.CreateChannel(input)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.SuccessResponse(c, channel)
}

// @Summary Get channel by ID
// @Description Returns a channel by its ID
// @Tags channels
// @Produce json
// @Param id path uint true "Channel ID"
// @Success 200 {object} models.Channel
// @Security BearerAuth
// @Router /channels/{id} [get]
func (h *ChannelHandler) GetChannel(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "Invalid channel ID")
		return
	}

	channel, err := h.channelService.GetChannelByID(uint(id))
	if err != nil {
		utils.NotFound(c, "Channel not found")
		return
	}

	utils.SuccessResponse(c, channel)
}

// @Summary Update channel
// @Description Updates a channel's properties
// @Tags channels
// @Accept json
// @Produce json
// @Param id path uint true "Channel ID"
// @Param channel body services.UpdateChannelInput true "Channel update data"
// @Success 200 {object} models.Channel
// @Security BearerAuth
// @Router /channels/{id} [patch]
func (h *ChannelHandler) UpdateChannel(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "Invalid channel ID")
		return
	}

	var input services.UpdateChannelInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	channel, err := h.channelService.UpdateChannel(uint(id), input)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.SuccessResponse(c, channel)
}

// @Summary Delete channel
// @Description Deletes a channel
// @Tags channels
// @Produce json
// @Param id path uint true "Channel ID"
// @Success 200 {string} string "Channel deleted"
// @Security BearerAuth
// @Router /channels/{id} [delete]
func (h *ChannelHandler) DeleteChannel(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "Invalid channel ID")
		return
	}

	if err := h.channelService.DeleteChannel(uint(id)); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.SuccessMessage(c, "Channel deleted")
}
