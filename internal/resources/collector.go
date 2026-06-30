package resources

import (
	"context"
	"fmt"
	"sync"
	"time"

	"nmsappsrv/pkg/logger"
	"nmsappsrv/pkg/redis"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/mem"
)

const (
	collectorInterval = 30 * time.Second
	redisKeyCPU       = "system:resource:cpu"
	redisKeyMem       = "system:resource:mem"
	redisTTL          = 2 * time.Minute
)

// Collector periodically samples CPU and memory usage and caches results to Redis.
type Collector struct {
	mu      sync.Mutex
	running bool
	stopCh  chan struct{}
}

// NewCollector creates a new resource Collector.
func NewCollector() *Collector {
	return &Collector{}
}

// Start begins the background collection loop.
func (c *Collector) Start() {
	c.mu.Lock()
	if c.running {
		c.mu.Unlock()
		return
	}
	c.running = true
	c.stopCh = make(chan struct{})
	c.mu.Unlock()

	logger.Info("resource collector starting")
	go c.loop()
}

// Stop stops the collector gracefully.
func (c *Collector) Stop() {
	c.mu.Lock()
	defer c.mu.Unlock()
	if !c.running {
		return
	}
	c.running = false
	close(c.stopCh)
	logger.Info("resource collector stopped")
}

// loop is the main collector loop.
func (c *Collector) loop() {
	ticker := time.NewTicker(collectorInterval)
	defer ticker.Stop()

	// Sample immediately on start so values are available right away.
	c.sample()

	for {
		select {
		case <-c.stopCh:
			return
		case <-ticker.C:
			c.sample()
		}
	}
}

// sample collects CPU and memory usage and caches to Redis.
func (c *Collector) sample() {
	ctx := context.Background()

	// CPU
	percent, err := cpu.Percent(1*time.Second, false)
	if err != nil {
		logger.Warnf("resource collector: CPU sample failed: %v", err)
	} else if len(percent) > 0 {
		if err := redis.Set(ctx, redisKeyCPU, fmt.Sprintf("%.1f", percent[0]), redisTTL); err != nil {
			logger.Warnf("resource collector: failed to cache CPU to Redis: %v", err)
		}
	}

	// Memory
	v, err := mem.VirtualMemory()
	if err != nil {
		logger.Warnf("resource collector: memory sample failed: %v", err)
	} else {
		if err := redis.Set(ctx, redisKeyMem, fmt.Sprintf("%.1f", v.UsedPercent), redisTTL); err != nil {
			logger.Warnf("resource collector: failed to cache memory to Redis: %v", err)
		}
	}
}
