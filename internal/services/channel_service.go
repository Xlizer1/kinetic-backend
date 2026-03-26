package services

import (
	"kinetic-backend/internal/models"
	"kinetic-backend/internal/repositories"
)

type ChannelService struct {
	channelRepo *repositories.ChannelRepository
}

func NewChannelService(channelRepo *repositories.ChannelRepository) *ChannelService {
	return &ChannelService{channelRepo: channelRepo}
}

type CreateChannelInput struct {
	Name     string `json:"name" binding:"required"`
	Type     string `json:"type"`
	Topic    string `json:"topic"`
	ServerID uint   `json:"server_id" binding:"required"`
}

type UpdateChannelInput struct {
	Name  string `json:"name"`
	Topic string `json:"topic"`
	Type  string `json:"type"`
}

func (s *ChannelService) CreateChannel(input CreateChannelInput) (*models.Channel, error) {
	channelType := "text"
	if input.Type != "" {
		channelType = input.Type
	}

	channel := &models.Channel{
		ServerID: input.ServerID,
		Name:     input.Name,
		Type:     channelType,
		Topic:    input.Topic,
	}

	if err := s.channelRepo.Create(channel); err != nil {
		return nil, err
	}

	return s.channelRepo.FindByID(channel.ID)
}

func (s *ChannelService) GetChannelByID(id uint) (*models.Channel, error) {
	return s.channelRepo.FindByID(id)
}

func (s *ChannelService) GetServerChannels(serverID uint) ([]models.Channel, error) {
	return s.channelRepo.FindByServerID(serverID)
}

func (s *ChannelService) UpdateChannel(id uint, input UpdateChannelInput) (*models.Channel, error) {
	channel, err := s.channelRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if input.Name != "" {
		channel.Name = input.Name
	}

	if input.Topic != "" {
		channel.Topic = input.Topic
	}

	if input.Type != "" {
		channel.Type = input.Type
	}

	if err := s.channelRepo.Update(channel); err != nil {
		return nil, err
	}

	return s.channelRepo.FindByID(id)
}

func (s *ChannelService) DeleteChannel(id uint) error {
	return s.channelRepo.Delete(id)
}
