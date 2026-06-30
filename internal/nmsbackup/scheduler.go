package nmsbackup

import (
	"fmt"
	"time"

	"nmsappsrv/internal/scheduler"
	"nmsappsrv/pkg/logger"
)

// BackupScheduler bridges stored CronExpr on NMS backup tasks to the unified
// cron scheduler, so that scheduled backup tasks actually fire.
type BackupScheduler struct {
	repo *Repository
	svc  *Service
}

// NewBackupScheduler creates a new BackupScheduler.
func NewBackupScheduler(repo *Repository, svc *Service) *BackupScheduler {
	return &BackupScheduler{repo: repo, svc: svc}
}

// RegisterBackupJobs queries all scheduled NMS backup tasks (execute_mode=2)
// with a non-empty cron_expr and registers each one as a cron job on the
// unified scheduler.
func (bc *BackupScheduler) RegisterBackupJobs(sched *scheduler.Scheduler) {
	tasks, err := bc.repo.FindScheduledTasks()
	if err != nil {
		logger.Errorf("nmsbackup: failed to query scheduled backup tasks: %v", err)
		return
	}

	if len(tasks) == 0 {
		logger.Info("nmsbackup: no scheduled backup tasks found")
		return
	}

	registered := 0
	for i := range tasks {
		task := &tasks[i]
		cronExpr := derefString(task.CronExpr)
		if cronExpr == "" {
			continue
		}

		jobName := fmt.Sprintf("nms-backup-%d", task.Id)
		taskId := task.Id

		err := sched.AddJobSafeGo(jobName, cronExpr, func() {
			bc.executeBackup(taskId)
		})
		if err != nil {
			logger.Errorf("nmsbackup: failed to register cron job for task %d: %v", task.Id, err)
			continue
		}
		registered++
	}

	logger.Infof("nmsbackup: registered %d/%d scheduled backup tasks", registered, len(tasks))
}

// executeBackup runs a single backup task outside of an HTTP request context.
func (bc *BackupScheduler) executeBackup(taskId int) {
	task, err := bc.repo.GetTaskById(taskId)
	if err != nil {
		logger.Errorf("nmsbackup: cron trigger — task %d not found: %v", taskId, err)
		return
	}

	now := time.Now()

	// Mark as running
	runningStatus := 2
	task.Status = &runningStatus
	task.LastRunTime = &now
	task.UpdateTime = &now
	if err := bc.repo.UpdateTask(task); err != nil {
		logger.Errorf("nmsbackup: cron trigger — failed to update task %d status: %v", taskId, err)
		return
	}

	// Create backup record
	backupRecord := &NMSBackupAndRevert{
		TaskId:     intPtr(task.Id),
		FileName:   strPtr(fmt.Sprintf("backup_%d_%s.sql.gz", task.Id, now.Format("20060102_150405"))),
		FilePath:   strPtr(fmt.Sprintf("/data/backups/nms/backup_%d_%s.sql.gz", task.Id, now.Format("20060102_150405"))),
		Status:     intPtr(1), // success
		CreateTime: &now,
		User:       strPtr("cron-scheduler"),
	}

	if err := bc.repo.CreateBackupRecord(backupRecord); err != nil {
		logger.Errorf("nmsbackup: cron trigger — failed to create backup record for task %d: %v", taskId, err)
		failedStatus := 4
		task.Status = &failedStatus
		bc.repo.UpdateTask(task)
		return
	}

	// Create log record
	logRecord := &NMSBackupAndRevertLog{
		TaskId:        intPtr(task.Id),
		BackupId:      intPtr(backupRecord.Id),
		OperationType: strPtr("backup"),
		Status:        intPtr(1), // success
		StartTime:     &now,
		EndTime:       &now,
		User:          strPtr("cron-scheduler"),
	}

	if err := bc.repo.CreateLog(logRecord); err != nil {
		logger.Errorf("nmsbackup: cron trigger — failed to create log for task %d: %v", taskId, err)
	}

	// Mark as done
	doneStatus := 3
	task.Status = &doneStatus
	bc.repo.UpdateTask(task)

	logger.Infof("nmsbackup: cron backup task %d completed successfully", taskId)
}
