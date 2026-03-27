package realtime

import (
	"encoding/json"
	"log"
	"sync"
	"time"
)

type Hub struct {
	Rooms           map[uint]*Room
	VoiceRooms      map[uint]*Room
	Clients         map[uint]*Client
	ClientList      map[*Client]bool
	Register        chan *Client
	Unregister      chan *Client
	Broadcast       chan BroadcastMessage
	JoinRoom        chan JoinRoomMessage
	LeaveRoom       chan LeaveRoomMessage
	SendMessage     chan SendMessageToHub
	Authenticate    chan AuthenticateMessage
	VoiceJoin       chan VoiceJoinMessage
	VoiceLeave      chan VoiceLeaveMessage
	VoiceSignal     chan VoiceSignalMessage
	PresenceUpdate  chan PresenceMessage
	UserAuth        func(token string) (uint, string, error)
	SaveMessage     func(channelID, authorID uint, content string) error
	JoinVoice       func(channelID, userID uint) error
	LeaveVoice      func(channelID, userID uint) error
	GetVoiceUsers   func(channelID uint) ([]VoiceUserInfo, error)
	GetPresenceList func() ([]PresenceUserInfo, error)
	SetPresence     func(userID uint, status string) error
	PubSub          *PubSubHub
	ServerID        string
	mutex           sync.RWMutex
}

type VoiceUserInfo struct {
	UserID   uint
	Username string
}

type PresenceUserInfo struct {
	UserID   uint
	Username string
	Status   string
}

type VoiceJoinMessage struct {
	Client    *Client
	ChannelID uint
}

type VoiceLeaveMessage struct {
	Client    *Client
	ChannelID uint
}

type VoiceSignalMessage struct {
	Client       *Client
	ChannelID    uint
	TargetUserID uint
	SDP          string
	Candidate    string
	Type         string
}

type PresenceMessage struct {
	Client *Client
	Status string
}

type BroadcastMessage struct {
	Message []byte
	Exclude uint
}

func NewHub() *Hub {
	return &Hub{
		Rooms:          make(map[uint]*Room),
		VoiceRooms:     make(map[uint]*Room),
		Clients:        make(map[uint]*Client),
		ClientList:     make(map[*Client]bool),
		Register:       make(chan *Client),
		Unregister:     make(chan *Client),
		Broadcast:      make(chan BroadcastMessage),
		JoinRoom:       make(chan JoinRoomMessage),
		LeaveRoom:      make(chan LeaveRoomMessage),
		SendMessage:    make(chan SendMessageToHub),
		Authenticate:   make(chan AuthenticateMessage),
		VoiceJoin:      make(chan VoiceJoinMessage),
		VoiceLeave:     make(chan VoiceLeaveMessage),
		VoiceSignal:    make(chan VoiceSignalMessage),
		PresenceUpdate: make(chan PresenceMessage),
	}
}

