package services

import (
	"errors"

	"kinetic-backend/internal/models"
	"kinetic-backend/internal/repositories"
)

type UserService struct {
	userRepo *repositories.UserRepository
}

func NewUserService(userRepo *repositories.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

type UpdateUserInput struct {
	Username  string `json:"username"`
	AvatarURL string `json:"avatar_url"`
}

type UpdateSettingsInput struct {
	Theme         string `json:"theme"`
	Language      string `json:"language"`
	Notifications *bool  `json:"notifications"`
	Privacy       string `json:"privacy"`
}

func (s *UserService) GetUserByID(id uint) (*models.User, error) {
	return s.userRepo.FindByID(id)
}

func (s *UserService) UpdateUser(id uint, input UpdateUserInput) (*models.User, error) {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if input.Username != "" {
		existing, _ := s.userRepo.FindByUsername(input.Username)
		if existing != nil && existing.ID != id {
			return nil, errors.New("username already taken")
		}
		user.Username = input.Username
	}

	if input.AvatarURL != "" {
		user.AvatarURL = input.AvatarURL
	}

	if err := s.userRepo.Update(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) UpdateUserSettings(userID uint, input UpdateSettingsInput) (*models.User, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}

	settings := user.GetSettings()

	if input.Theme != "" {
		if input.Theme == "light" || input.Theme == "dark" {
			settings.Theme = input.Theme
		}
	}

	if input.Language != "" {
		settings.Language = input.Language
	}

	if input.Notifications != nil {
		settings.Notifications = *input.Notifications
	}

	if input.Privacy != "" {
		if input.Privacy == "everyone" || input.Privacy == "friends" || input.Privacy == "none" {
			settings.Privacy = input.Privacy
		}
	}

	user.SetSettings(settings)

	if err := s.userRepo.Update(user); err != nil {
		return nil, err
	}

	return user, nil
}
