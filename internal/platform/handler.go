package platform

import (
	"nmsappsrv/pkg/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Handler provides HTTP handlers for platform settings endpoints
type Handler struct {
	svc *Service
}

// NewHandler creates a new Handler
func NewHandler(db *gorm.DB, aesKeyHex string) *Handler {
	return &Handler{svc: NewService(NewRepository(db), aesKeyHex)}
}

// GetDate handles POST /api/v1/getDate
func (h *Handler) GetDate(c *gin.Context) {
	utils.Success(c, h.svc.GetTime())
}

// GetSupportedZone handles GET /api/v1/getSupportedZone
func (h *Handler) GetSupportedZone(c *gin.Context) {
	utils.Success(c, h.svc.GetSupportedZone())
}

// GetLogo handles GET /api/v1/getLogo
func (h *Handler) GetLogo(c *gin.Context) {
	// Try to get logo from license table
	logo := ""
	// The logo can come from the license table; for now return empty
	// In Java, it reads from tenancy.license_logo_base64
	utils.Success(c, logo)
}

// ListLogConfig handles GET /api/v1/listLogConfig
func (h *Handler) ListLogConfig(c *gin.Context) {
	cfg, err := h.svc.GetLogConfig()
	if err != nil {
		utils.Error(c, 500, err.Error())
		return
	}
	utils.Success(c, cfg)
}

// UpdateLogConfig handles POST /api/v1/updateLogConfig
func (h *Handler) UpdateLogConfig(c *gin.Context) {
	var cfg LogConfig
	if err := c.ShouldBindJSON(&cfg); err != nil {
		utils.Error(c, 400, "invalid request body")
		return
	}
	if err := h.svc.UpdateLogConfig(&cfg); err != nil {
		utils.Error(c, 400, err.Error())
		return
	}
	utils.Success(c, nil)
}

// GetFTPTransferLogConfig handles GET /api/v1/getFTPTransferLogConfig
func (h *Handler) GetFTPTransferLogConfig(c *gin.Context) {
	cfg, err := h.svc.GetFTPTransferLogConfig()
	if err != nil {
		utils.Error(c, 500, err.Error())
		return
	}
	utils.Success(c, cfg)
}

// UpdateFTPTransferLogConfig handles POST /api/v1/updateFTPTransferLogConfig
func (h *Handler) UpdateFTPTransferLogConfig(c *gin.Context) {
	var cfg FTPTransferLogConfig
	if err := c.ShouldBindJSON(&cfg); err != nil {
		utils.Error(c, 400, "invalid request body")
		return
	}
	if err := h.svc.UpdateFTPTransferLogConfig(&cfg); err != nil {
		utils.Error(c, 400, err.Error())
		return
	}
	utils.Success(c, nil)
}

// GetHECConfig handles GET /api/v1/getHECConfig
func (h *Handler) GetHECConfig(c *gin.Context) {
	cfg, err := h.svc.GetHECConfig()
	if err != nil {
		utils.Error(c, 500, err.Error())
		return
	}
	utils.Success(c, cfg)
}

// UpdateHECConfig handles POST /api/v1/updateHECConfig
func (h *Handler) UpdateHECConfig(c *gin.Context) {
	var cfg HECConfig
	if err := c.ShouldBindJSON(&cfg); err != nil {
		utils.Error(c, 400, "invalid request body")
		return
	}
	if err := h.svc.UpdateHECConfig(&cfg); err != nil {
		utils.Error(c, 400, err.Error())
		return
	}
	utils.Success(c, nil)
}

// ListNMSSecret handles GET /api/v1/listNMSSecret
func (h *Handler) ListNMSSecret(c *gin.Context) {
	secret, err := h.svc.ListNMSSecret()
	if err != nil {
		utils.Error(c, 500, err.Error())
		return
	}
	utils.Success(c, secret)
}

// UpdateNMSSecret handles POST /api/v1/updateNMSSecret
func (h *Handler) UpdateNMSSecret(c *gin.Context) {
	var secret NMSSecret
	if err := c.ShouldBindJSON(&secret); err != nil {
		utils.Error(c, 400, "invalid request body")
		return
	}
	if err := h.svc.UpdateNMSSecret(&secret); err != nil {
		utils.Error(c, 400, err.Error())
		return
	}
	utils.Success(c, nil)
}

// DownloadPasswordRSAPublicKey handles GET /api/v1/downloadPasswordRSAPublicKey
func (h *Handler) DownloadPasswordRSAPublicKey(c *gin.Context) {
	// File download from /home/cert/password/publicKey.pem
	// TODO: implement file streaming when deployment paths are configured
	c.File("/home/cert/password/publicKey.pem")
}

// DownloadPlatformLogs handles POST /api/v1/downloadPlatformLogs
func (h *Handler) DownloadPlatformLogs(c *gin.Context) {
	// TODO: implement async log collection and ZIP download
	// This requires background goroutine + WebSocket notification pattern
	utils.Success(c, nil)
}

// DownloadNMSManualDocument handles GET /api/v1/downloadNMSManualDocument
func (h *Handler) DownloadNMSManualDocument(c *gin.Context) {
	// File download from /home/nms_manual.pdf (or configured path)
	// TODO: implement when manual file path is configured
	c.File("/home/nms_manual.pdf")
}
