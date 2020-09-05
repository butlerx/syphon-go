package connection

import "github.com/prometheus/client_golang/prometheus"

// Max size of send buffer.
const bufferSendSize = 900

var messageSentCounter = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name:      "messages_sent",
		Namespace: "syphon",
		Help:      "How many messages have been sent, partitioned by Protocol, endpoint.",
	},
	[]string{"protocol", "endpoint"},
)

// RegisterMetrics for collection.
func RegisterMetrics() {
	prometheus.MustRegister(messageSentCounter)
}
