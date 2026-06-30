package resources

import (
	"encoding/json"
	"time"
)

// Service provides resource monitoring operations
type Service struct {
	repo *Repository
}

// NewService creates a new Service
func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

// GetCpuAndMemUsage returns current CPU and memory usage
func (s *Service) GetCpuAndMemUsage() ResourcesVO {
	// TODO: integrate actual OS-level CPU/memory sampling
	return ResourcesVO{
		CPU:       0.0,
		Mem:       0.0,
		Timestamp: time.Now().Format(time.RFC3339),
	}
}

// GetTableStatus returns MySQL table sizes
func (s *Service) GetTableStatus() ([]TableStatusVO, error) {
	return s.repo.GetTableStatus()
}

// GetDiskUsage returns disk partition usage
func (s *Service) GetDiskUsage() []DiskUsageVO {
	// TODO: integrate actual disk usage sampling (e.g. syscall.Statfs)
	return []DiskUsageVO{}
}

// GetThreshold returns alarm threshold configuration
func (s *Service) GetThreshold() (*ThresholdConfig, error) {
	value, err := s.repo.GetSystemConfig("sysThreshold")
	if err != nil {
		return nil, err
	}

	defaults := &ThresholdConfig{
		CPU:       60,
		Mem:       60,
		Disk:      80,
		DiskClear: 80,
		Table:     5.0,
	}

	if value == "" {
		return defaults, nil
	}

	var cfg ThresholdConfig
	if err := json.Unmarshal([]byte(value), &cfg); err != nil {
		return defaults, nil
	}
	return &cfg, nil
}

// UpdateThreshold updates alarm threshold configuration
func (s *Service) UpdateThreshold(cfg *ThresholdConfig) error {
	data, err := json.Marshal(cfg)
	if err != nil {
		return err
	}
	return s.repo.SaveSystemConfig("sysThreshold", string(data))
}
