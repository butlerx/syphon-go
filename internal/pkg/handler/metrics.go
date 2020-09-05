package handler

import "github.com/prometheus/client_golang/prometheus"

var receiversGauge = prometheus.NewGaugeVec(
	prometheus.GaugeOpts{
		Name:      "receiver_worker",
		Namespace: "syphon",
		Help:      "How many receiver workers are running, partitioned by Protocol.",
	},
	[]string{"protocol"},
)

var messageRecievedCounter = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name:      "messages_received",
		Namespace: "syphon",
		Help:      "How many messages have been received, partitioned by Protocol.",
	},
	[]string{"protocol"},
)

// RegisterMetrics for collection.
func RegisterMetrics() {
	prometheus.MustRegister(receiversGauge)
	prometheus.MustRegister(messageRecievedCounter)
}
