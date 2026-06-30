package diagnostics

import "time"

// ---------- parameter_value table (shared with TR-069 engine) ----------

// ParameterValue 对应 parameter_value 表，存储 CPE 的 TR-069 参数值
type ParameterValue struct {
	ParamId    int64      `gorm:"primaryKey;autoIncrement;column:param_id" json:"param_id"`
	ElementId  *int64     `gorm:"column:element_id;index" json:"element_id"`
	ParamName  *string    `gorm:"column:param_name;type:varchar(512);index" json:"param_name"`
	ParamValue *string    `gorm:"column:param_value;type:longtext" json:"param_value"`
	Type       *string    `gorm:"column:type;type:varchar(255)" json:"type"`
	Writable   *bool      `gorm:"column:writable" json:"writable"`
	IsLeaf     *bool      `gorm:"column:is_leaf" json:"is_leaf"`
	UpdateTime *time.Time `gorm:"column:update_time" json:"update_time"`
}

func (ParameterValue) TableName() string { return "parameter_value" }

// ---------- Request DTOs ----------

// PingRequest Ping 诊断请求
type PingRequest struct {
	Server    string `json:"server" binding:"required"`
	Count     int    `json:"count" binding:"required"`
	ElementId int64  `json:"elementId" binding:"required"`
}

// TraceRouteRequest 路由追踪诊断请求
type TraceRouteRequest struct {
	Server    string `json:"server" binding:"required"`
	ElementId int64  `json:"elementId" binding:"required"`
}

// DownloadRequest 下载诊断请求
type DownloadRequest struct {
	DownloadUrl string `json:"downloadUrl"`
	ElementId   int64  `json:"elementId" binding:"required"`
}

// UploadRequest 上传诊断请求
type UploadRequest struct {
	DownloadUrl string `json:"downloadUrl"`
	FileSize    *int64 `json:"fileSize"` // MB, default 10
	ElementId   int64  `json:"elementId" binding:"required"`
}

// IdRequest 通用 ID 请求 (用于查询结果/状态)
type IdRequest struct {
	Id int64 `json:"id" binding:"required"`
}

// ---------- Response VOs ----------

// DiagnosticsResultVO 诊断结果聚合
type DiagnosticsResultVO struct {
	Ping     *PingResultVO     `json:"diagnosticsPingVO"`
	Download *DownloadResultVO `json:"diagnosticsDownloadResultVO"`
	Upload   *UploadResultVO   `json:"diagnosticsUploadResultVO"`
	Trace    *TraceResultVO    `json:"diagnosticsTraceRoutResultVO"`
}

// PingResultVO Ping 诊断结果
type PingResultVO struct {
	TestTime            int64  `json:"testTime"`
	Server              string `json:"server"`
	Status              string `json:"status"`
	SuccessCount        *int   `json:"successCount"`
	FailureCount        *int   `json:"failureCount"`
	AverageResponseTime *int   `json:"averageResponseTime"`
}

// DownloadResultVO 下载诊断结果
type DownloadResultVO struct {
	TestTime        int64    `json:"testTime"`
	Server          string   `json:"server"`
	Statues         string   `json:"statues"`
	ConnectionTimes *int64   `json:"connectionTimes"` // ms
	DownloadTimes   *int64   `json:"downloadTimes"`   // ms
	FileSize        *float64 `json:"fileSize"`         // MB
	SpeedSize       *float64 `json:"speedSize"`        // Mbps
}

// UploadResultVO 上传诊断结果
type UploadResultVO struct {
	TestTime        int64    `json:"testTime"`
	Server          string   `json:"server"`
	Statues         string   `json:"statues"`
	ConnectionTimes *int64   `json:"connectionTimes"` // ms
	UploadTimes     *int64   `json:"uploadTimes"`     // ms
	FileSize        *float64 `json:"fileSize"`         // MB
	SpeedSize       *float64 `json:"speedSize"`        // Mbps
}

// TraceResultVO 路由追踪诊断结果
type TraceResultVO struct {
	TestTime int64     `json:"testTime"`
	Server   string    `json:"server"`
	Status   string    `json:"status"`
	RouteVOS []RouteVO `json:"routeVOS"`
}

// RouteVO 路由跳
type RouteVO struct {
	HopHost        string `json:"hopHost"`
	HopHostAddress string `json:"hopHostAddress"`
	HopErrorCode   *int   `json:"hopErrorCode"`
	HopRTTimes     string `json:"hopRTTimes"`
}

// ---------- Internal types ----------

// paramVo 构建 TR-069 SetParameterValues 参数
type paramVo struct {
	ParamName  string `json:"paramName"`
	ParamValue string `json:"paramValue"`
	ParamType  string `json:"paramType"`
}

// operationMessage Redis operation_queue 消息体
type operationMessage struct {
	EventType      string `json:"eventType"`
	NeNeid         int64  `json:"neNeid"`
	Operation      string `json:"operation"`
	OperationParam string `json:"operationParam"`
	OperationUser  string `json:"operationUser"`
	CommandTrackId int64  `json:"commandTrackId"`
	ExpiredAt      int64  `json:"expiredAt"`
}
