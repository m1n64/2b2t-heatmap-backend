package metrics

import (
	"github.com/alexcesaro/statsd"
	"github.com/dgraph-io/ristretto"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/process"
	"os"
	"runtime"
	"time"
)

var startTime = time.Now()

type SystemCollector struct {
	metrics *statsd.Client
	proc    *process.Process
	mem     runtime.MemStats
	cache   *ristretto.Cache
}

func NewSystemCollector(m *statsd.Client, cache *ristretto.Cache) (*SystemCollector, error) {
	p, err := process.NewProcess(int32(os.Getpid()))
	if err != nil {
		return nil, err
	}
	return &SystemCollector{
		metrics: m,
		proc:    p,
		mem:     runtime.MemStats{},
		cache:   cache,
	}, nil
}

func (c *SystemCollector) Run(interval time.Duration) {
	ticker := time.NewTicker(interval)

	go func() {
		for range ticker.C {
			c.collect()
		}
	}()
}

func (c *SystemCollector) collect() {
	// CPU %
	if cpuPercent, err := c.proc.CPUPercent(); err == nil {
		c.metrics.Gauge("cpu.percent", int(cpuPercent))
	}

	runtime.ReadMemStats(&c.mem)

	// Memory
	if memInfo, err := c.proc.MemoryInfo(); err == nil {
		c.metrics.Gauge("mem.rss_bytes", int(memInfo.RSS))
		c.metrics.Gauge("mem.heap_bytes", int(memInfo.VMS))
	}

	c.metrics.Gauge("go.heap_alloc_bytes", int(c.mem.HeapAlloc))
	c.metrics.Gauge("go.heap_sys_bytes", int(c.mem.HeapSys))

	lastPause := c.mem.PauseNs[(c.mem.NumGC+255)%256]
	c.metrics.Gauge("go.gc_pause_ns", int(lastPause))
	c.metrics.Gauge("go.mallocs_total", int(c.mem.Mallocs))
	c.metrics.Gauge("go.gc_count", int(c.mem.NumGC))

	// Goroutines / threads
	c.metrics.Gauge("go.goroutines", runtime.NumGoroutine())

	if threads, err := c.proc.NumThreads(); err == nil {
		c.metrics.Gauge("go.threads", int(threads))
	}

	// System memory (container-aware)
	if vm, err := mem.VirtualMemory(); err == nil {
		c.metrics.Gauge("fs.used_bytes", int(vm.Used))
		c.metrics.Gauge("fs.free_bytes", int(vm.Free))
	}

	// Uptime
	c.metrics.Gauge("go.uptime_sec", int(time.Since(startTime).Seconds()))

	// Memory cache stats
	if c.cache != nil {
		c.metrics.Gauge("cache.hits", float64(c.cache.Metrics.Hits()))
		c.metrics.Gauge("cache.misses", float64(c.cache.Metrics.Misses()))
		c.metrics.Gauge("cache.ratio", c.cache.Metrics.Ratio())
		c.metrics.Histogram("cache.life_expectancy_sec", c.cache.Metrics.LifeExpectancySeconds())
	}
}