func (h *Hub) Run() {
	cleanupTicker := time.NewTicker(60 * time.Second)
	defer cleanupTicker.Stop()

	for {
		select {
		case <-cleanupTicker.C:
			h.cleanupEmptyVoiceRooms()

		case client := <-h.Register:
			h.mutex.Lock()
			h.ClientList[client] = true
			h.mutex.Unlock()
			log.Printf("Client connected: %p", client)

		case client := <-h.Unregister:
			h.mutex.Lock()
			if _, ok := h.ClientList[client]; ok {
				delete(h.ClientList, client)
				if client.ID != 0 {
					delete(h.Clients, client.ID)
				}
				close(client.Send)
				for channelID := range client.Rooms {
					if room, ok := h.Rooms[channelID]; ok {
						room.Leave(client)
					}
				}
			}
			h.mutex.Unlock()
			log.Printf("Client disconnected: %p", client)

		case message := <-h.Broadcast:
			h.mutex.RLock()
			for client := range h.ClientList {
				if client.ID != message.Exclude {
					select {
					case client.Send <- message.Message:
					default:
						h.mutex.RUnlock()
						h.removeClient(client)
						h.mutex.RLock()
					}
				}
			}
			h.mutex.RUnlock()

		case authMsg := <-h.Authenticate:
			if h.UserAuth != nil {
				userID, username, err := h.UserAuth(authMsg.Token)
				if err != nil {
					authMsg.Client.SendEvent(Event{
						Type: EventError,
						Payload: ErrorPayload{
							Message: "Authentication failed",
						},
					})
					continue
				}
				authMsg.Client.ID = userID
				authMsg.Client.Username = username
				h.mutex.Lock()
				h.Clients[userID] = authMsg.Client
				h.mutex.Unlock()

				if h.SetPresence != nil {
					h.SetPresence(userID, "online")
				}

				if h.GetPresenceList != nil {
					users, _ := h.GetPresenceList()
					presenceUsers := make([]PresencePayload, len(users))
					for i, u := range users {
						presenceUsers[i] = PresencePayload{
							UserID:   u.UserID,
							Username: u.Username,
							Status:   u.Status,
						}
					}
					authMsg.Client.SendEvent(Event{
						Type: EventPresenceList,
						Payload: PresenceListPayload{
							Users: presenceUsers,
						},
					})
				}

				authMsg.Client.SendEvent(Event{
					Type: EventAuthenticate,
					Payload: map[string]interface{}{
						"success":  true,
						"user_id":  userID,
						"username": username,
					},
				})
			}

		case joinMsg := <-h.JoinRoom:
			log.Printf("[Hub] Received JOIN_ROOM for channel %d from user %d", joinMsg.ChannelID, joinMsg.Client.ID)
			room, ok := h.Rooms[joinMsg.ChannelID]
			if !ok {
				log.Printf("[Hub] Creating new room for channel %d", joinMsg.ChannelID)
				room = NewRoom(joinMsg.ChannelID, "text")
				h.Rooms[joinMsg.ChannelID] = room
			}
			room.Join(joinMsg.Client)
			log.Printf("[Hub] User %d joined room for channel %d", joinMsg.Client.ID, joinMsg.ChannelID)
			room.BroadcastMessage(MustMarshal(Event{
				Type: EventUserJoined,
				Payload: UserPayload{
					UserID:   joinMsg.Client.ID,
					Username: joinMsg.Client.Username,
				},
			}))

		case leaveMsg := <-h.LeaveRoom:
			if room, ok := h.Rooms[leaveMsg.ChannelID]; ok {
				room.Leave(leaveMsg.Client)
				room.BroadcastMessage(MustMarshal(Event{
					Type: EventUserLeft,
					Payload: UserPayload{
						UserID:   leaveMsg.Client.ID,
						Username: leaveMsg.Client.Username,
					},
				}))
			}

		case msg := <-h.SendMessage:
			log.Printf("[Hub] Received SEND_MESSAGE from user %d in channel %d: %s", msg.Client.ID, msg.ChannelID, msg.Content)
			if room, ok := h.Rooms[msg.ChannelID]; ok {
				log.Printf("[Hub] Room found for channel %d, saving message", msg.ChannelID)
				if h.SaveMessage != nil {
					if err := h.SaveMessage(msg.ChannelID, msg.Client.ID, msg.Content); err != nil {
						log.Printf("[Hub] Error saving message: %v", err)
						msg.Client.SendEvent(Event{
							Type: EventError,
							Payload: ErrorPayload{
								Message: "Failed to save message",
							},
						})
						continue
					}
					log.Printf("[Hub] Message saved successfully")
				}

				event := Event{
					Type: EventNewMessage,
					Payload: MessagePayload{
						ChannelID: msg.ChannelID,
						AuthorID:  msg.Client.ID,
						Content:   msg.Content,
						Username:  msg.Client.Username,
					},
				}
				log.Printf("[Hub] Broadcasting message to room channel %d", msg.ChannelID)
				room.BroadcastMessage(MustMarshal(event))
				log.Printf("[Hub] Message broadcast complete")
			} else {
				log.Printf("[Hub] ERROR: Room not found for channel %d", msg.ChannelID)
			}

		case voiceJoinMsg := <-h.VoiceJoin:
			if h.JoinVoice != nil {
				if err := h.JoinVoice(voiceJoinMsg.ChannelID, voiceJoinMsg.Client.ID); err != nil {
					voiceJoinMsg.Client.SendEvent(Event{
						Type: EventError,
						Payload: ErrorPayload{
							Message: err.Error(),
						},
					})
					continue
				}
			}

			room, ok := h.VoiceRooms[voiceJoinMsg.ChannelID]
			if !ok {
				room = NewRoom(voiceJoinMsg.ChannelID, "voice")
				h.VoiceRooms[voiceJoinMsg.ChannelID] = room
			}
			room.Join(voiceJoinMsg.Client)

			users := []VoiceUserInfo{}
			if h.GetVoiceUsers != nil {
				users, _ = h.GetVoiceUsers(voiceJoinMsg.ChannelID)
			}

			room.BroadcastMessage(MustMarshal(Event{
				Type: EventVoiceUserJoined,
				Payload: VoiceUserPayload{
					ChannelID: voiceJoinMsg.ChannelID,
					UserID:    voiceJoinMsg.Client.ID,
					Username:  voiceJoinMsg.Client.Username,
				},
			}))

			voiceJoinMsg.Client.SendEvent(Event{
				Type: EventVoiceJoin,
				Payload: map[string]interface{}{
					"channel_id": voiceJoinMsg.ChannelID,
					"users":      users,
				},
			})

		case voiceLeaveMsg := <-h.VoiceLeave:
			if h.LeaveVoice != nil {
				h.LeaveVoice(voiceLeaveMsg.ChannelID, voiceLeaveMsg.Client.ID)
			}

			if room, ok := h.VoiceRooms[voiceLeaveMsg.ChannelID]; ok {
				room.Leave(voiceLeaveMsg.Client)
				room.BroadcastMessage(MustMarshal(Event{
					Type: EventVoiceUserLeft,
					Payload: VoiceUserPayload{
						ChannelID: voiceLeaveMsg.ChannelID,
						UserID:    voiceLeaveMsg.Client.ID,
						Username:  voiceLeaveMsg.Client.Username,
					},
				}))

				if room.Count() == 0 {
					close(room.Broadcast)
					delete(h.VoiceRooms, voiceLeaveMsg.ChannelID)
				}
			}

		case signalMsg := <-h.VoiceSignal:
			h.mutex.RLock()
			targetClient, ok := h.Clients[signalMsg.TargetUserID]
			h.mutex.RUnlock()

			if !ok || targetClient == nil {
				continue
			}

			var eventType string
			var payload interface{}

			switch signalMsg.Type {
			case "offer":
				eventType = EventVoiceOffer
				payload = VoiceOfferPayload{
					ChannelID:    signalMsg.ChannelID,
					TargetUserID: signalMsg.Client.ID,
					SDP:          signalMsg.SDP,
				}
			case "answer":
				eventType = EventVoiceAnswer
				payload = VoiceAnswerPayload{
					ChannelID:    signalMsg.ChannelID,
					TargetUserID: signalMsg.Client.ID,
					SDP:          signalMsg.SDP,
				}
			case "ice":
				eventType = EventIceCandidate
				payload = IceCandidatePayload{
					ChannelID:    signalMsg.ChannelID,
					TargetUserID: signalMsg.Client.ID,
					Candidate:    signalMsg.Candidate,
				}
			}

			targetClient.SendEvent(Event{
				Type:    eventType,
				Payload: payload,
			})

		case presenceMsg := <-h.PresenceUpdate:
			if h.SetPresence != nil {
				h.SetPresence(presenceMsg.Client.ID, presenceMsg.Status)
			}

			h.BroadcastToAll(Event{
				Type: EventPresenceUpdate,
				Payload: PresencePayload{
					UserID:   presenceMsg.Client.ID,
					Username: presenceMsg.Client.Username,
					Status:   presenceMsg.Status,
				},
			}, presenceMsg.Client.ID)
		}
	}
}

