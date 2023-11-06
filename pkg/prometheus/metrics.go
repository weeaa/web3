package prometheus

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

const (
	DefaultPort = ":2115"
)

type PromMetrics struct {
	DatabaseQueryDuration *prometheus.SummaryVec

	GoroutineCount prometheus.Counter

	FriendTechMetrics      FriendTechMetrics
	EthereumWatcherMetrics EthereumWatcher
}

type DatabaseQueriesDuration struct {
}

type EthereumWatcher struct {
}

type FriendTechMetrics struct {
	BlockchainDataLatency prometheus.Histogram
}

func (p *PromMetrics) Initialize(port string) error {
	if err := http.ListenAndServe(port, promhttp.Handler()); err != nil {
		return fmt.Errorf("error initializing PromMetrics: %w", err)
	}
	return nil
}

func NewPromMetrics() *PromMetrics {
	return &PromMetrics{
		DatabaseQueryDuration: promauto.NewSummaryVec(prometheus.SummaryOpts{
			Name: "db_query_duration_seconds",
			Help: "Database query execution time in seconds",
		},
			[]string{"query_name"},
		),
		GoroutineCount: promauto.NewCounter(prometheus.CounterOpts{
			Name: "",
			Help: "",
		}),
		EthereumWatcherMetrics: EthereumWatcher{},
		FriendTechMetrics: FriendTechMetrics{
			BlockchainDataLatency: promauto.NewHistogram(prometheus.HistogramOpts{
				Name: "blockchain_data_latency",
				Help: "friend tech blockchain data treatment's latency",
			}),
		},
	}
}
