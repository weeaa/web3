package prometheus

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

type PromMetrics struct {
	GoroutineCount prometheus.Counter

	FriendTechMetrics      FriendTechMetrics
	EthereumWatcherMetrics EthereumWatcher
}

type EthereumWatcher struct {
}

type FriendTechMetrics struct {
	BlockchainDataLatency prometheus.Histogram
}

func (p *PromMetrics) Initialize() {
	go http.ListenAndServe(":9091", promhttp.Handler())
}

func NewPromMetrics() *PromMetrics {
	return &PromMetrics{
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
