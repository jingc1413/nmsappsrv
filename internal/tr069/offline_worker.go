package tr069

import (
	"context"
	"fmt"
	"sync"
	"time"

	"nmsappsrv/internal/alarm"
	"nmsappsrv/internal/device"
	"nmsappsrv/pkg/logger"
	"nmsappsrv/pkg/redis"
	"nmsappsrv/pkg/utils"

	"gorm.io/gorm"
)

// OfflineWorker periodically checks device online status via Redis keys
// and creates/clears offline alarms accordingly.
type OfflineWorker struct {
	mu      sync.Mutex
	running bool
	db      *gorm.DB
}

// NewOfflineWorker creates a new OfflineWorker.
func NewOfflineWorker(db *gorm.DB) *OfflineWorker {
	return &OfflineWorker{db: db}
}

// Start begins the periodic offline detection loop.
func (w *OfflineWorker) Start() {
	w.mu.Lock()
	if w.running {
		w.mu.Unlock()
		return
	}
	w.running = true
	w.mu.Unlock()

	logger.Info("offline detection worker starting")

	utils.SafeGo("offline-worker", func() {
		w.pollLoop()
	})
}

// Stop stops the worker.
func (w *OfflineWorker) Stop() {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.running = false
	logger.Info("offline detection worker stopped")
}

// IsRunning returns whether the worker is currently running.
func (w *OfflineWorker) IsRunning() bool {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.running
}

// pollLoop runs the detection loop every 5 seconds.
func (w *OfflineWorker) pollLoop() {
	for w.IsRunning() {
		w.checkAllDevices()
		time.Sleep(5 * time.Second)
	}
}

// checkAllDevices queries all non-deleted devices and checks their online status.
func (w *OfflineWorker) checkAllDevices() {
	var devices []device.CpeElement
	if err := w.db.Where("deleted = ?", false).Find(&devices).Error; err != nil {
		logger.Errorf("offline worker: failed to query devices: %v", err)
		return
	}

	ctx := context.Background()
	now := time.Now()

	for i := range devices {
		dev := &devices[i]
		if dev.SerialNumber == nil || *dev.SerialNumber == "" {
			continue
		}
		sn := *dev.SerialNumber

		// Check the TR069 online key (written by event_processor on each Inform)
		onlineKey := fmt.Sprintf("device:online:%s", sn)
		isOnline := redis.Exists(ctx, onlineKey)

		if isOnline {
			// Bridge: set the dashboard/PM-facing key pattern online_{neId}
			dashboardKey := fmt.Sprintf("online_%d", dev.NeNeid)
			redis.Set(ctx, dashboardKey, "1", 5*time.Minute)

			// Clear any active offline alarm for this device
			w.clearOfflineAlarm(dev.NeNeid, &now)
		} else {
			// Device is offline: create or ensure an offline alarm exists
			w.createOfflineAlarm(dev, &now)

			// Also clear the dashboard key
			dashboardKey := fmt.Sprintf("online_%d", dev.NeNeid)
			redis.Del(ctx, dashboardKey)
		}
	}
}

// createOfflineAlarm creates an omc_device_offline alarm if one doesn't already exist.
func (w *OfflineWorker) createOfflineAlarm(dev *device.CpeElement, now *time.Time) {
	// Check if there's already an active (uncleared) offline alarm for this device
	var existing alarm.Alarm
	err := w.db.Where("element_id = ? AND alarm_identifier = ? AND alarm_status != 0",
		dev.NeNeid, "omc_device_offline").First(&existing).Error
	if err == nil {
		// Already has an active offline alarm, no need to create another
		return
	}

	severity := "Major"
	alarmIdentifier := "omc_device_offline"
	probableCause := "EquipmentOffLine"
	alarmSource := "TR069"
	eventType := "CommunicationsAlarm"
	alarmStatus := 1 // 1=active
	alarmType := 1   // communications alarm

	deviceName := ""
	if dev.DeviceName != nil {
		deviceName = *dev.DeviceName
	}
	networkElement := deviceName
	if networkElement == "" && dev.SerialNumber != nil {
		networkElement = *dev.SerialNumber
	}

	newAlarm := alarm.Alarm{
		Severity:        &severity,
		AlarmIdentifier: &alarmIdentifier,
		ProbableCause:   &probableCause,
		AlarmSource:     &alarmSource,
		NetworkElement:  &networkElement,
		EventType:       &eventType,
		AlarmStatus:     &alarmStatus,
		AlarmType:       &alarmType,
		EventTime:       now,
		UpdateTime:      now,
		CreateTime:      now,
		ElementId:       &dev.NeNeid,
		LicenseId:       dev.LicenseId,
		SpecificProblem: stringPtr("Device offline detected by periodic check"),
	}

	if dev.SerialNumber != nil {
		alarmIdStr := fmt.Sprintf("offline_%d_%d", dev.NeNeid, now.Unix())
		newAlarm.AlarmId = &alarmIdStr
	}

	if err := w.db.Create(&newAlarm).Error; err != nil {
		logger.Errorf("offline worker: failed to create offline alarm for device %d: %v", dev.NeNeid, err)
	} else {
		sn := ""
		if dev.SerialNumber != nil {
			sn = *dev.SerialNumber
		}
		logger.Infof("offline worker: created offline alarm for device %d (SN=%s)", dev.NeNeid, sn)
	}
}

// clearOfflineAlarm clears any active offline alarm for the given device.
func (w *OfflineWorker) clearOfflineAlarm(neId int64, now *time.Time) {
	result := w.db.Model(&alarm.Alarm{}).
		Where("element_id = ? AND alarm_identifier = ? AND alarm_status != 0",
			neId, "omc_device_offline").
		Updates(map[string]interface{}{
			"alarm_status":  0,
			"cleared_time":  now,
			"update_time":   now,
		})
	if result.RowsAffected > 0 {
		logger.Infof("offline worker: cleared offline alarm for device %d", neId)
	}
}
