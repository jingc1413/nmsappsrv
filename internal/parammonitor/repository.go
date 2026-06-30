package parammonitor

import (
	"nmsappsrv/pkg/logger"

	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateConfig(config *ParameterMonitorConfig) error {
	return r.db.Create(config).Error
}

func (r *Repository) UpdateConfig(config *ParameterMonitorConfig) error {
	return r.db.Save(config).Error
}

func (r *Repository) DeleteConfig(id int) error {
	return r.db.Delete(&ParameterMonitorConfig{}, id).Error
}

func (r *Repository) GetConfig(id int) (*ParameterMonitorConfig, error) {
	var config ParameterMonitorConfig
	err := r.db.First(&config, id).Error
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func (r *Repository) ListConfigs(licenseId int, page, pageSize int) ([]ParameterMonitorConfig, int64, error) {
	var configs []ParameterMonitorConfig
	var total int64

	query := r.db.Where("license_id = ?", licenseId)
	query.Model(&ParameterMonitorConfig{}).Count(&total)

	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).Order("id DESC").Find(&configs).Error
	if err != nil {
		logger.Errorf("ListConfigs error: %v", err)
		return nil, 0, err
	}

	return configs, total, nil
}

func (r *Repository) SetConfigParameters(configId int, parameterIds []string) error {
	// Delete old associations
	err := r.db.Where("config_id = ?", configId).Delete(&MonitorConfigHasParameter{}).Error
	if err != nil {
		return err
	}

	// Insert new associations
	for _, paramId := range parameterIds {
		assoc := MonitorConfigHasParameter{
			ConfigId:    &configId,
			ParameterId: &paramId,
		}
		if err := r.db.Create(&assoc).Error; err != nil {
			return err
		}
	}

	return nil
}

func (r *Repository) GetConfigParameters(configId int) ([]string, error) {
	var assocs []MonitorConfigHasParameter
	err := r.db.Where("config_id = ?", configId).Find(&assocs).Error
	if err != nil {
		return nil, err
	}

	paramIds := make([]string, 0, len(assocs))
	for _, assoc := range assocs {
		if assoc.ParameterId != nil {
			paramIds = append(paramIds, *assoc.ParameterId)
		}
	}

	return paramIds, nil
}

func (r *Repository) GetParameterByIds(ids []string) (map[string]string, error) {
	if len(ids) == 0 {
		return make(map[string]string), nil
	}

	var params []struct {
		Id   string
		Path *string
	}
	err := r.db.Table("parameter").Select("id, path").Where("id IN ?", ids).Find(&params).Error
	if err != nil {
		return nil, err
	}

	result := make(map[string]string)
	for _, p := range params {
		if p.Path != nil {
			result[p.Id] = *p.Path
		}
	}

	return result, nil
}
