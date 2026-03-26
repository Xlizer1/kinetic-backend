package services

import (
	"kinetic-backend/internal/models"
	"kinetic-backend/internal/repositories"
)

type MessageService struct {
	messageRepo *repositories.MessageRepository
}

func NewMessageService(messageRepo *repositories.MessageRepository) *MessageService {
	return &MessageService{messageRepo: messageRepo}
}

type CreateMessageInput struct {
	ChannelID uint   `json:"channel_id" binding:"required"`
	AuthorID  uint   `json:"author_id" binding:"required"`
	Content   string `json:"content" binding:"required"`
}

func (s *MessageService) CreateMessage(input CreateMessageInput) (*models.Message, error) {
	message := &models.Message{
		ChannelID: input.ChannelID,
		AuthorID:  input.AuthorID,
		Content:   input.Content,
	}

	if err := s.messageRepo.Create(message); err != nil {
		return nil, err
	}

	return s.messageRepo.FindByID(message.ID)
}

func (s *MessageService) GetChannelMessages(channelID uint, limit, offset int) ([]models.Message, error) {
	if limit <= 0 {
		limit = 50
	}
	return s.messageRepo.FindByChannelID(channelID, limit, offset)
}

func (s *MessageService) DeleteMessage(id uint) error {
	return s.messageRepo.Delete(id)
}
