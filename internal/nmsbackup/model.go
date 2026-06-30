package nmsbackup

import "time"

// NMSBackupAndRevertTask 对应 nms_backup_and_revert_task 表
type NMSBackupAndRevertTask struct {
	Id          int        `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        *string    `gorm:"column:name;type:varchar(255)" json:"name"`
	TaskType    *string    `gorm:"column:task_type;type:varchar(255)" json:"task_type"` // "backup" or "revert"
	ExecuteMode *int       `gorm:"column:execute_mode" json:"execute_mode"`             // 1=once, 2=scheduled
	CronExpr    *string    `gorm:"column:cron_expr;type:varchar(255)" json:"cron_expr"`
	Status      *int       `gorm:"column:status" json:"status"` // 1=waiting, 2=running, 3=done, 4=cancelled
	LastRunTime *time.Time `gorm:"column:last_run_time" json:"last_run_time"`
	CreateTime  *time.Time `gorm:"column:create_time" json:"create_time"`
	UpdateTime  *time.Time `gorm:"column:update_time" json:"update_time"`
	User        *string    `gorm:"column:user;type:varchar(255)" json:"user"`
	LicenseId   *int       `gorm:"column:license_id" json:"license_id"`
}

func (NMSBackupAndRevertTask) TableName() string { return "nms_backup_and_revert_task" }

// NMSBackupAndRevert 对应 nms_backup_and_revert 表 (backup file records)
type NMSBackupAndRevert struct {
	Id         int        `gorm:"primaryKey;autoIncrement" json:"id"`
	TaskId     *int       `gorm:"column:task_id" json:"task_id"`
	FileName   *string    `gorm:"column:file_name;type:varchar(255)" json:"file_name"`
	FilePath   *string    `gorm:"column:file_path;type:varchar(255)" json:"file_path"`
	FileSize   *int64     `gorm:"column:file_size" json:"file_size"`
	Status     *int       `gorm:"column:status" json:"status"` // 1=success, 2=failed
	CreateTime *time.Time `gorm:"column:create_time" json:"create_time"`
	User       *string    `gorm:"column:user;type:varchar(255)" json:"user"`
}

func (NMSBackupAndRevert) TableName() string { return "nms_backup_and_revert" }

// NMSBackupAndRevertLog 对应 nms_backup_and_revert_log 表
type NMSBackupAndRevertLog struct {
	Id            int        `gorm:"primaryKey;autoIncrement" json:"id"`
	TaskId        *int       `gorm:"column:task_id" json:"task_id"`
	BackupId      *int       `gorm:"column:backup_id" json:"backup_id"`
	OperationType *string    `gorm:"column:operation_type;type:varchar(255)" json:"operation_type"` // "backup" or "revert"
	Status        *int       `gorm:"column:status" json:"status"`
	StartTime     *time.Time `gorm:"column:start_time" json:"start_time"`
	EndTime       *time.Time `gorm:"column:end_time" json:"end_time"`
	User          *string    `gorm:"column:user;type:varchar(255)" json:"user"`
	FailureReason *string    `gorm:"column:failure_reason;type:text" json:"failure_reason"`
}

func (NMSBackupAndRevertLog) TableName() string { return "nms_backup_and_revert_log" }

// BackupRetentionConfig stored in system_config as key "nms_backup_retention"
type BackupRetentionConfig struct {
	MaxBackupCount *int  `json:"maxBackupCount"`
	RetentionDays  *int  `json:"retentionDays"`
	AutoCleanup    *bool `json:"autoCleanup"`
}

// --- DTOs ---

type AddNMSBackupTaskRequest struct {
	Name        string `json:"name" binding:"required"`
	ExecuteMode int    `json:"executeMode"` // 1=once, 2=scheduled
	CronExpr    string `json:"cronExpr,omitempty"`
}

type ModifyNMSBackupTaskRequest struct {
	Id          int     `json:"id" binding:"required"`
	Name        string  `json:"name"`
	ExecuteMode *int    `json:"executeMode"`
	CronExpr    *string `json:"cronExpr"`
}

type RunNMSBackupTaskRequest struct {
	TaskId int `json:"taskId" binding:"required"`
}

type DeleteNMSBackupTaskRequest struct {
	TaskId int `json:"taskId" binding:"required"`
}

type RevertNMSBackupTaskRequest struct {
	BackupId int `json:"backupId" binding:"required"`
}

type ListNMSBackupTaskRequest struct {
	Page     int `json:"page"`
	PageSize int `json:"pageSize"`
}

type NMSBackupTaskVo struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	ExecuteMode int    `json:"executeMode"`
	CronExpr    string `json:"cronExpr"`
	Status      int    `json:"status"`
	LastRunTime string `json:"lastRunTime"`
	CreateTime  string `json:"createTime"`
}

type ListNMSBackupLogsRequest struct {
	TaskId   int `json:"taskId"`
	Page     int `json:"page"`
	PageSize int `json:"pageSize"`
}

type NMSBackupLogVo struct {
	Id            int    `json:"id"`
	TaskId        int    `json:"taskId"`
	OperationType string `json:"operationType"`
	Status        int    `json:"status"`
	StartTime     string `json:"startTime"`
	EndTime       string `json:"endTime"`
	User          string `json:"user"`
	FailureReason string `json:"failureReason"`
}

type GetNMSBackupLogDetailRequest struct {
	LogId int `json:"logId" binding:"required"`
}

type UpdateBackupRetentionRequest struct {
	MaxBackupCount *int  `json:"maxBackupCount"`
	RetentionDays  *int  `json:"retentionDays"`
	AutoCleanup    *bool `json:"autoCleanup"`
}
