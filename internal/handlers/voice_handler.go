package handlers

import (
	"strconv"

	"kinetic-backend/internal/realtime"
	"kinetic-backend/internal/services"
	"kinetic-backend/internal/utils"

	"github.com/gin-gonic/gin"
)

type VoiceHandler struct {
	voiceService *services.VoiceService
	hub          *realtime.Hub
}

func NewVoiceHandler(voiceService *services.VoiceService, hub *realtime.Hub) *VoiceHandler {
	return &VoiceHandler{voiceService: voiceService, hub: hub}
}

// @Summary Join a voice channel
// @Description Joins a voice channel
// @Tags voice
// @Accept json
// @Produce json
// @Param id path uint true "Channel ID"
// @Success 200 {object} models.VoiceState
// @Security BearerAuth
// @Router /channels/{id}/voice/join [post]
func (h *VoiceHandler) JoinVoice(c *gin.Context) {
	userID := c.GetUint("userID")
	username := c.GetString("username")

	channelID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "Invalid channel ID")
		return
	}

	state, err := h.voiceService.JoinVoice(services.JoinVoiceInput{
		ChannelID: uint(channelID),
		UserID:    userID,
	})
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	// Broadcast to WebSocket room
	if h.hub != nil {
		event := realtime.Event{
			Type: realtime.EventVoiceUserJoined,
			Payload: realtime.VoiceUserPayload{
				ChannelID: state.ChannelID,
				UserID:    state.UserID,
				Username:  username,
				IsMuted:   state.IsMuted,
				IsDeaf:    state.IsDeaf,
			},
		}
		if room, ok := h.hub.Rooms[state.ChannelID]; ok {
			room.BroadcastMessage(realtime.MustMarshal(event))
		}
	}

	utils.SuccessResponse(c, state)
}

// @Summary Leave a voice channel
// @Description Leaves a voice channel
// @Tags voice
// @Produce json
// @Param id path uint true "Channel ID"
// @Success 200 {string} string "Left voice channel"
// @Security BearerAuth
// @Router /channels/{id}/voice/leave [post]
func (h *VoiceHandler) LeaveVoice(c *gin.Context) {
	userID := c.GetUint("userID")
	username := c.GetString("username")

	channelID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "Invalid channel ID")
		return
	}

	if err := h.voiceService.LeaveVoice(uint(channelID), userID); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	// Broadcast to WebSocket room
	if h.hub != nil {
		event := realtime.Event{
			Type: realtime.EventVoiceUserLeft,
			Payload: realtime.VoiceUserPayload{
				ChannelID: uint(channelID),
				UserID:    userID,
				Username:  username,
			},
		}
		if room, ok := h.hub.Rooms[uint(channelID)]; ok {
			room.BroadcastMessage(realtime.MustMarshal(event))
		}
	}

	utils.SuccessMessage(c, "Left voice channel")
}

// @Summary Get voice channel users
// @Description Returns users in a voice channel
// @Tags voice
// @Produce json
// @Param id path uint true "Channel ID"
// @Success 200 {array} models.VoiceState
// @Security BearerAuth
// @Router /channels/{id}/voice [get]
func (h *VoiceHandler) GetVoiceUsers(c *gin.Context) {
	channelID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "Invalid channel ID")
		return
	}

	users, err := h.voiceService.GetChannelUsers(uint(channelID))
	if err != nil {
		utils.InternalServerError(c, "Failed to get voice users")
		return
	}

	utils.SuccessResponse(c, users)
}
