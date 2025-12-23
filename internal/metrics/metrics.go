package metrics

import "github.com/prometheus/client_golang/prometheus"

var (
	Success *prometheus.CounterVec
	Failure *prometheus.CounterVec
	Last    *prometheus.GaugeVec
)
