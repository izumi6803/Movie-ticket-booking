package repository

import (
	"cinema-backend/internal/models"

	"gorm.io/gorm"
)

type SettingRepository struct {
	db *gorm.DB
}

func NewSettingRepository(db *gorm.DB) *SettingRepository {
	return &SettingRepository{db: db}
}

func (r *SettingRepository) GetAll() ([]models.SystemSetting, error) {
	var settings []models.SystemSetting
	err := r.db.Find(&settings).Error
	return settings, err
}

func (r *SettingRepository) GetByKey(key string) (*models.SystemSetting, error) {
	var setting models.SystemSetting
	err := r.db.First(&setting, "key = ?", key).Error
	if err != nil {
		return nil, err
	}
	return &setting, nil
}

func (r *SettingRepository) Upsert(setting *models.SystemSetting) error {
	var existing models.SystemSetting
	result := r.db.First(&existing, "key = ?", setting.Key)
	if result.Error != nil {
		return r.db.Create(setting).Error
	}
	return r.db.Model(&models.SystemSetting{}).Where("key = ?", setting.Key).Update("value", setting.Value).Error
}
