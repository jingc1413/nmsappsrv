package mml

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"nmsappsrv/internal/mq"
	"nmsappsrv/internal/tr069"
	"nmsappsrv/internal/tr069/soap"
	"nmsappsrv/pkg/logger"
	"nmsappsrv/pkg/redis"
	"nmsappsrv/pkg/utils"

	"gorm.io/gorm"
)

// MMLMessage is the JSON payload pushed to the queue:mml Redis queue.
type MMLMessage struct {
	ElementId int64                  `json:"element_id"`
	Command   string                 `json:"command"`
	Params    map[string]interface{} `json:"params"`
	ResultId  int                    `json:"result_id"`
}

// MMLWorker consumes messages from the MML queue and dispatches
// MML commands to devices via TR-069 SetParameterValues.
type MMLWorker struct {
	db       *gorm.DB
	opSender *tr069.OperationSender
	mu       sync.Mutex
	running  bool
	stopCh   chan struct{}
}

// NewMMLWorker creates a new MMLWorker.
func NewMMLWorker(db *gorm.DB, msgManager *tr069.MessageManager) *MMLWorker {
	return &MMLWorker{
		db:       db,
		opSender: tr069.NewOperationSender(db, msgManager),
	}
}

// Start begins the MML command processing loop.
func (w *MMLWorker) Start() {
	w.mu.Lock()
	if w.running {
		w.mu.Unlock()
		return
	}
	w.running = true
	w.stopCh = make(chan struct{})
	w.mu.Unlock()

	logger.Info("MML worker starting")

	utils.SafeGo("mml-worker", func() {
		w.pollLoop()
	})
}

// Stop stops the worker gracefully.
func (w *MMLWorker) Stop() {
	w.mu.Lock()
	defer w.mu.Unlock()
	if !w.running {
		return
	}
	w.running = false
	close(w.stopCh)
	logger.Info("MML worker stopped")
}

// IsRunning returns whether the worker is currently running.
func (w *MMLWorker) IsRunning() bool {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.running
}

// pollLoop continuously polls the MML Redis queue for command messages.
func (w *MMLWorker) pollLoop() {
	for w.IsRunning() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

		result, err := redis.BRPop(ctx, 5*time.Second, mq.MMLQueue)
		cancel()

		if err != nil {
			if err.Error() == "redis: nil" {
				continue
			}
			if !w.IsRunning() {
				return
			}
			logger.Debugf("MML worker queue poll error (may be timeout): %v", err)
			time.Sleep(1 * time.Second)
			continue
		}

		if len(result) < 2 {
			continue
		}

		// result[0] is the queue name, result[1] is the message
		msgJSON := result[1]

		var msg MMLMessage
		if err := json.Unmarshal([]byte(msgJSON), &msg); err != nil {
			logger.Errorf("MML worker failed to unmarshal message: %v, data: %s", err, msgJSON)
			continue
		}

		logger.Infof("MML worker: processing command=%s for element=%d resultId=%d", msg.Command, msg.ElementId, msg.ResultId)
		w.processMmlCommand(&msg)
	}
}

// processMmlCommand handles a single MML command message.
// It looks up the device, builds TR-069 parameters from the MML command,
// and sends them via SetParameterValues.
func (w *MMLWorker) processMmlCommand(msg *MMLMessage) {
	// 1. Look up the device serial number from element_id
	var sn string
	type deviceLookup struct {
		SerialNumber *string `gorm:"column:serial_number"`
	}
	var dev deviceLookup
	if err := w.db.Table("cpe_element").
		Select("serial_number").
		Where("ne_neid = ? AND deleted = ?", msg.ElementId, false).
		First(&dev).Error; err != nil {
		faultMsg := fmt.Sprintf("device %d not found: %v", msg.ElementId, err)
		logger.Errorf("MML worker: %s", faultMsg)
		w.updateResultStatus(msg.ResultId, 3, faultMsg)
		return
	}

	if dev.SerialNumber == nil || *dev.SerialNumber == "" {
		faultMsg := fmt.Sprintf("device %d has no serial number", msg.ElementId)
		logger.Errorf("MML worker: %s", faultMsg)
		w.updateResultStatus(msg.ResultId, 3, faultMsg)
		return
	}
	sn = *dev.SerialNumber

	// 2. Build TR-069 parameter values from the MML command.
	// The MML command string is sent as a vendor-specific parameter value.
	// The parameter path encodes the command name so the device can interpret it.
	paramName := fmt.Sprintf("Device.VendorMML.Command.%s", msg.Command)
	paramValue := msg.Command

	// If there are additional parameters, encode them as JSON in a sibling parameter
	if len(msg.Params) > 0 {
		paramsJSON, err := json.Marshal(msg.Params)
		if err == nil {
			paramValue = string(paramsJSON)
		}
	}

	spvParams := []soap.ParameterValueStruct{
		{
			Name:  paramName,
			Value: paramValue,
		},
	}

	// 3. Send SetParameterValues to the device
	operationId := fmt.Sprintf("mml_%d_%d", msg.ResultId, time.Now().Unix())
	paramKey := fmt.Sprintf("mml_%d", msg.ResultId)

	if err := w.opSender.SendSetParameterValues(sn, spvParams, paramKey, operationId); err != nil {
		faultMsg := fmt.Sprintf("SPV send failed for device %d (SN=%s): %v", msg.ElementId, sn, err)
		logger.Errorf("MML worker: %s", faultMsg)
		w.updateResultStatus(msg.ResultId, 3, faultMsg)
		return
	}

	// 4. Update result status to 2 (sent)
	w.updateResultStatus(msg.ResultId, 2, "")
	logger.Infof("MML worker: command=%s sent to device %d (SN=%s), operationId=%s", msg.Command, msg.ElementId, sn, operationId)
}

// updateResultStatus updates the MmlExecuteResult status in the database.
// status: 0=pending, 2=sent, 3=failed
func (w *MMLWorker) updateResultStatus(resultId int, status int, faultString string) {
	updates := map[string]interface{}{
		"status": status,
	}
	if faultString != "" {
		updates["fault_string"] = faultString
		updates["has_fault"] = true
	}
	if err := w.db.Model(&MmlExecuteResult{}).Where("id = ?", resultId).Updates(updates).Error; err != nil {
		logger.Errorf("MML worker: failed to update result %d status: %v", resultId, err)
	}
}
