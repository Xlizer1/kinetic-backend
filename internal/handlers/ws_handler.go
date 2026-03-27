package handlers

import (
	"log"
	"net/http"

	"kinetic-backend/internal/realtime"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type WsHandler struct {
	Hub *realtime.Hub
}

func NewWsHandler(hub *realtime.Hub) *WsHandler {
	return &WsHandler{Hub: hub}
}

func (h *WsHandler) HandleWebSocket(c *gin.Context) {
	log.Printf("[WsHandler] WebSocket connection request from %s", c.Request.RemoteAddr)
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("[WsHandler] Failed to upgrade WebSocket: %v", err)
		return
	}

	log.Printf("[WsHandler] WebSocket upgraded successfully for %s", c.Request.RemoteAddr)
	client := realtime.NewClient(conn, h.Hub)
	h.Hub.Register <- client

	log.Printf("[WsHandler] Client registered, starting pumps")
	go client.WritePump()
	go client.ReadPump()
}
