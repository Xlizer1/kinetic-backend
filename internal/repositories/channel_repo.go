package repositories

import (
	"kinetic-backend/internal/models"

	"gorm.io/gorm"
)

type ChannelRepository struct {
	db *gorm.DB
}

func NewChannelRepository(db *gorm.DB) *ChannelRepository {
	return &ChannelRepository{db: db}
}

func (r *ChannelRepository) Create(channel *models.Channel) error {
	return r.db.Create(channel).Error
}

func (r *ChannelRepository) FindByID(id uint) (*models.Channel, error) {
	var channel models.Channel
	err := r.db.First(&channel, id).Error
	if err != nil {
		return nil, err
	}
	return &channel, nil
}

func (r *ChannelRepository) FindByServerID(serverID uint) ([]models.Channel, error) {
	var channels []models.Channel
	err := r.db.Where("server_id = ?", serverID).Order("position").Find(&channels).Error
	return channels, err
}

func (r *ChannelRepository) Update(channel *models.Channel) error {
	return r.db.Save(channel).Error
}

func (r *ChannelRepository) Delete(id uint) error {
	return r.db.Delete(&models.Channel{}, id).Error
}
