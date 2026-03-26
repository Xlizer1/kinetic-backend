package models

import (
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

type UserSettings struct {
	Theme         string `json:"theme"`
	Language      string `json:"language"`
	Notifications bool   `json:"notifications"`
	Privacy       string `json:"privacy"`
}

func (s *UserSettings) ToJSON() string {
	data, _ := json.Marshal(s)
	return string(data)
}

func (s *UserSettings) FromJSON(data string) error {
	return json.Unmarshal([]byte(data), s)
}

type User struct {
	ID                   uint           `gorm:"primarykey" json:"id"`
	Email                string         `gorm:"uniqueIndex;not null" json:"email"`
	Username             string         `gorm:"uniqueIndex;not null" json:"username"`
	PasswordHash         string         `gorm:"not null" json:"-"`
	AvatarURL            string         `json:"avatar_url"`
	Settings             string         `gorm:"type:text" json:"settings"`
	PasswordResetToken   string         `gorm:"-" json:"-"`
	PasswordResetExpires time.Time      `gorm:"type:timestamp" json:"-"`
	CreatedAt            time.Time      `json:"created_at"`
	UpdatedAt            time.Time      `json:"updated_at"`
	DeletedAt            gorm.DeletedAt `gorm:"index" json:"-"`

	Servers []ServerMember `gorm:"foreignKey:UserID" json:"servers,omitempty"`
}

func (u *User) GetSettings() UserSettings {
	var settings UserSettings
	settings.FromJSON(u.Settings)
	return settings
}

func (u *User) SetSettings(settings UserSettings) {
	u.Settings = settings.ToJSON()
}

type ServerMember struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	UserID    uint      `gorm:"not null" json:"user_id"`
	ServerID  uint      `gorm:"not null" json:"server_id"`
	Nickname  string    `json:"nickname"`
	Role      string    `gorm:"default:'member'" json:"role"`
	JoinedAt  time.Time `json:"joined_at"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	User   User   `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Server Server `gorm:"foreignKey:ServerID" json:"server,omitempty"`
}
