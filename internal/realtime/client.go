package realtime

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512 * 1024
)

type Client struct {
	ID       uint
	Username string
	Conn     *websocket.Conn
	Send     chan []byte
	Hub      *Hub
	Rooms    map[uint]bool
}

func NewClient(conn *websocket.Conn, hub *Hub) *Client {
	return &Client{
		Conn:  conn,
		Send:  make(chan []byte, 256),
		Hub:   hub,
		Rooms: make(map[uint]bool),
	}
}

func (c *Client) ReadPump() {
	defer func() {
		log.Printf("[Client] ReadPump: Client %d (%s) disconnecting", c.ID, c.Username)
		c.Hub.Unregister <- c
		c.Conn.Close()
	}()

	log.Printf("[Client] ReadPump: Starting for client")
	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("[Client] WebSocket error: %v", err)
			}
			break
		}

		log.Printf("[Client] Received message: %s", string(message))
		var event Event
		if err := json.Unmarshal(message, &event); err != nil {
			log.Printf("[Client] Error unmarshaling message: %v", err)
			continue
		}

		log.Printf("[Client] Handling event type: %s", event.Type)
		c.HandleEvent(event)
	}
}

func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (c *Client) HandleEvent(event Event) {
	log.Printf("[Client] HandleEvent: received event type = %s", event.Type)
	switch event.Type {
	case "PING":
		log.Printf("[Client] Received PING, sending PONG")
		c.SendEvent(Event{Type: "PONG", Payload: map[string]interface{}{}})
	case EventAuthenticate:
		c.handleAuthenticate(event.Payload)
	case EventJoinRoom:
		c.handleJoinRoom(event.Payload)
	case EventLeaveRoom:
		c.handleLeaveRoom(event.Payload)
	case EventSendMessage:
		c.handleSendMessage(event.Payload)
	case EventTypingStart:
		c.handleTypingStart(event.Payload)
	case EventTypingStop:
		c.handleTypingStop(event.Payload)
	case EventVoiceJoin:
		c.handleVoiceJoin(event.Payload)
	case EventVoiceLeave:
		c.handleVoiceLeave(event.Payload)
	case EventVoiceOffer:
		c.handleVoiceOffer(event.Payload)
	case EventVoiceAnswer:
		c.handleVoiceAnswer(event.Payload)
	case EventIceCandidate:
		c.handleIceCandidate(event.Payload)
	case EventPresenceUpdate:
		c.handlePresenceUpdate(event.Payload)
	}
}

func (c *Client) handleAuthenticate(payload interface{}) {
	data, ok := payload.(map[string]interface{})
	if !ok {
		c.SendError("Invalid authentication payload")
		return
	}

	token, ok := data["token"].(string)
	if !ok {
		c.SendError("Missing token")
		return
	}

	c.Hub.Authenticate <- AuthenticateMessage{
		Client: c,
		Token:  token,
	}
}

func (c *Client) handleJoinRoom(payload interface{}) {
	log.Printf("[Client] handleJoinRoom called for client %d", c.ID)
	data, ok := payload.(map[string]interface{})
	if !ok {
		log.Printf("[Client] handleJoinRoom: Invalid payload type")
		c.SendError("Invalid join room payload")
		return
	}

	log.Printf("[Client] handleJoinRoom: payload = %+v", data)

	channelID, ok := data["channel_id"].(float64)
	if !ok {
		log.Printf("[Client] handleJoinRoom: Missing or invalid channel_id")
		c.SendError("Missing channel_id")
		return
	}

	channelIDUint := uint(channelID)
	log.Printf("[Client] handleJoinRoom: Sending to hub - channel=%d", channelIDUint)
	c.Hub.JoinRoom <- JoinRoomMessage{
		Client:    c,
		ChannelID: channelIDUint,
	}
	log.Printf("[Client] handleJoinRoom: Join request sent to hub")
}

func (c *Client) handleLeaveRoom(payload interface{}) {
	data, ok := payload.(map[string]interface{})
	if !ok {
		c.SendError("Invalid leave room payload")
		return
	}

	channelID, ok := data["channel_id"].(float64)
	if !ok {
		c.SendError("Missing channel_id")
		return
	}

	channelIDUint := uint(channelID)
	c.Hub.LeaveRoom <- LeaveRoomMessage{
		Client:    c,
		ChannelID: channelIDUint,
	}
}

func (c *Client) handleSendMessage(payload interface{}) {
	log.Printf("[Client] handleSendMessage called for client %d", c.ID)
	data, ok := payload.(map[string]interface{})
	if !ok {
		log.Printf("[Client] handleSendMessage: Invalid payload type")
		c.SendError("Invalid send message payload")
		return
	}

	log.Printf("[Client] handleSendMessage: payload = %+v", data)

	channelID, ok := data["channel_id"].(float64)
	if !ok {
		log.Printf("[Client] handleSendMessage: Missing or invalid channel_id, type = %T", data["channel_id"])
		c.SendError("Missing channel_id")
		return
	}

	content, ok := data["content"].(string)
	if !ok {
		log.Printf("[Client] handleSendMessage: Missing content")
		c.SendError("Missing content")
		return
	}

	log.Printf("[Client] handleSendMessage: Sending to hub - channel=%d, content=%s", uint(channelID), content)
	c.Hub.SendMessage <- SendMessageToHub{
		Client:    c,
		ChannelID: uint(channelID),
		Content:   content,
	}
	log.Printf("[Client] handleSendMessage: Message sent to hub")
}

