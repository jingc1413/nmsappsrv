package ssh

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"gorm.io/gorm"

	"nmsappsrv/pkg/logger"
	"nmsappsrv/pkg/redis"
)

// Repository provides data access for SSH management.
type Repository struct {
	db *gorm.DB
}

// NewRepository creates a new Repository.
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// ---------- SSH Label ----------

func (r *Repository) FindLabelByName(tenancyId int, name string) (*SSHLabel, error) {
	var label SSHLabel
	err := r.db.Where("license_id = ? AND name = ?", tenancyId, name).First(&label).Error
	if err != nil {
		return nil, err
	}
	return &label, nil
}

func (r *Repository) FindLabelByID(id int) (*SSHLabel, error) {
	var label SSHLabel
	err := r.db.First(&label, id).Error
	if err != nil {
		return nil, err
	}
	return &label, nil
}

func (r *Repository) CreateLabel(label *SSHLabel) error {
	return r.db.Create(label).Error
}

func (r *Repository) DeleteLabel(id int) error {
	return r.db.Delete(&SSHLabel{}, id).Error
}

func (r *Repository) UpdateLabel(label *SSHLabel) error {
	return r.db.Save(label).Error
}

func (r *Repository) ListLabels(tenancyId int) ([]SSHLabel, error) {
	var labels []SSHLabel
	err := r.db.Where("license_id = ?", tenancyId).Find(&labels).Error
	return labels, err
}

// ---------- SSH Access Timer ----------

func (r *Repository) FindTimerByElementId(elementId int64) (*SSHAccessTimerTask, error) {
	var task SSHAccessTimerTask
	err := r.db.Where("element_id = ?", elementId).First(&task).Error
	if err != nil {
		return nil, err
	}
	return &task, nil
}

func (r *Repository) CreateTimer(task *SSHAccessTimerTask) error {
	return r.db.Create(task).Error
}

func (r *Repository) UpdateTimer(task *SSHAccessTimerTask) error {
	return r.db.Save(task).Error
}

func (r *Repository) ListTimers(page, pageSize int, elementId *int64) ([]SSHAccessTimerTask, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}

	query := r.db.Model(&SSHAccessTimerTask{})
	if elementId != nil {
		query = query.Where("element_id = ?", *elementId)
	}

	var total int64
	query.Count(&total)

	var tasks []SSHAccessTimerTask
	err := query.Order("deadline DESC").
		Limit(pageSize).Offset((page - 1) * pageSize).
		Find(&tasks).Error
	return tasks, total, err
}

func (r *Repository) FindExpiredTimers() ([]SSHAccessTimerTask, error) {
	var tasks []SSHAccessTimerTask
	now := time.Now()
	err := r.db.Where("ssh_status = ? AND deadline < ?", "1", now).Find(&tasks).Error
	return tasks, err
}

func (r *Repository) FindElementInfo(elementId int64) (sn string, deviceName string, err error) {
	var row struct {
		SN         string `gorm:"column:serial_number"`
		DeviceName string `gorm:"column:device_name"`
	}
	err = r.db.Table("cpe_element").
		Select("serial_number, device_name").
		Where("ne_neid = ? AND deleted = 0", elementId).
		Scan(&row).Error
	return row.SN, row.DeviceName, err
}

func (r *Repository) FindElementIdsByGroup(groupIds []string) ([]int64, error) {
	if len(groupIds) == 0 {
		return nil, nil
	}
	var ids []int64
	err := r.db.Table("device_group_element_rel").
		Select("DISTINCT element_id").
		Where("group_id IN ?", groupIds).
		Pluck("element_id", &ids).Error
	return ids, err
}

func (r *Repository) GetTenancyName(tenancyId int) string {
	var name string
	r.db.Table("tenancy").Where("id = ?", tenancyId).Pluck("name", &name)
	return name
}

// ---------- Service ----------

// Service contains SSH business logic.
type Service struct {
	repo *Repository
}

// NewService creates a new SSH service.
func NewService(db *gorm.DB) *Service {
	return &Service{repo: NewRepository(db)}
}

// ---------- SSH Label methods ----------

func (s *Service) AddLabel(req *AddSSHLabelRequest, tenancyId int) error {
	existing, _ := s.repo.FindLabelByName(tenancyId, req.Name)
	if existing != nil {
		return fmt.Errorf("label name already exists")
	}
	label := &SSHLabel{
		Name:      &req.Name,
		Content:   &req.Content,
		LicenseId: &tenancyId,
	}
	return s.repo.CreateLabel(label)
}

func (s *Service) DeleteLabel(id int) error {
	return s.repo.DeleteLabel(id)
}

func (s *Service) ListLabels(tenancyId int) ([]SSHLabel, error) {
	return s.repo.ListLabels(tenancyId)
}

func (s *Service) UpdateLabel(req *UpdateSSHLabelRequest, tenancyId int) error {
	existing, _ := s.repo.FindLabelByName(tenancyId, req.Name)
	if existing != nil && existing.Id != req.Id {
		return fmt.Errorf("label name already exists")
	}
	label, err := s.repo.FindLabelByID(req.Id)
	if err != nil {
		return fmt.Errorf("label not found")
	}
	label.Name = &req.Name
	label.Content = &req.Content
	return s.repo.UpdateLabel(label)
}

