package ntp

// SystemConfig mirrors system_config table for key-value config storage.
type SystemConfig struct {
	Id    int     `gorm:"primaryKey;autoIncrement" json:"id"`
	Key   *string `gorm:"column:config_key;type:varchar(255);uniqueIndex" json:"key"`
	Value *string `gorm:"column:config_value;type:longtext" json:"value"`
}

func (SystemConfig) TableName() string { return "system_config" }

// NTPConfig is the JSON payload stored in system_config (key="ntp").
type NTPConfig struct {
	NTPServer string `json:"ntpServer"`
	Enable    bool   `json:"enable"`
}

// NTPConfigRequest is the JSON body for POST /updateNTPConfig.
type NTPConfigRequest struct {
	NTPServer string `json:"ntpServer"`
	Enable    bool   `json:"enable"`
}

// NTPStatusResponse is the response for POST /getNTPStatus.
type NTPStatusResponse struct {
	Status string `json:"status"`
}
