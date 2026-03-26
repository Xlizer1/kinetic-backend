package repositories

import (
	"kinetic-backend/internal/models"

	"gorm.io/gorm"
)

type MessageRepository struct {
	db *gorm.DB
}

func NewMessageRepository(db *gorm.DB) *MessageRepository {
	return &MessageRepository{db: db}
}

func (r *MessageRepository) Create(message *models.Message) error {
	return r.db.Create(message).Error
}

func (r *MessageRepository) FindByID(id uint) (*models.Message, error) {
	var message models.Message
	err := r.db.Preload("Author").First(&message, id).Error
	if err != nil {
		return nil, err
	}
	return &message, nil
}

func (r *MessageRepository) FindByChannelID(channelID uint, limit, offset int) ([]models.Message, error) {
	var messages []models.Message
	err := r.db.Preload("Author").
		Where("channel_id = ?", channelID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&messages).Error
	return messages, err
}

func (r *MessageRepository) Update(message *models.Message) error {
	return r.db.Save(message).Error
}

func (r *MessageRepository) Delete(id uint) error {
	return r.db.Delete(&models.Message{}, id).Error
}
