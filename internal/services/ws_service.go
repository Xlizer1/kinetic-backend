package services

import (
	"kinetic-backend/internal/models"
	"kinetic-backend/internal/repositories"
	"kinetic-backend/internal/utils"
)

type WsService struct {
	userRepo    *repositories.UserRepository
	messageRepo *repositories.MessageRepository
}

func NewWsService(userRepo *repositories.UserRepository, messageRepo *repositories.MessageRepository) *WsService {
	return &WsService{
		userRepo:    userRepo,
		messageRepo: messageRepo,
	}
}

func (s *WsService) AuthenticateToken(token string) (uint, string, error) {
	claims, err := utils.ValidateToken(token)
	if err != nil {
		return 0, "", err
	}
	return claims.UserID, claims.Username, nil
}

func (s *WsService) SaveMessage(channelID, authorID uint, content string) error {
	message := &models.Message{
		ChannelID: channelID,
		AuthorID:  authorID,
		Content:   content,
	}
	return s.messageRepo.Create(message)
}

func (s *WsService) CreateMessage(channelID, authorID uint, content string) error {
	message := &models.Message{
		ChannelID: channelID,
		AuthorID:  authorID,
		Content:   content,
	}
	return s.messageRepo.Create(message)
}
