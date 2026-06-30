package scheduler

import (
	"fmt"
	"sync"

	"nmsappsrv/pkg/logger"
	"nmsappsrv/pkg/utils"

	"github.com/robfig/cron/v3"
	"gorm.io/gorm"
)

// Scheduler is a unified cron scheduler that manages all periodic jobs
// across the application using a single cron engine with seconds precision.
type Scheduler struct {
	cron    *cron.Cron
	db      *gorm.DB
	mu      sync.Mutex
	running bool
	entries map[string]cron.EntryID
}

// NewScheduler creates a new unified Scheduler with seconds-precision cron.
func NewScheduler(db *gorm.DB) *Scheduler {
	return &Scheduler{
		cron:    cron.New(cron.WithSeconds()),
		db:      db,
		entries: make(map[string]cron.EntryID),
	}
}

// Start begins the unified cron scheduler.
func (s *Scheduler) Start() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.running {
		return
	}
	s.running = true
	s.cron.Start()
	logger.Info("unified cron scheduler started")
}

// Stop stops the unified cron scheduler gracefully, waiting for running jobs.
func (s *Scheduler) Stop() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.running {
		return
	}
	s.running = false
	ctx := s.cron.Stop()
	<-ctx.Done()
	logger.Info("unified cron scheduler stopped")
}

// IsRunning returns whether the scheduler is currently running.
func (s *Scheduler) IsRunning() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.running
}

// AddJob registers a named cron job with the given spec and command function.
// The spec uses 6-field cron format (seconds minute hour dom month dow) since
// the scheduler is configured with cron.WithSeconds().
//
// The command is wrapped with panic recovery and logging so that a single job
// failure does not affect other scheduled jobs.
func (s *Scheduler) AddJob(name, spec string, cmd func()) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.entries[name]; exists {
		return fmt.Errorf("cron job %q already registered", name)
	}

	wrappedCmd := func() {
		defer func() {
			if r := recover(); r != nil {
				logger.Errorf("cron job %q panicked: %v", name, r)
			}
		}()
		logger.Infof("cron job %q triggered", name)
		cmd()
	}

	entryID, err := s.cron.AddFunc(spec, wrappedCmd)
	if err != nil {
		return fmt.Errorf("failed to add cron job %q with spec %q: %w", name, spec, err)
	}

	s.entries[name] = entryID
	logger.Infof("cron job %q registered with spec %q (entryID=%d)", name, spec, entryID)
	return nil
}

// RemoveJob removes a previously registered cron job by name.
func (s *Scheduler) RemoveJob(name string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	entryID, exists := s.entries[name]
	if !exists {
		return
	}
	s.cron.Remove(entryID)
	delete(s.entries, name)
	logger.Infof("cron job %q removed", name)
}

// AddJobSafeGo registers a named cron job that runs the command in a
// goroutine managed by utils.SafeGo, providing panic recovery at both
// the cron wrapper and the goroutine level.
func (s *Scheduler) AddJobSafeGo(name, spec string, cmd func()) error {
	return s.AddJob(name, spec, func() {
		utils.SafeGo("cron-"+name, cmd)
	})
}

// DB returns the database handle for use by job implementations.
func (s *Scheduler) DB() *gorm.DB {
	return s.db
}
