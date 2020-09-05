package uploader

import "github.com/prometheus/client_golang/prometheus"

const bufferSendTimer = 15

var uploaderGuage = prometheus.NewGaugeVec(
	prometheus.GaugeOpts{
		Namespace: "syphon",
		Name:      "uploaders",
		Help:      "Number of Uploaders running, partitioned by type.",
	},
	[]string{"type", "pattern"},
)

// RegisterMetrics with collector.
func RegisterMetrics() {
	prometheus.MustRegister(uploaderGuage)
}
