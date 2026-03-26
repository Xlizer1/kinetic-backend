package models

import (
	"time"

	"gorm.io/gorm"
)

type VoiceState struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	ChannelID uint           `gorm:"not null;index" json:"channel_id"`
	UserID    uint           `gorm:"not null;index" json:"user_id"`
	IsMuted   bool           `gorm:"default:false" json:"is_muted"`
	IsDeaf    bool           `gorm:"default:false" json:"is_deaf"`
	JoinedAt  time.Time      `json:"joined_at"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Channel Channel `gorm:"foreignKey:ChannelID" json:"channel,omitempty"`
	User    User    `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

func (VoiceState) TableName() string {
	return "voice_states"
}
