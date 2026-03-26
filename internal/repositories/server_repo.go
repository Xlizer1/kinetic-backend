package repositories

import (
	"kinetic-backend/internal/models"
	"math/rand"
	"time"

	"gorm.io/gorm"
)

type ServerRepository struct {
	db *gorm.DB
}

func NewServerRepository(db *gorm.DB) *ServerRepository {
	return &ServerRepository{db: db}
}

func (r *ServerRepository) Create(server *models.Server) error {
	server.InviteCode = generateInviteCode()
	return r.db.Create(server).Error
}

func (r *ServerRepository) FindByID(id uint) (*models.Server, error) {
	var server models.Server
	err := r.db.Preload("Owner").Preload("Members").Preload("Channels").First(&server, id).Error
	if err != nil {
		return nil, err
	}
	return &server, nil
}

func (r *ServerRepository) FindByInviteCode(code string) (*models.Server, error) {
	var server models.Server
	err := r.db.Where("invite_code = ?", code).First(&server).Error
	if err != nil {
		return nil, err
	}
	return &server, nil
}

func (r *ServerRepository) FindByUserID(userID uint) ([]models.Server, error) {
	var servers []models.Server
	err := r.db.Joins("JOIN server_members ON server_members.server_id = servers.id").
		Where("server_members.user_id = ?", userID).
		Preload("Owner").
		Find(&servers).Error
	return servers, err
}

func (r *ServerRepository) Update(server *models.Server) error {
	return r.db.Save(server).Error
}

func (r *ServerRepository) Delete(id uint) error {
	return r.db.Delete(&models.Server{}, id).Error
}

func (r *ServerRepository) AddMember(member *models.ServerMember) error {
	return r.db.Create(member).Error
}

func (r *ServerRepository) RemoveMember(userID, serverID uint) error {
	return r.db.Where("user_id = ? AND server_id = ?", userID, serverID).Delete(&models.ServerMember{}).Error
}

func (r *ServerRepository) FindMember(userID, serverID uint) (*models.ServerMember, error) {
	var member models.ServerMember
	err := r.db.Where("user_id = ? AND server_id = ?", userID, serverID).First(&member).Error
	if err != nil {
		return nil, err
	}
	return &member, nil
}

func generateInviteCode() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	rand.Seed(time.Now().UnixNano())
	code := make([]byte, 8)
	for i := range code {
		code[i] = charset[rand.Intn(len(charset))]
	}
	return string(code)
}
