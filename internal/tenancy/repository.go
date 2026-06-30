package tenancy

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

// Repository provides database operations for tenancy management
type Repository struct {
	db *gorm.DB
}

// NewRepository creates a new Repository
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// Create inserts a new tenancy record
func (r *Repository) Create(t *tenancyModel) error {
	return r.db.Create(t).Error
}

// Update updates an existing tenancy record
func (r *Repository) Update(t *tenancyModel) error {
	return r.db.Save(t).Error
}

// FindByID returns a tenancy by ID
func (r *Repository) FindByID(id int) (*tenancyModel, error) {
	var t tenancyModel
	if err := r.db.Where("id = ?", id).First(&t).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &t, nil
}

// DeleteByID deletes a tenancy by ID
func (r *Repository) DeleteByID(id int) error {
	return r.db.Delete(&tenancyModel{}, "id = ?", id).Error
}

// ExistsByName checks if a tenancy with the given name already exists
func (r *Repository) ExistsByName(name string) (bool, error) {
	var count int64
	err := r.db.Model(&tenancyModel{}).Where("license_name = ?", name).Count(&count).Error
	return count > 0, err
}

// ExistsByNameExcluding checks if a tenancy with the given name exists, excluding a specific ID
func (r *Repository) ExistsByNameExcluding(name string, excludeID int) (bool, error) {
	var count int64
	err := r.db.Model(&tenancyModel{}).Where("license_name = ? AND id != ?", name, excludeID).Count(&count).Error
	return count > 0, err
}

// List returns paginated tenancies with optional name filter
func (r *Repository) List(nameFilter string, page, pageSize int) ([]tenancyModel, int64, error) {
	var items []tenancyModel
	var total int64

	query := r.db.Model(&tenancyModel{})
	if nameFilter != "" {
		query = query.Where("license_name LIKE ?", fmt.Sprintf("%%%s%%", nameFilter))
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("id ASC").Find(&items).Error; err != nil {
		return nil, 0, err
	}

	return items, total, nil
}

// strPtr returns a pointer to the given string
func strPtr(s string) *string {
	return &s
}

// timeFromMillis converts a millisecond timestamp to time.Time
func timeFromMillis(ms int64) *time.Time {
	t := time.UnixMilli(ms)
	return &t
}

// strOrEmpty safely dereferences a string pointer
func strOrEmpty(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// millisFromTime converts a time.Time to milliseconds
func millisFromTime(t *time.Time) int64 {
	if t == nil {
		return 0
	}
	return t.UnixMilli()
}
