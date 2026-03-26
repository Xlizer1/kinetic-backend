package repositories

import (
	"kinetic-backend/internal/models"

	"gorm.io/gorm"
)

type VoiceRepository struct {
	db *gorm.DB
}

func NewVoiceRepository(db *gorm.DB) *VoiceRepository {
	return &VoiceRepository{db: db}
}

func (r *VoiceRepository) Create(state *models.VoiceState) error {
	return r.db.Create(state).Error
}

func (r *VoiceRepository) FindByChannelAndUser(channelID, userID uint) (*models.VoiceState, error) {
	var state models.VoiceState
	err := r.db.Where("channel_id = ? AND user_id = ?", channelID, userID).First(&state).Error
	if err != nil {
		return nil, err
	}
	return &state, nil
}

func (r *VoiceRepository) FindByChannelID(channelID uint) ([]models.VoiceState, error) {
	var states []models.VoiceState
	err := r.db.Where("channel_id = ?", channelID).Find(&states).Error
	return states, err
}

func (r *VoiceRepository) FindByUserID(userID uint) ([]models.VoiceState, error) {
	var states []models.VoiceState
	err := r.db.Where("user_id = ?", userID).Find(&states).Error
	return states, err
}

func (r *VoiceRepository) Update(state *models.VoiceState) error {
	return r.db.Save(state).Error
}

func (r *VoiceRepository) Delete(channelID, userID uint) error {
	return r.db.Where("channel_id = ? AND user_id = ?", channelID, userID).Delete(&models.VoiceState{}).Error
}

func (r *VoiceRepository) DeleteByChannelID(channelID uint) error {
	return r.db.Where("channel_id = ?", channelID).Delete(&models.VoiceState{}).Error
}
