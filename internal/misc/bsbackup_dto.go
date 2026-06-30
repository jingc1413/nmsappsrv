package misc

// ---------------------------------------------------------------------------
// DTOs for BaseStation Backup & Restore module
// These extend the existing misc package without redefining any existing types.
// ---------------------------------------------------------------------------

// BaseStationBackupInfoVo is the response VO for listBaseStationBackupLatestFileInfo.
type BaseStationBackupInfoVo struct {
	ElementId        int64  `json:"elementId"`
	DeviceName       string `json:"deviceName"`
	SerialNumber     string `json:"serialNumber"`
	ConfigFile       string `json:"configFile"`
	ConfigFileTime   string `json:"configFileTime"`
	HasBackup        bool   `json:"hasBackup"`
	LatestBackupTime string `json:"latestBackupTime"`
}

// ListBaseStationBackupInfoRequest is the request body for listing device backup info.
type ListBaseStationBackupInfoRequest struct {
	SearchText string `json:"searchText"`
	Page       int    `json:"page"`
	PageSize   int    `json:"pageSize"`
}

// ImportConfigFileResult is returned after a successful config file import.
type ImportConfigFileResult struct {
	ElementId int64  `json:"elementId"`
	FileName  string `json:"fileName"`
	Success   bool   `json:"success"`
	Message   string `json:"message,omitempty"`
}

// ExportConfigFileRequest is the request body for exporting config files as a zip.
type ExportConfigFileRequest struct {
	ElementIds []int64 `json:"elementIds" binding:"required"`
}

// AddBSBackupTaskRequest is the request body for creating a device-specific backup task.
type AddBSBackupTaskRequest struct {
	Name               string   `json:"name" binding:"required"`
	ExecuteMode        int      `json:"executeMode" binding:"required"` // 1=immediate, 2=awaiting, 3=scheduled
	TriggerTime        string   `json:"triggerTime,omitempty"`
	ElementIds         []int64  `json:"elementIds"`
	ExecuteOnAllDevice bool     `json:"executeOnAllDevice"`
	Scope              string   `json:"scope"`
	DeviceGroupIds     []string `json:"deviceGroupIds"`
}

// CancelTaskRequest is the request body for cancelling a backup or restore task.
type CancelTaskRequest struct {
	TaskId int `json:"taskId" binding:"required"`
}

// AddBSRestoreTaskRequest is the request body for creating a device-specific restore task.
type AddBSRestoreTaskRequest struct {
	Name               string   `json:"name" binding:"required"`
	ExecuteMode        int      `json:"executeMode" binding:"required"` // 1=immediate, 2=awaiting, 3=scheduled
	TriggerTime        string   `json:"triggerTime,omitempty"`
	ElementIds         []int64  `json:"elementIds"`
	ExecuteOnAllDevice bool     `json:"executeOnAllDevice"`
	Scope              string   `json:"scope"`
	DeviceGroupIds     []string `json:"deviceGroupIds"`
}

// StartTaskRequest is the request body for manually starting an awaiting task.
type StartTaskRequest struct {
	TaskId int `json:"taskId" binding:"required"`
}

// ListBSBackupTaskRequest is the request body for paginated task listing.
type ListBSBackupTaskRequest struct {
	Page     int `json:"page"`
	PageSize int `json:"pageSize"`
}

// ListDeviceBackupResultRequest is the request body for per-device result listing.
type ListDeviceBackupResultRequest struct {
	TaskId   int `json:"taskId" binding:"required"`
	Page     int `json:"page"`
	PageSize int `json:"pageSize"`
}

// DeviceBackupResultVo is the per-device result VO for a backup/restore task.
type DeviceBackupResultVo struct {
	ElementId         int64  `json:"elementId"`
	DeviceName        string `json:"deviceName"`
	SerialNumber      string `json:"serialNumber"`
	Result            *int   `json:"result"` // null=pending, 1=success, 2=failure
	FailureReason     string `json:"failureReason"`
	StartTime         string `json:"startTime"`
	EndTime           string `json:"endTime"`
	ConfigurationFile string `json:"configurationFile"`
}
