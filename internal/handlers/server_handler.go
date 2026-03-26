package handlers

import (
	"strconv"

	"kinetic-backend/internal/services"
	"kinetic-backend/internal/utils"

	"github.com/gin-gonic/gin"
)

type ServerHandler struct {
	serverService *services.ServerService
}

func NewServerHandler(serverService *services.ServerService) *ServerHandler {
	return &ServerHandler{serverService: serverService}
}

// @Summary Get user servers
// @Description Returns all servers the authenticated user is a member of
// @Tags servers
// @Produce json
// @Success 200 {array} models.Server
// @Security BearerAuth
// @Router /servers [get]
func (h *ServerHandler) GetServers(c *gin.Context) {
	userID := c.GetUint("userID")

	servers, err := h.serverService.GetUserServers(userID)
	if err != nil {
		utils.InternalServerError(c, "Failed to fetch servers")
		return
	}

	utils.SuccessResponse(c, servers)
}

// @Summary Create a new server
// @Description Creates a new server and adds the user as owner
// @Tags servers
// @Accept json
// @Produce json
// @Param server body services.CreateServerInput true "Server data"
// @Success 200 {object} models.Server
// @Security BearerAuth
// @Router /servers [post]
func (h *ServerHandler) CreateServer(c *gin.Context) {
	userID := c.GetUint("userID")

	var input services.CreateServerInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	server, err := h.serverService.CreateServer(userID, input)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.SuccessResponse(c, server)
}

// @Summary Get server by ID
// @Description Returns a server by its ID
// @Tags servers
// @Produce json
// @Param id path uint true "Server ID"
// @Success 200 {object} models.Server
// @Security BearerAuth
// @Router /servers/{id} [get]
func (h *ServerHandler) GetServer(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "Invalid server ID")
		return
	}

	server, err := h.serverService.GetServerByID(uint(id))
	if err != nil {
		utils.NotFound(c, "Server not found")
		return
	}

	utils.SuccessResponse(c, server)
}

// @Summary Update server
// @Description Updates a server's properties
// @Tags servers
// @Accept json
// @Produce json
// @Param id path uint true "Server ID"
// @Param server body services.UpdateServerInput true "Server update data"
// @Success 200 {object} models.Server
// @Security BearerAuth
// @Router /servers/{id} [patch]
func (h *ServerHandler) UpdateServer(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "Invalid server ID")
		return
	}

	var input services.UpdateServerInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	server, err := h.serverService.UpdateServer(uint(id), input)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.SuccessResponse(c, server)
}

// @Summary Delete server
// @Description Deletes a server (owner only)
// @Tags servers
// @Produce json
// @Param id path uint true "Server ID"
// @Success 200 {string} string "Server deleted"
// @Security BearerAuth
// @Router /servers/{id} [delete]
func (h *ServerHandler) DeleteServer(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "Invalid server ID")
		return
	}

	if err := h.serverService.DeleteServer(uint(id)); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.SuccessMessage(c, "Server deleted")
}

// @Summary Join server
// @Description Joins a server using invite code
// @Tags servers
// @Accept json
// @Produce json
// @Param invite_code body object{invite_code=string} true "Invite code"
// @Success 200 {object} models.Server
// @Security BearerAuth
// @Router /servers/join [post]
func (h *ServerHandler) JoinServer(c *gin.Context) {
	userID := c.GetUint("userID")

	var input struct {
		InviteCode string `json:"invite_code" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	server, err := h.serverService.JoinServerByInviteCode(userID, input.InviteCode)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.SuccessResponse(c, server)
}

// @Summary Leave server
// @Description Removes the user from a server
// @Tags servers
// @Produce json
// @Param id path uint true "Server ID"
// @Success 200 {string} string "Left server successfully"
// @Security BearerAuth
// @Router /servers/{id}/leave [post]
func (h *ServerHandler) LeaveServer(c *gin.Context) {
	userID := c.GetUint("userID")
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "Invalid server ID")
		return
	}

	if err := h.serverService.LeaveServer(userID, uint(id)); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.SuccessMessage(c, "Left server successfully")
}
