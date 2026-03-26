package services

import (
	"kinetic-backend/internal/models"
	"kinetic-backend/internal/repositories"
)

type PresenceService struct {
	presenceRepo *repositories.PresenceRepository
	userRepo     *repositories.UserRepository
}

func NewPresenceService(presenceRepo *repositories.PresenceRepository, userRepo *repositories.UserRepository) *PresenceService {
	return &PresenceService{
		presenceRepo: presenceRepo,
		userRepo:     userRepo,
	}
}

type PresenceUser struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Status   string `json:"status"`
}

func (s *PresenceService) SetPresence(userID uint, status string) error {
	if !s.isValidStatus(status) {
		status = models.PresenceOnline
	}
	return s.presenceRepo.Upsert(userID, status)
}

func (s *PresenceService) GetPresence(userID uint) (*PresenceUser, error) {
	presence, err := s.presenceRepo.FindByUserID(userID)
	if err != nil {
		return nil, err
	}

	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}

	return &PresenceUser{
		UserID:   presence.UserID,
		Username: user.Username,
		Status:   presence.Status,
	}, nil
}

func (s *PresenceService) GetAllOnline() ([]PresenceUser, error) {
	presences, err := s.presenceRepo.FindAllOnline()
	if err != nil {
		return nil, err
	}

	users := make([]PresenceUser, len(presences))
	for i, p := range presences {
		user, _ := s.userRepo.FindByID(p.UserID)
		username := ""
		if user != nil {
			username = user.Username
		}
		users[i] = PresenceUser{
			UserID:   p.UserID,
			Username: username,
			Status:   p.Status,
		}
	}

	return users, nil
}

func (s *PresenceService) SetOnline(userID uint) error {
	return s.SetPresence(userID, models.PresenceOnline)
}

func (s *PresenceService) SetOffline(userID uint) error {
	return s.presenceRepo.SetOffline(userID)
}

func (s *PresenceService) isValidStatus(status string) bool {
	switch status {
	case models.PresenceOnline, models.PresenceIdle, models.PresenceDND, models.PresenceOffline:
		return true
	}
	return false
}
