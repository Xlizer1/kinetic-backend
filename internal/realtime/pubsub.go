package realtime

import (
	"context"
	"encoding/json"
	"log"
	"sync"

	"github.com/redis/go-redis/v9"
)

const (
	WSChannel = "kinetic:ws:events"
)

type PubSubHub struct {
	redis    *redis.Client
	pubsub   *redis.PubSub
	localHub *Hub
	serverID string
	mu       sync.RWMutex
	stopped  bool
}

func NewPubSubHub(redisClient *redis.Client, localHub *Hub, serverID string) *PubSubHub {
	return &PubSubHub{
		redis:    redisClient,
		localHub: localHub,
		serverID: serverID,
	}
}

func (p *PubSubHub) Start(ctx context.Context) error {
	pubsub := p.redis.Subscribe(ctx, WSChannel)
	ch := pubsub.Channel()

	go func() {
		for msg := range ch {
			p.handleMessage(msg.Payload)
		}
	}()

	p.pubsub = pubsub
	log.Printf("Redis PubSub started, listening on channel: %s", WSChannel)
	return nil
}

func (p *PubSubHub) handleMessage(payload string) {
	var pubSubMsg PubSubMessage
	if err := json.Unmarshal([]byte(payload), &pubSubMsg); err != nil {
		log.Printf("Failed to unmarshal pubsub message: %v", err)
		return
	}

	if pubSubMsg.ServerID == p.serverID {
		return
	}

	switch pubSubMsg.Type {
	case "broadcast":
		var event Event
		if err := json.Unmarshal([]byte(pubSubMsg.Payload), &event); err != nil {
			log.Printf("Failed to unmarshal event: %v", err)
			return
		}
		p.localHub.BroadcastToAll(event, pubSubMsg.ExcludeUserID)

	case "room_broadcast":
		var msg RoomBroadcastPayload
		if err := json.Unmarshal([]byte(payload), &msg); err != nil {
			log.Printf("Failed to unmarshal room broadcast: %v", err)
			return
		}
		var event Event
		if err := json.Unmarshal([]byte(msg.Payload), &event); err != nil {
			log.Printf("Failed to unmarshal event: %v", err)
			return
		}
		p.localHub.BroadcastToRoom(msg.ChannelID, event, msg.ExcludeUserID)

	case "user_message":
		var msg UserMessagePayload
		if err := json.Unmarshal([]byte(payload), &msg); err != nil {
			log.Printf("Failed to unmarshal user message: %v", err)
			return
		}
		var event Event
		if err := json.Unmarshal([]byte(msg.Payload), &event); err != nil {
			log.Printf("Failed to unmarshal event: %v", err)
			return
		}
		p.localHub.SendToUser(msg.TargetUserID, event)
	}
}

func (p *PubSubHub) Publish(eventType string, payload interface{}, excludeUserID uint) error {
	msg := PubSubMessage{
		Type:          eventType,
		ServerID:      p.serverID,
		ExcludeUserID: excludeUserID,
		Payload:       "",
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	msg.Payload = string(payloadBytes)

	msgBytes, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	return p.redis.Publish(context.Background(), WSChannel, msgBytes).Err()
}

func (p *PubSubHub) PublishRoomBroadcast(channelID uint, event Event, excludeUserID uint) error {
	msg := PubSubMessage{
		Type:          "room_broadcast",
		ServerID:      p.serverID,
		ExcludeUserID: excludeUserID,
		ChannelID:     channelID,
	}

	payloadBytes, err := json.Marshal(event)
	if err != nil {
		return err
	}
	msg.Payload = string(payloadBytes)

	msgBytes, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	return p.redis.Publish(context.Background(), WSChannel, msgBytes).Err()
}

func (p *PubSubHub) PublishUserMessage(targetUserID uint, event Event) error {
	msg := PubSubMessage{
		Type:         "user_message",
		ServerID:     p.serverID,
		TargetUserID: targetUserID,
	}

	payloadBytes, err := json.Marshal(event)
	if err != nil {
		return err
	}
	msg.Payload = string(payloadBytes)

	msgBytes, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	return p.redis.Publish(context.Background(), WSChannel, msgBytes).Err()
}

func (p *PubSubHub) Stop() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.pubsub != nil {
		return p.pubsub.Close()
	}
	return nil
}

type PubSubMessage struct {
	Type          string `json:"type"`
	ServerID      string `json:"server_id"`
	TargetUserID  uint   `json:"target_user_id,omitempty"`
	ExcludeUserID uint   `json:"exclude_user_id,omitempty"`
	ChannelID     uint   `json:"channel_id,omitempty"`
	Payload       string `json:"payload"`
}

type RoomBroadcastPayload struct {
	Type          string `json:"type"`
	ServerID      string `json:"server_id"`
	ChannelID     uint   `json:"channel_id"`
	ExcludeUserID uint   `json:"exclude_user_id"`
	Payload       string `json:"payload"`
}

type UserMessagePayload struct {
	Type         string `json:"type"`
	ServerID     string `json:"server_id"`
	TargetUserID uint   `json:"target_user_id"`
	Payload      string `json:"payload"`
}
