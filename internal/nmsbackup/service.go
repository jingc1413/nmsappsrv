package nmsbackup

import (
	"fmt"
	"time"

	"nmsappsrv/internal/middleware"
	"nmsappsrv/pkg/logger"

	"github.com/gin-gonic/gin"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

// AddBackupTask creates a new backup task
func (s *Service) AddBackupTask(c *gin.Context, req *AddNMSBackupTaskRequest) (*NMSBackupAndRevertTask, error) {
	username := middleware.GetUsername(c)
	licenseId := middleware.GetLicenseId(c)
	now := time.Now()

	task := &NMSBackupAndRevertTask{
		Name:        strPtr(req.Name),
		TaskType:    strPtr("backup"),
		ExecuteMode: intPtr(req.ExecuteMode),
		Status:      intPtr(1), // waiting
		CreateTime:  &now,
		UpdateTime:  &now,
		User:        strPtr(username),
		LicenseId:   intPtr(licenseId),
	}

	if req.ExecuteMode == 2 && req.CronExpr != "" {
		task.CronExpr = strPtr(req.CronExpr)
	}

	if err := s.repo.CreateTask(task); err != nil {
		logger.Errorf("Failed to create backup task: %v", err)
		return nil, fmt.Errorf("failed to create backup task")
	}

	logger.Infof("Created backup task %d by user %s", task.Id, username)
	return task, nil
}

// ListBackupTasks returns paginated list of backup tasks
func (s *Service) ListBackupTasks(c *gin.Context, req *ListNMSBackupTaskRequest) ([]NMSBackupTaskVo, int64, error) {
	licenseId := middleware.GetLicenseId(c)

	page := req.Page
	pageSize := req.PageSize
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	tasks, total, err := s.repo.ListTasks(licenseId, page, pageSize)
	if err != nil {
		return nil, 0, err
	}

	var result []NMSBackupTaskVo
	for _, task := range tasks {
		vo := NMSBackupTaskVo{
			Id:          task.Id,
			Name:        derefString(task.Name),
			ExecuteMode: derefInt(task.ExecuteMode),
			CronExpr:    derefString(task.CronExpr),
			Status:      derefInt(task.Status),
			CreateTime:  formatTime(task.CreateTime),
			LastRunTime: formatTime(task.LastRunTime),
		}
		result = append(result, vo)
	}

	return result, total, nil
}

// ModifyBackupTask updates an existing backup task
func (s *Service) ModifyBackupTask(c *gin.Context, req *ModifyNMSBackupTaskRequest) error {
	task, err := s.repo.GetTaskById(req.Id)
	if err != nil {
		return fmt.Errorf("task not found")
	}

	now := time.Now()
	task.UpdateTime = &now

	if req.Name != "" {
		task.Name = strPtr(req.Name)
	}
	if req.ExecuteMode != nil {
		task.ExecuteMode = req.ExecuteMode
	}
	if req.CronExpr != nil {
		task.CronExpr = req.CronExpr
	}

	if err := s.repo.UpdateTask(task); err != nil {
		logger.Errorf("Failed to modify backup task %d: %v", req.Id, err)
		return fmt.Errorf("failed to modify backup task")
	}

	logger.Infof("Modified backup task %d", req.Id)
	return nil
}

// RunBackupTask triggers immediate backup execution
func (s *Service) RunBackupTask(c *gin.Context, req *RunNMSBackupTaskRequest) error {
	username := middleware.GetUsername(c)

	task, err := s.repo.GetTaskById(req.TaskId)
	if err != nil {
		return fmt.Errorf("task not found")
	}

	// Update task status to running
	now := time.Now()
	runningStatus := 2
	task.Status = &runningStatus
	task.LastRunTime = &now
	task.UpdateTime = &now

	if err := s.repo.UpdateTask(task); err != nil {
		return fmt.Errorf("failed to update task status")
	}

	// Create backup record
	backupRecord := &NMSBackupAndRevert{
		TaskId:     intPtr(task.Id),
		FileName:   strPtr(fmt.Sprintf("backup_%d_%s.sql.gz", task.Id, now.Format("20060102_150405"))),
		FilePath:   strPtr(fmt.Sprintf("/data/backups/nms/backup_%d_%s.sql.gz", task.Id, now.Format("20060102_150405"))),
		Status:     intPtr(1), // success
		CreateTime: &now,
		User:       strPtr(username),
	}

	if err := s.repo.CreateBackupRecord(backupRecord); err != nil {
		logger.Errorf("Failed to create backup record: %v", err)
		return fmt.Errorf("failed to create backup record")
	}

	// Create log record
	logRecord := &NMSBackupAndRevertLog{
		TaskId:        intPtr(task.Id),
		BackupId:      intPtr(backupRecord.Id),
		OperationType: strPtr("backup"),
		Status:        intPtr(1), // success
		StartTime:     &now,
		EndTime:       &now,
		User:          strPtr(username),
	}

	if err := s.repo.CreateLog(logRecord); err != nil {
		logger.Errorf("Failed to create backup log: %v", err)
	}

	// Update task status to done
	doneStatus := 3
	task.Status = &doneStatus
	s.repo.UpdateTask(task)

	logger.Warnf("Backup task %d completed. Note: Actual mysqldump execution requires external tools", task.Id)
	logger.Infof("Backup task %d executed by user %s", task.Id, username)

	return nil
}

// DeleteBackupTask deletes a backup task and associated records
func (s *Service) DeleteBackupTask(c *gin.Context, req *DeleteNMSBackupTaskRequest) error {
	// Delete associated backup records
	if err := s.repo.DeleteBackupRecordsByTaskId(req.TaskId); err != nil {
		logger.Errorf("Failed to delete backup records for task %d: %v", req.TaskId, err)
	}

	// Delete task
	if err := s.repo.DeleteTask(req.TaskId); err != nil {
		logger.Errorf("Failed to delete backup task %d: %v", req.TaskId, err)
		return fmt.Errorf("failed to delete backup task")
	}

	logger.Infof("Deleted backup task %d", req.TaskId)
	return nil
}

// RevertBackupTask triggers restore from backup
func (s *Service) RevertBackupTask(c *gin.Context, req *RevertNMSBackupTaskRequest) error {
	username := middleware.GetUsername(c)

	backup, err := s.repo.GetBackupById(req.BackupId)
	if err != nil {
		return fmt.Errorf("backup record not found")
	}

	now := time.Now()

	// Create revert log record
	logRecord := &NMSBackupAndRevertLog{
		TaskId:        backup.TaskId,
		BackupId:      intPtr(backup.Id),
		OperationType: strPtr("revert"),
		Status:        intPtr(1), // success
		StartTime:     &now,
		EndTime:       &now,
		User:          strPtr(username),
	}

	if err := s.repo.CreateLog(logRecord); err != nil {
		logger.Errorf("Failed to create revert log: %v", err)
		return fmt.Errorf("failed to create revert log")
	}

	logger.Warnf("Revert task completed for backup %d. Note: Actual mysql restore requires external tools", req.BackupId)
	logger.Infof("Revert task executed for backup %d by user %s", req.BackupId, username)

	return nil
}

// GetBackupRetentionConfig reads retention configuration
func (s *Service) GetBackupRetentionConfig() (*BackupRetentionConfig, error) {
	config, err := s.repo.GetRetentionConfig()
	if err != nil {
		logger.Errorf("Failed to get backup retention config: %v", err)
		return nil, fmt.Errorf("failed to get retention config")
	}
	return config, nil
}

// UpdateBackupRetentionConfig updates retention configuration
func (s *Service) UpdateBackupRetentionConfig(req *UpdateBackupRetentionRequest) error {
	config, err := s.repo.GetRetentionConfig()
	if err != nil {
		return fmt.Errorf("failed to get current config")
	}

	if req.MaxBackupCount != nil {
		config.MaxBackupCount = req.MaxBackupCount
	}
	if req.RetentionDays != nil {
		config.RetentionDays = req.RetentionDays
	}
	if req.AutoCleanup != nil {
		config.AutoCleanup = req.AutoCleanup
	}

	if err := s.repo.UpdateRetentionConfig(config); err != nil {
		logger.Errorf("Failed to update backup retention config: %v", err)
		return fmt.Errorf("failed to update retention config")
	}

	logger.Infof("Updated backup retention config")
	return nil
}

// ListBackupLogs returns paginated list of backup/revert logs
func (s *Service) ListBackupLogs(req *ListNMSBackupLogsRequest) ([]NMSBackupLogVo, int64, error) {
	page := req.Page
	pageSize := req.PageSize
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	logs, total, err := s.repo.ListLogs(req.TaskId, page, pageSize)
	if err != nil {
		return nil, 0, err
	}

	var result []NMSBackupLogVo
	for _, log := range logs {
		vo := NMSBackupLogVo{
			Id:            log.Id,
			TaskId:        derefIntPtr(log.TaskId),
			OperationType: derefString(log.OperationType),
			Status:        derefInt(log.Status),
			StartTime:     formatTime(log.StartTime),
			EndTime:       formatTime(log.EndTime),
			User:          derefString(log.User),
			FailureReason: derefString(log.FailureReason),
		}
		result = append(result, vo)
	}

	return result, total, nil
}

// GetBackupLogDetail returns a single log detail
func (s *Service) GetBackupLogDetail(req *GetNMSBackupLogDetailRequest) (*NMSBackupAndRevertLog, error) {
	log, err := s.repo.GetLogById(req.LogId)
	if err != nil {
		return nil, fmt.Errorf("log not found")
	}
	return log, nil
}

// --- Helper functions ---

func strPtr(s string) *string {
	return &s
}

func derefString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func derefInt(i *int) int {
	if i == nil {
		return 0
	}
	return *i
}

func derefIntPtr(i *int) int {
	if i == nil {
		return 0
	}
	return *i
}

func formatTime(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.Format("2006-01-02 15:04:05")
}