// ---------- SSH Access Timer methods ----------

func (s *Service) SetAccessTimer(req *SSHAccessTimerRequest, tenancyId int, username string) error {
	elementIds := req.ElementIds
	if req.DeviceGroupIds != nil && len(req.DeviceGroupIds) > 0 {
		groupIds, err := s.repo.FindElementIdsByGroup(req.DeviceGroupIds)
		if err != nil {
			return fmt.Errorf("resolve device groups: %w", err)
		}
		elementIds = append(elementIds, groupIds...)
	}
	if len(elementIds) == 0 {
		return fmt.Errorf("no devices selected")
	}

	deadline := time.Now().Add(time.Duration(req.Deadline) * time.Minute)
	tenancyName := s.repo.GetTenancyName(tenancyId)
	ctx := context.Background()
	now := time.Now()

	for _, eid := range elementIds {
		sn, deviceName, err := s.repo.FindElementInfo(int64(eid))
		if err != nil {
			logger.Errorf("ssh_timer: find element %d: %v", eid, err)
			continue
		}

		existing, _ := s.repo.FindTimerByElementId(int64(eid))
		if existing != nil {
			existing.Deadline = &deadline
			existing.SshStatus = strPtr("1")
			existing.LatestModifyTime = &now
			s.repo.UpdateTimer(existing)
		} else {
			task := &SSHAccessTimerTask{
				TenancyName:      &tenancyName,
				TenancyId:        &tenancyId,
				ElementId:        int64Ptr(int64(eid)),
				SshStatus:        strPtr("1"),
				DeviceName:       &deviceName,
				SerialNumber:     &sn,
				Deadline:         &deadline,
				LatestModifyTime: &now,
			}
			s.repo.CreateTimer(task)
		}

		// Push TR-069 SetParamValue to enable SSH
		msg := operationMessage{
			EventType:      "SetParamValue",
			NeNeid:         int64(eid),
			Operation:      "SetParamValue",
			OperationParam: "Device.SecurityManagement.SshEnable=true",
			OperationUser:  username,
			ExpiredAt:      now.Add(5 * time.Minute).UnixMilli(),
		}
		msgJSON, _ := json.Marshal(msg)
		if err := redis.LPush(ctx, "operation_queue", string(msgJSON)); err != nil {
			logger.Errorf("ssh_timer: push enable SSH for %d: %v", eid, err)
		}
	}
	return nil
}

func (s *Service) ListAccessTimers(req *ListSSHAccessTimerRequest) ([]SSHAccessTimerVO, int64, error) {
	tasks, total, err := s.repo.ListTimers(req.Page, req.PageSize, req.ElementId)
	if err != nil {
		return nil, 0, err
	}

	vos := make([]SSHAccessTimerVO, len(tasks))
	for i, t := range tasks {
		vo := SSHAccessTimerVO{
			Id:          t.Id,
			SshStatus:   strVal(t.SshStatus),
			Deadline:    t.Deadline,
			TenancyName: strVal(t.TenancyName),
		}
		if t.ElementId != nil {
			vo.ElementId = *t.ElementId
		}
		if t.DeviceName != nil {
			vo.DeviceName = *t.DeviceName
		}
		if t.SerialNumber != nil {
			vo.SerialNumber = *t.SerialNumber
		}
		vos[i] = vo
	}
	return vos, total, nil
}

// StartExpiredChecker runs a background goroutine that checks every 2 minutes
// for expired SSH access timers and disables SSH on those devices.
func (s *Service) StartExpiredChecker() {
	go func() {
		ticker := time.NewTicker(2 * time.Minute)
		defer ticker.Stop()
		for range ticker.C {
			s.processExpiredTasks()
		}
	}()
}

func (s *Service) processExpiredTasks() {
	tasks, err := s.repo.FindExpiredTimers()
	if err != nil {
		logger.Errorf("ssh_timer: find expired: %v", err)
		return
	}
	ctx := context.Background()
	now := time.Now()

	for _, t := range tasks {
		eid := int64(0)
		if t.ElementId != nil {
			eid = *t.ElementId
		}

		// Push TR-069 to disable SSH
		msg := operationMessage{
			EventType:      "SetParamValue",
			NeNeid:         eid,
			Operation:      "SetParamValue",
			OperationParam: "Device.SecurityManagement.SshEnable=false",
			OperationUser:  "system",
			ExpiredAt:      now.Add(5 * time.Minute).UnixMilli(),
		}
		msgJSON, _ := json.Marshal(msg)
		if err := redis.LPush(ctx, "operation_queue", string(msgJSON)); err != nil {
			logger.Errorf("ssh_timer: push disable SSH for %d: %v", eid, err)
		}

		disabled := "0"
		t.SshStatus = &disabled
		t.LatestModifyTime = &now
		s.repo.UpdateTimer(&t)
	}
}

// ---------- helpers ----------

func strPtr(s string) *string { return &s }
func int64Ptr(i int64) *int64 { return &i }
func strVal(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
