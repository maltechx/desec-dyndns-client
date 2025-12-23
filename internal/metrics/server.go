package metrics

import (
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func Init() {
	Success = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "dyndns_success_total",
			Help: "Successful DNS updates",
		},
		[]string{"type", "host"},
	)

	Failure = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "dyndns_failure_total",
			Help: "Failed DNS updates",
		},
		[]string{"type", "host"},
	)

	Last = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "dyndns_last_update_timestamp",
			Help: "Last successful update timestamp",
		},
		[]string{"type", "host"},
	)

	prometheus.MustRegister(Success, Failure, Last)
}

func Start(addr string) {
	go func() {
		log.Printf("INFO: metrics listening on %s", addr)
		http.Handle("/metrics", promhttp.Handler())
		log.Fatal(http.ListenAndServe(addr, nil))
	}()
}
