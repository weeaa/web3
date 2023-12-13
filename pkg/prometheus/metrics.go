package prometheus

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"
	"github.com/shirou/gopsutil/disk"
	"net/http"
	"runtime"
	"time"
)

const (
	DefaultPort = ":2115"
)

type Metrics struct {
	logger zerolog.Logger
}

type PromMetrics struct {
	DatabaseQueryDuration *prometheus.SummaryVec
	
	DatabaseMetrics       DatabaseMetrics
	InfrastructureMetrics InfrastructureMetrics
	
	FriendTechMetrics      FriendTechMetrics
	EthereumWatcherMetrics EthereumWatcher
	
	PanicsRecovered prometheus.Counter
}

type InfrastructureMetrics struct {
	GarbageCollection prometheus.Histogram
	MemoryUsage       prometheus.Gauge
	CpuUsage          prometheus.Gauge
	GoroutineCount    prometheus.Gauge
	DiskOperations    DiskOperations
}

type DiskOperations struct {
	Read  prometheus.Gauge
	Write prometheus.Gauge
}

type DatabaseMetrics struct {
	QueryExecutionTime *prometheus.SummaryVec
	ErrorRate          prometheus.Counter
}

type EthereumWatcher struct {
}

type FriendTechMetrics struct {
	BlockchainDataLatency prometheus.Histogram
}

type UnisatMetrics struct {
	Latency prometheus.Histogram
}

func (p *PromMetrics) Initialize(port string) error {
	if err := http.ListenAndServe(port, promhttp.Handler()); err != nil {
		return fmt.Errorf("error initializing PromMetrics: %w", err)
	}
	return nil
}

func NewPromMetrics() *PromMetrics {
	
	return &PromMetrics{
		DatabaseMetrics: DatabaseMetrics{
			QueryExecutionTime: promauto.NewSummaryVec(prometheus.SummaryOpts{
				Name: "db_query_duration_ms",
				Help: "Database query execution time in milliseconds",
			},
				[]string{
					"",
				},
			),
		},
		InfrastructureMetrics: InfrastructureMetrics{
			GarbageCollection: promauto.NewHistogram(prometheus.HistogramOpts{
				Name: "garbage_collection_frequency_duration",
				Help: "Track garbage collection frequency and duration",
			}),
			MemoryUsage: promauto.NewGauge(prometheus.GaugeOpts{
				Name: "memory_usage",
				Help: "Monitor heap, stack, and system memory usage",
			}),
			CpuUsage: promauto.NewGauge(prometheus.GaugeOpts{
				Name: "cpu_usage",
				Help: "Monitor the CPU usage of the servers running the application",
			}),
			GoroutineCount: promauto.NewGauge(prometheus.GaugeOpts{
				Name: "goroutine_count",
				Help: "Number of running goroutines currently running",
			}),
			/*
			NetworkTraffic: promauto.NewCounter(prometheus.CounterOpts{
				Name: "network_traffic",
				Help: "Incoming and outgoing network traffic statistics",
			}),
			*/
			DiskOperations: promauto.NewCounter(prometheus.CounterOpts{
				Name: "disk_operations",
				Help: "Disk read/write operations statistics",
			}),
		},
		
		EthereumWatcherMetrics: EthereumWatcher{},
		FriendTechMetrics: FriendTechMetrics{
			BlockchainDataLatency: promauto.NewHistogram(prometheus.HistogramOpts{
				Name: "blockchain_data_latency",
				Help: "friend tech blockchain data treatment's latency",
			}),
		},
	}
}

// MonitorGoRoutinesCount monitors the number of goroutines running in the InfrastructureMetrics instance.
// It updates the GoroutineCount gauge with the current count of goroutines every second.
func (i *InfrastructureMetrics) monitorGoRoutinesCount() {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				i.monitorGoRoutinesCount()
			}
		}()
		for {
			i.GoroutineCount.Set(float64(runtime.NumGoroutine()))
			time.Sleep(time.Second)
		}
	}()
}

// MonitorGarbageCollection monitors the garbage collection statistics of the InfrastructureMetrics instance.
// It periodically prints the current memory allocation, total allocation, system memory usage, and number of garbage collections.

// I have no clue if it is even useful but here we go ;)
func (i *InfrastructureMetrics) monitorInfrastructure() {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				i.monitorInfrastructure()
			}
		}()
		for {
			var m runtime.MemStats
			runtime.ReadMemStats(&m)
			
			i.CpuUsage.Set(m.GCCPUFraction)
			i.GarbageCollection.Observe(float64(m.NumGC))
			i.MemoryUsage.Set(float64(m.Sys))
			
			partitions, _ := disk.Partitions(false)
			for _, partition := range partitions {
				ioStat, _ := disk.IOCounters(partition.Mountpoint)
				for _, io := range ioStat {
					i.DiskOperations.Read.Set(float64(io.ReadCount))
					i.DiskOperations.Write.Set(float64(io.WriteCount))
					//log.Printf("ReadCount: %v, WriteCount: %v\n", io.ReadCount, io.WriteCount)
				}
			}
			time.Sleep(time.Minute)
		}
	}()
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}
