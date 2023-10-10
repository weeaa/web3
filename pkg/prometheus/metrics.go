package prometheus

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type Prometrics struct {
	msgCounter prometheus.Counter
	msgLatency prometheus.Histogram
}

func (p *Prometrics) WithMetrics() {

}

func NewPromMetrics() *Prometrics {
	return &Prometrics{
		msgCounter: promauto.NewCounter(prometheus.CounterOpts{
			Name: "",
			Help: "",
		}),
	}
}
