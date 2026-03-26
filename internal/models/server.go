package models

import (
	"time"

	"gorm.io/gorm"
)

type Server struct {
	ID         uint           `gorm:"primarykey" json:"id"`
	OwnerID    uint           `gorm:"not null" json:"owner_id"`
	Name       string         `gorm:"not null" json:"name"`
	IconURL    string         `json:"icon_url"`
	InviteCode string         `gorm:"uniqueIndex" json:"invite_code"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`

	Owner    User           `gorm:"foreignKey:OwnerID" json:"owner,omitempty"`
	Members  []ServerMember `gorm:"foreignKey:ServerID" json:"members,omitempty"`
	Channels []Channel      `gorm:"foreignKey:ServerID" json:"channels,omitempty"`
}

type Channel struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	ServerID  uint           `gorm:"not null" json:"server_id"`
	Name      string         `gorm:"not null" json:"name"`
	Type      string         `gorm:"default:'text'" json:"type"`
	Topic     string         `json:"topic"`
	Position  int            `gorm:"default:0" json:"position"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Server   Server    `gorm:"foreignKey:ServerID" json:"server,omitempty"`
	Messages []Message `gorm:"foreignKey:ChannelID" json:"messages,omitempty"`
}

type Message struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	ChannelID uint           `gorm:"not null" json:"channel_id"`
	AuthorID  uint           `gorm:"not null" json:"author_id"`
	Content   string         `gorm:"not null" json:"content"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Channel Channel `gorm:"foreignKey:ChannelID" json:"channel,omitempty"`
	Author  User    `gorm:"foreignKey:AuthorID" json:"author,omitempty"`
}
