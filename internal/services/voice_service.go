package services

import (
	"errors"
	"time"

	"kinetic-backend/internal/models"
	"kinetic-backend/internal/repositories"
)

type VoiceService struct {
	voiceRepo   *repositories.VoiceRepository
	channelRepo *repositories.ChannelRepository
}

func NewVoiceService(voiceRepo *repositories.VoiceRepository, channelRepo *repositories.ChannelRepository) *VoiceService {
	return &VoiceService{
		voiceRepo:   voiceRepo,
		channelRepo: channelRepo,
	}
}

type JoinVoiceInput struct {
	ChannelID uint
	UserID    uint
}

func (s *VoiceService) Join(channelID, userID uint) error {
	_, err := s.JoinVoice(JoinVoiceInput{
		ChannelID: channelID,
		UserID:    userID,
	})
	return err
}

func (s *VoiceService) Leave(channelID, userID uint) error {
	return s.LeaveVoice(channelID, userID)
}

func (s *VoiceService) JoinVoice(input JoinVoiceInput) (*models.VoiceState, error) {
	channel, err := s.channelRepo.FindByID(input.ChannelID)
	if err != nil {
		return nil, err
	}

	if channel.Type != "voice" {
		return nil, errors.New("channel is not a voice channel")
	}

	existing, _ := s.voiceRepo.FindByChannelAndUser(input.ChannelID, input.UserID)
	if existing != nil {
		return existing, nil
	}

	state := &models.VoiceState{
		ChannelID: input.ChannelID,
		UserID:    input.UserID,
		IsMuted:   false,
		IsDeaf:    false,
		JoinedAt:  time.Now(),
	}

	if err := s.voiceRepo.Create(state); err != nil {
		return nil, err
	}

	return state, nil
}

func (s *VoiceService) LeaveVoice(channelID, userID uint) error {
	return s.voiceRepo.Delete(channelID, userID)
}

func (s *VoiceService) GetChannelUsers(channelID uint) ([]models.VoiceState, error) {
	return s.voiceRepo.FindByChannelID(channelID)
}

func (s *VoiceService) UpdateVoiceState(channelID, userID uint, isMuted, isDeaf bool) (*models.VoiceState, error) {
	state, err := s.voiceRepo.FindByChannelAndUser(channelID, userID)
	if err != nil {
		return nil, err
	}

	state.IsMuted = isMuted
	state.IsDeaf = isDeaf

	if err := s.voiceRepo.Update(state); err != nil {
		return nil, err
	}

	return state, nil
}

func (s *VoiceService) GetUserVoiceState(userID uint) ([]models.VoiceState, error) {
	return s.voiceRepo.FindByUserID(userID)
}
