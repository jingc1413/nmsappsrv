package devicelog

import (
	"nmsappsrv/pkg/utils"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Handler struct {
	svc *Service
}

func NewHandler(db *gorm.DB) *Handler {
	return &Handler{
		svc: NewService(db),
	}
}

func (h *Handler) AddLogCollectionTask(c *gin.Context) {
	var req AddLogCollectionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, 400, "Invalid request: "+err.Error())
		return
	}

	if err := h.svc.AddLogCollectionTask(c, &req); err != nil {
		utils.Error(c, 500, "Failed to add log collection task: "+err.Error())
		return
	}

	utils.Success(c, nil)
}

func (h *Handler) ListLogCollectionResults(c *gin.Context) {
	var req ListLogCollectionResultRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, 400, "Invalid request: "+err.Error())
		return
	}

	results, total, err := h.svc.ListLogCollectionResults(c, &req)
	if err != nil {
		utils.Error(c, 500, "Failed to list log collection results: "+err.Error())
		return
	}

	page := req.Page
	pageSize := req.PageSize
	utils.Paginated(c, results, total, page, pageSize)
}

func (h *Handler) DeleteAllLogFile(c *gin.Context) {
	var req DeleteAllLogFileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, 400, "Invalid request: "+err.Error())
		return
	}

	if err := h.svc.DeleteAllLogFile(&req); err != nil {
		utils.Error(c, 500, "Failed to delete all log files: "+err.Error())
		return
	}

	utils.Success(c, nil)
}

func (h *Handler) DeleteLogFile(c *gin.Context) {
	var req DeleteLogFileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, 400, "Invalid request: "+err.Error())
		return
	}

	if err := h.svc.DeleteLogFile(&req); err != nil {
		utils.Error(c, 500, "Failed to delete log file: "+err.Error())
		return
	}

	utils.Success(c, nil)
}

func (h *Handler) DownloadLogFile(c *gin.Context) {
	logIdStr := c.Query("id")
	if logIdStr == "" {
		utils.Error(c, 400, "Missing log id parameter")
		return
	}

	logId, err := strconv.ParseInt(logIdStr, 10, 64)
	if err != nil {
		utils.Error(c, 400, "Invalid log id parameter")
		return
	}

	filePath, err := h.svc.GetLogFile(logId)
	if err != nil {
		utils.Error(c, 500, "Failed to get log file: "+err.Error())
		return
	}

	// Serve the file for download
	c.Header("Content-Disposition", "attachment; filename="+filePath)
	c.Header("Content-Type", "application/octet-stream")
	c.File(filePath)
}

func (h *Handler) ListLogFiles(c *gin.Context) {
	var req ListLogFileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, 400, "Invalid request: "+err.Error())
		return
	}

	files, total, err := h.svc.ListLogFiles(&req)
	if err != nil {
		utils.Error(c, 500, "Failed to list log files: "+err.Error())
		return
	}

	page := req.Page
	pageSize := req.PageSize
	utils.Paginated(c, files, total, page, pageSize)
}

func (h *Handler) EnablePeriodicUpload(c *gin.Context) {
	var req EnablePeriodicUploadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, 400, "Invalid request: "+err.Error())
		return
	}

	if err := h.svc.EnablePeriodicUpload(c, &req); err != nil {
		utils.Error(c, 500, "Failed to enable periodic upload: "+err.Error())
		return
	}

	utils.Success(c, nil)
}

func (h *Handler) DisablePeriodicUpload(c *gin.Context) {
	var req DisablePeriodicUploadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, 400, "Invalid request: "+err.Error())
		return
	}

	if err := h.svc.DisablePeriodicUpload(c, &req); err != nil {
		utils.Error(c, 500, "Failed to disable periodic upload: "+err.Error())
		return
	}

	utils.Success(c, nil)
}
