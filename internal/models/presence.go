package models

import (
	"time"

	"gorm.io/gorm"
)

type Presence struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	UserID    uint           `gorm:"uniqueIndex;not null" json:"user_id"`
	Status    string         `gorm:"default:'offline'" json:"status"`
	LastSeen  time.Time      `json:"last_seen"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

const (
	PresenceOnline  = "online"
	PresenceIdle    = "idle"
	PresenceDND     = "dnd"
	PresenceOffline = "offline"
)

func (Presence) TableName() string {
	return "presences"
}
