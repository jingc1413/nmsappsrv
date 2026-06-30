package platform

import (
	"errors"

	"gorm.io/gorm"
)

// Repository provides database operations for platform settings
type Repository struct {
	db *gorm.DB
}

// NewRepository creates a new Repository
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// GetSystemConfig reads a system_config entry by config_key
func (r *Repository) GetSystemConfig(key string) (string, error) {
	var cfg systemConfigModel
	if err := r.db.Where("config_key = ?", key).First(&cfg).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", nil
		}
		return "", err
	}
	if cfg.Value == nil {
		return "", nil
	}
	return *cfg.Value, nil
}

// SaveSystemConfig upserts a system_config entry
func (r *Repository) SaveSystemConfig(key, value string) error {
	var cfg systemConfigModel
	err := r.db.Where("config_key = ?", key).First(&cfg).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			cfg = systemConfigModel{
				Key:   &key,
				Value: &value,
			}
			return r.db.Create(&cfg).Error
		}
		return err
	}
	cfg.Value = &value
	return r.db.Save(&cfg).Error
}
