package systemsettings

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"nmsappsrv/internal/middleware"
	"nmsappsrv/pkg/logger"
	"nmsappsrv/pkg/utils"
)

// SystemSettingsHandler handles system settings endpoints.
type SystemSettingsHandler struct {
	svc *SystemSettingsService
}

// NewSystemSettingsHandler creates a new SystemSettingsHandler.
func NewSystemSettingsHandler(db *gorm.DB, aesKey string) *SystemSettingsHandler {
	return &SystemSettingsHandler{
		svc: NewSystemSettingsService(db, aesKey),
	}
}

// ListDeviceSettings returns the device configuration for the current tenancy.
func (h *SystemSettingsHandler) ListDeviceSettings(c *gin.Context) {
	tenancyId := middleware.GetLicenseId(c)

	cfg, err := h.svc.GetDeviceSettings(tenancyId)
	if err != nil {
		logger.Errorf("Failed to get device settings: %v", err)
		utils.Error(c, 500, "Failed to get device settings")
		return
	}

	utils.Success(c, cfg)
}

// UpdateDeviceSettings updates the device configuration for the current tenancy.
func (h *SystemSettingsHandler) UpdateDeviceSettings(c *gin.Context) {
	var req UpdateDeviceSettingsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, 400, "Invalid request: "+err.Error())
		return
	}

	tenancyId := middleware.GetLicenseId(c)

	if err := h.svc.UpdateDeviceSettings(&req, tenancyId); err != nil {
		logger.Errorf("Failed to update device settings: %v", err)
		utils.Error(c, 500, "Failed to update device settings")
		return
	}

	utils.OK(c, nil)
}

// ListACSSettings returns the ACS configuration.
func (h *SystemSettingsHandler) ListACSSettings(c *gin.Context) {
	cfg, err := h.svc.GetACSConfig()
	if err != nil {
		logger.Errorf("Failed to get ACS config: %v", err)
		utils.Error(c, 500, "Failed to get ACS config")
		return
	}

	utils.Success(c, cfg)
}

// UpdateACSSettings updates the ACS configuration.
func (h *SystemSettingsHandler) UpdateACSSettings(c *gin.Context) {
	var req UpdateACSConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, 400, "Invalid request: "+err.Error())
		return
	}

	if err := h.svc.UpdateACSConfig(&req); err != nil {
		logger.Errorf("Failed to update ACS config: %v", err)
		utils.Error(c, 500, "Failed to update ACS config")
		return
	}

	utils.OK(c, nil)
}

// ListLogSettings returns the log configuration.
func (h *SystemSettingsHandler) ListLogSettings(c *gin.Context) {
	cfg, err := h.svc.GetLogConfig()
	if err != nil {
		logger.Errorf("Failed to get log config: %v", err)
		utils.Error(c, 500, "Failed to get log config")
		return
	}

	utils.Success(c, cfg)
}

// UpdateLogSettings updates the log configuration.
func (h *SystemSettingsHandler) UpdateLogSettings(c *gin.Context) {
	var req UpdateLogConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, 400, "Invalid request: "+err.Error())
		return
	}

	if err := h.svc.UpdateLogConfig(&req); err != nil {
		logger.Errorf("Failed to update log config: %v", err)
		utils.Error(c, 500, "Failed to update log config")
		return
	}

	utils.OK(c, nil)
}
