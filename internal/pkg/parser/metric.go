package parser

import (
	"bytes"
	"fmt"
	"strings"

	api "github.com/lomik/carbon-clickhouse/grpc"
)

// Metric Structure of parsed message.
type Metric struct {
	Path      string
	Value     float64
	Labels    map[string]string
	Timestamp int64
}

func (m *Metric) labels() string {
	labels := new(bytes.Buffer)

	for key, value := range m.Labels {
		fmt.Fprintf(labels, ";%s=%s", key, value)
	}

	return labels.String()
}

func (m *Metric) String() string {
	return fmt.Sprintf(
		"%s%s %s %d\n",
		m.Path,
		m.labels(),
		trimTrailingZero(m.Value),
		m.Timestamp,
	)
}

// Grpc converts metrics to grpc format.
func (m *Metric) Grpc() *api.Metric {
	return &api.Metric{
		Metric: fmt.Sprintf("%s%s", m.Path, m.labels()),
		Points: []*api.Point{{Timestamp: uint32(m.Timestamp), Value: m.Value}},
	}
}

func trimTrailingZero(num float64) string {
	return strings.TrimRight(
		strings.TrimRight(
			fmt.Sprintf("%f", num),
			"0",
		),
		".",
	)
}
