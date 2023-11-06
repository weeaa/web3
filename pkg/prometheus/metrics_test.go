package prometheus

import (
	"testing"
)

func TestPrometheusMetrics(t *testing.T) {
	c := make(chan struct{})
	prom := PromMetrics{}
	prom.Initialize()
	<-c
}
