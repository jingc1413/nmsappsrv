package upgrade

// AddShutdownTaskRequest represents the request to create a shutdown task.
type AddShutdownTaskRequest struct {
	Name        string  `json:"name" binding:"required"`
	ExecuteMode int     `json:"executeMode" binding:"required"` // 1=immediate, 2=awaiting, 3=scheduled
	TriggerTime string  `json:"triggerTime,omitempty"`          // RFC3339 for mode 3
	ElementIds  []int64 `json:"elementIds" binding:"required"`
}

// ListShutdownTaskRequest represents the request to list shutdown tasks.
type ListShutdownTaskRequest struct {
	Page     int `json:"page"`
	PageSize int `json:"pageSize"`
}

// ShutdownTaskVo represents a shutdown task in the list response.
type ShutdownTaskVo struct {
	Id            int    `json:"id"`
	Name          string `json:"name"`
	OperationUser string `json:"operationUser"`
	OperationTime string `json:"operationTime"`
	Status        int    `json:"status"`
	ExecuteMode   int    `json:"executeMode"`
	DeviceCount   int    `json:"deviceCount"`
	Progress      string `json:"progress"` // "success/total"
}

// ViewShutdownTaskRequest represents the request to view a shutdown task detail.
type ViewShutdownTaskRequest struct {
	TaskId int `json:"taskId" binding:"required"`
}

// ViewShutdownTaskVo represents the detail of a shutdown task with device list.
type ViewShutdownTaskVo struct {
	Id            int                `json:"id"`
	Name          string             `json:"name"`
	OperationUser string             `json:"operationUser"`
	OperationTime string             `json:"operationTime"`
	Status        int                `json:"status"`
	ExecuteMode   int                `json:"executeMode"`
	TriggerTime   string             `json:"triggerTime"`
	Devices       []ShutdownDeviceVo `json:"devices"`
}

// ShutdownDeviceVo represents a device in the shutdown task detail.
type ShutdownDeviceVo struct {
	ElementId    int64  `json:"elementId"`
	DeviceName   string `json:"deviceName"`
	SerialNumber string `json:"serialNumber"`
}

// ListShutdownResultRequest represents the request to list shutdown results.
type ListShutdownResultRequest struct {
	TaskId   int `json:"taskId" binding:"required"`
	Page     int `json:"page"`
	PageSize int `json:"pageSize"`
}

// ShutdownResultVo represents the shutdown result for a single device.
type ShutdownResultVo struct {
	ElementId    int64  `json:"elementId"`
	DeviceName   string `json:"deviceName"`
	SerialNumber string `json:"serialNumber"`
	Status       int    `json:"status"`
	Time         string `json:"time"`
}

// DeleteShutdownTaskRequest represents the request to delete a shutdown task.
type DeleteShutdownTaskRequest struct {
	TaskId int `json:"taskId" binding:"required"`
}