func (c *Client) handleTypingStart(payload interface{}) {
	data, ok := payload.(map[string]interface{})
	if !ok {
		return
	}

	channelID, ok := data["channel_id"].(float64)
	if !ok {
		return
	}

	c.Hub.BroadcastToRoom(uint(channelID), Event{
		Type: EventTyping,
		Payload: TypingPayload{
			ChannelID: uint(channelID),
			UserID:    c.ID,
			Username:  c.Username,
		},
	}, c.ID)
}

func (c *Client) handleTypingStop(payload interface{}) {
}

func (c *Client) handleVoiceJoin(payload interface{}) {
	data, ok := payload.(map[string]interface{})
	if !ok {
		c.SendError("Invalid voice join payload")
		return
	}

	channelID, ok := data["channel_id"].(float64)
	if !ok {
		c.SendError("Missing channel_id")
		return
	}

	c.Hub.VoiceJoin <- VoiceJoinMessage{
		Client:    c,
		ChannelID: uint(channelID),
	}
}

func (c *Client) handleVoiceLeave(payload interface{}) {
	data, ok := payload.(map[string]interface{})
	if !ok {
		c.SendError("Invalid voice leave payload")
		return
	}

	channelID, ok := data["channel_id"].(float64)
	if !ok {
		c.SendError("Missing channel_id")
		return
	}

	c.Hub.VoiceLeave <- VoiceLeaveMessage{
		Client:    c,
		ChannelID: uint(channelID),
	}
}

func (c *Client) handleVoiceOffer(payload interface{}) {
	data, ok := payload.(map[string]interface{})
	if !ok {
		c.SendError("Invalid voice offer payload")
		return
	}

	channelID, _ := data["channel_id"].(float64)
	targetUserID, _ := data["target_user_id"].(float64)
	sdp, _ := data["sdp"].(string)

	c.Hub.VoiceSignal <- VoiceSignalMessage{
		Client:       c,
		ChannelID:    uint(channelID),
		TargetUserID: uint(targetUserID),
		SDP:          sdp,
		Type:         "offer",
	}
}

func (c *Client) handleVoiceAnswer(payload interface{}) {
	data, ok := payload.(map[string]interface{})
	if !ok {
		c.SendError("Invalid voice answer payload")
		return
	}

	channelID, _ := data["channel_id"].(float64)
	targetUserID, _ := data["target_user_id"].(float64)
	sdp, _ := data["sdp"].(string)

	c.Hub.VoiceSignal <- VoiceSignalMessage{
		Client:       c,
		ChannelID:    uint(channelID),
		TargetUserID: uint(targetUserID),
		SDP:          sdp,
		Type:         "answer",
	}
}

func (c *Client) handleIceCandidate(payload interface{}) {
	data, ok := payload.(map[string]interface{})
	if !ok {
		c.SendError("Invalid ICE candidate payload")
		return
	}

	channelID, _ := data["channel_id"].(float64)
	targetUserID, _ := data["target_user_id"].(float64)
	candidate, _ := data["candidate"].(string)

	c.Hub.VoiceSignal <- VoiceSignalMessage{
		Client:       c,
		ChannelID:    uint(channelID),
		TargetUserID: uint(targetUserID),
		Candidate:    candidate,
		Type:         "ice",
	}
}

func (c *Client) handlePresenceUpdate(payload interface{}) {
	data, ok := payload.(map[string]interface{})
	if !ok {
		c.SendError("Invalid presence update payload")
		return
	}

	status, ok := data["status"].(string)
	if !ok {
		c.SendError("Missing status")
		return
	}

	validStatuses := map[string]bool{
		"online":  true,
		"idle":    true,
		"dnd":     true,
		"offline": true,
	}

	if !validStatuses[status] {
		status = "online"
	}

	c.Hub.PresenceUpdate <- PresenceMessage{
		Client: c,
		Status: status,
	}
}

func (c *Client) SendError(message string) {
	c.Send <- MustMarshal(Event{
		Type: EventError,
		Payload: ErrorPayload{
			Message: message,
		},
	})
}

func (c *Client) SendEvent(event Event) {
	data, err := json.Marshal(event)
	if err != nil {
		log.Printf("Error marshaling event: %v", err)
		return
	}
	c.Send <- data
}

type AuthenticateMessage struct {
	Client *Client
	Token  string
}

type JoinRoomMessage struct {
	Client    *Client
	ChannelID uint
}

type LeaveRoomMessage struct {
	Client    *Client
	ChannelID uint
}

type SendMessageToHub struct {
	Client    *Client
	ChannelID uint
	Content   string
}