func (h *Hub) BroadcastToRoom(channelID uint, event Event, exclude uint) {
	if room, ok := h.Rooms[channelID]; ok {
		msg := MustMarshal(event)
		for client := range room.Clients {
			if client.ID != exclude {
				select {
				case client.Send <- msg:
				default:
					room.Leave(client)
				}
			}
		}
	}
}

func (h *Hub) BroadcastToAll(event Event, exclude uint) {
	msg := BroadcastMessage{
		Message: MustMarshal(event),
		Exclude: exclude,
	}
	h.Broadcast <- msg

	if h.PubSub != nil {
		h.PubSub.Publish("broadcast", event, exclude)
	}
}

func (h *Hub) SendToUser(userID uint, event Event) {
	h.mutex.RLock()
	client, ok := h.Clients[userID]
	h.mutex.RUnlock()

	if ok && client != nil {
		client.SendEvent(event)
		return
	}

	if h.PubSub != nil {
		h.PubSub.PublishUserMessage(userID, event)
	}
}

func (h *Hub) removeClient(client *Client) {
	delete(h.ClientList, client)
	if client.ID != 0 {
		delete(h.Clients, client.ID)
	}
	close(client.Send)
	for channelID := range client.Rooms {
		if room, ok := h.Rooms[channelID]; ok {
			room.Leave(client)
		}
	}
}

func MustMarshal(v interface{}) []byte {
	data, err := json.Marshal(v)
	if err != nil {
		log.Printf("Error marshaling: %v", err)
		return []byte{}
	}
	return data
}

func (h *Hub) cleanupEmptyVoiceRooms() {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	for channelID, room := range h.VoiceRooms {
		if len(room.Clients) == 0 {
			close(room.Broadcast)
			delete(h.VoiceRooms, channelID)
			log.Printf("Cleaned up empty voice room: %d", channelID)
		}
	}
}
