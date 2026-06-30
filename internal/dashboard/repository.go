package dashboard

import (
	"context"
	"fmt"
	"strings"

	"gorm.io/gorm"
)

// Repository handles database operations for Dashboard module
type Repository struct {
	db *gorm.DB
}

// NewRepository creates a new Repository
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// ListDevicesByMode returns devices filtered by mode and tenancyId
func (r *Repository) ListDevicesByMode(ctx context.Context, mode string, tenancyId *int) ([]map[string]interface{}, error) {
	var results []map[string]interface{}

	q := r.db.WithContext(ctx).Table("cpe_element").
		Select("model_name, ne_neid").
		Where("deleted = ?", false)

	if tenancyId != nil {
		q = q.Where("license_id = ?", *tenancyId)
	}

	switch mode {
	case "CPE":
		q = q.Where("device_type = ?", "cpe")
	case "eNB":
		q = q.Where("device_type = ? AND generation = ?", "enb", "LTE")
	case "gNB":
		q = q.Where("device_type = ? AND generation = ?", "enb", "NR")
	default:
		return nil, fmt.Errorf("invalid mode: %s", mode)
	}

	err := q.Scan(&results).Error
	return results, err
}

// ListAllDevices returns all non-deleted devices with type info
func (r *Repository) ListAllDevices(ctx context.Context, tenancyId *int) ([]map[string]interface{}, error) {
	var results []map[string]interface{}

	q := r.db.WithContext(ctx).Table("cpe_element").
		Select("ne_neid, device_type, generation").
		Where("deleted = ?", false)

	if tenancyId != nil {
		q = q.Where("license_id = ?", *tenancyId)
	}

	err := q.Scan(&results).Error
	return results, err
}

// ListPDCPTraffic returns PDCP traffic statistics grouped by PLMN
func (r *Repository) ListPDCPTraffic(ctx context.Context, startTime, endTime string, tenancyId *int) ([]map[string]interface{}, error) {
	var results []map[string]interface{}

	sql := `SELECT plmn, SUM(dl_traffic), SUM(ul_traffic) FROM pdcp_traffic WHERE `
	args := []interface{}{}

	if tenancyId != nil {
		sql += "tenancy_id = ? AND "
		args = append(args, *tenancyId)
	} else {
		sql += "tenancy_id IS NULL AND "
	}

	if startTime != "" {
		sql += "statistic_time >= ? AND "
		args = append(args, startTime)
	}
	if endTime != "" {
		sql += "statistic_time < ? AND "
		args = append(args, endTime)
	}

	// Remove trailing " AND "
	sql = strings.TrimSuffix(sql, " AND ")
	sql += " GROUP BY plmn"

	err := r.db.WithContext(ctx).Raw(sql, args...).Scan(&results).Error
	return results, err
}

// GetDeviceByIds returns devices by their IDs
func (r *Repository) GetDeviceByIds(ctx context.Context, ids []int64) ([]map[string]interface{}, error) {
	var results []map[string]interface{}
	err := r.db.WithContext(ctx).Table("cpe_element").
		Select("ne_neid, serial_number").
		Where("ne_neid IN ?", ids).
		Scan(&results).Error
	return results, err
}
