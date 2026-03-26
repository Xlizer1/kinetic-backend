package repositories

import (
	"kinetic-backend/internal/models"

	"gorm.io/gorm"
)

type PresenceRepository struct {
	db *gorm.DB
}

func NewPresenceRepository(db *gorm.DB) *PresenceRepository {
	return &PresenceRepository{db: db}
}

func (r *PresenceRepository) Create(presence *models.Presence) error {
	return r.db.Create(presence).Error
}

func (r *PresenceRepository) FindByUserID(userID uint) (*models.Presence, error) {
	var presence models.Presence
	err := r.db.Where("user_id = ?", userID).First(&presence).Error
	if err != nil {
		return nil, err
	}
	return &presence, nil
}

func (r *PresenceRepository) FindAllOnline() ([]models.Presence, error) {
	var presences []models.Presence
	err := r.db.Where("status != ?", models.PresenceOffline).Find(&presences).Error
	return presences, err
}

func (r *PresenceRepository) Update(presence *models.Presence) error {
	return r.db.Save(presence).Error
}

func (r *PresenceRepository) Upsert(userID uint, status string) error {
	var presence models.Presence
	err := r.db.Where("user_id = ?", userID).First(&presence).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}

	if err == gorm.ErrRecordNotFound {
		presence = models.Presence{
			UserID: userID,
			Status: status,
		}
		return r.db.Create(&presence).Error
	}

	presence.Status = status
	return r.db.Save(&presence).Error
}

func (r *PresenceRepository) SetOffline(userID uint) error {
	return r.Upsert(userID, models.PresenceOffline)
}

func (r *PresenceRepository) Delete(userID uint) error {
	return r.db.Where("user_id = ?", userID).Delete(&models.Presence{}).Error
}
