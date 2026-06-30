package systemsettings

// DeviceConfig represents the device configuration stored per-tenancy in system_config.
type DeviceConfig struct {
	AutoRegistration    *bool   `json:"autoRegistration"`
	AutoRegistrationKey *string `json:"autoRegistrationKey"`
	MaxDeviceCount      *int    `json:"maxDeviceCount"`
	DefaultDeviceName   *string `json:"defaultDeviceName"`
}

// ACSConfig represents the NMS ACS configuration stored in system_config as key "nms_config".
type ACSConfig struct {
	AcsUrl            *string `json:"acsUrl"`
	AcsUsername       *string `json:"acsUsername"`
	AcsPassword       *string `json:"acsPassword"` // encrypted
	ConnectionTimeout *int    `json:"connectionTimeout"`
	InformInterval    *int    `json:"informInterval"`
	UDPPort           *int    `json:"udpPort"`
	TR069Enabled      *bool   `json:"tr069Enabled"`
}

// LogConfig represents the NMS log configuration stored in system_config as key "nms_log_config".
type LogConfig struct {
	RetentionDays *int  `json:"retentionDays"`
	MaxFileSizeMb *int  `json:"maxFileSizeMb"`
	AutoCleanup   *bool `json:"autoCleanup"`
}

// SysParameter represents a row from sys_parameter table (key-value system params).
type SysParameter struct {
	Id    int     `gorm:"primaryKey;autoIncrement" json:"id"`
	Key   *string `gorm:"column:config_key;type:varchar(255);uniqueIndex" json:"key"`
	Value *string `gorm:"column:config_value;type:longtext" json:"value"`
}

func (SysParameter) TableName() string { return "sys_parameter" }

// UpdateDeviceSettingsRequest represents the request to update device settings.
type UpdateDeviceSettingsRequest struct {
	AutoRegistration    *bool   `json:"autoRegistration"`
	AutoRegistrationKey *string `json:"autoRegistrationKey"`
	MaxDeviceCount      *int    `json:"maxDeviceCount"`
	DefaultDeviceName   *string `json:"defaultDeviceName"`
}

// UpdateACSConfigRequest represents the request to update ACS configuration.
type UpdateACSConfigRequest struct {
	AcsUrl            *string `json:"acsUrl"`
	AcsUsername       *string `json:"acsUsername"`
	AcsPassword       *string `json:"acsPassword"`
	ConnectionTimeout *int    `json:"connectionTimeout"`
	InformInterval    *int    `json:"informInterval"`
	UDPPort           *int    `json:"udpPort"`
	TR069Enabled      *bool   `json:"tr069Enabled"`
}

// UpdateLogConfigRequest represents the request to update log configuration.
type UpdateLogConfigRequest struct {
	RetentionDays *int  `json:"retentionDays"`
	MaxFileSizeMb *int  `json:"maxFileSizeMb"`
	AutoCleanup   *bool `json:"autoCleanup"`
}
