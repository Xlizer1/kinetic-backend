package services

import (
	"errors"
	"time"

	"kinetic-backend/internal/models"
	"kinetic-backend/internal/repositories"
)

type ServerService struct {
	serverRepo *repositories.ServerRepository
}

func NewServerService(serverRepo *repositories.ServerRepository) *ServerService {
	return &ServerService{serverRepo: serverRepo}
}

type CreateServerInput struct {
	Name    string `json:"name" binding:"required"`
	IconURL string `json:"icon_url"`
}

type UpdateServerInput struct {
	Name    string `json:"name"`
	IconURL string `json:"icon_url"`
}

func (s *ServerService) CreateServer(ownerID uint, input CreateServerInput) (*models.Server, error) {
	server := &models.Server{
		OwnerID: ownerID,
		Name:    input.Name,
		IconURL: input.IconURL,
	}

	if err := s.serverRepo.Create(server); err != nil {
		return nil, err
	}

	member := &models.ServerMember{
		UserID:   ownerID,
		ServerID: server.ID,
		Role:     "owner",
		JoinedAt: time.Now(),
	}

	if err := s.serverRepo.AddMember(member); err != nil {
		return nil, err
	}

	return s.serverRepo.FindByID(server.ID)
}

func (s *ServerService) GetServerByID(id uint) (*models.Server, error) {
	return s.serverRepo.FindByID(id)
}

func (s *ServerService) GetUserServers(userID uint) ([]models.Server, error) {
	return s.serverRepo.FindByUserID(userID)
}

func (s *ServerService) UpdateServer(id uint, input UpdateServerInput) (*models.Server, error) {
	server, err := s.serverRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if input.Name != "" {
		server.Name = input.Name
	}

	if input.IconURL != "" {
		server.IconURL = input.IconURL
	}

	if err := s.serverRepo.Update(server); err != nil {
		return nil, err
	}

	return s.serverRepo.FindByID(id)
}

func (s *ServerService) DeleteServer(id uint) error {
	return s.serverRepo.Delete(id)
}

func (s *ServerService) JoinServerByInviteCode(userID uint, inviteCode string) (*models.Server, error) {
	server, err := s.serverRepo.FindByInviteCode(inviteCode)
	if err != nil {
		return nil, errors.New("invalid invite code")
	}

	existingMember, _ := s.serverRepo.FindMember(userID, server.ID)
	if existingMember != nil {
		return nil, errors.New("already a member of this server")
	}

	member := &models.ServerMember{
		UserID:   userID,
		ServerID: server.ID,
		Role:     "member",
		JoinedAt: time.Now(),
	}

	if err := s.serverRepo.AddMember(member); err != nil {
		return nil, err
	}

	return s.serverRepo.FindByID(server.ID)
}

func (s *ServerService) LeaveServer(userID, serverID uint) error {
	member, err := s.serverRepo.FindMember(userID, serverID)
	if err != nil {
		return errors.New("not a member of this server")
	}

	if member.Role == "owner" {
		return errors.New("owner cannot leave server, transfer ownership first")
	}

	return s.serverRepo.RemoveMember(userID, serverID)
}
