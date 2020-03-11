package parser

import (
	"context"
)

// PromMetric is raw prometheus Metric
type PromMetric struct {
	Metric    []string //["name", "key1", "value1", ...]
	Value     float64
	Timestamp int64
}

// Parse returns parser.Metric
func (m *PromMetric) Parse() Metric {
	labels := make(map[string]string)
	for i := 1; i < len(m.Metric); i += 2 {
		labels[m.Metric[i]] = m.Metric[i+1]
	}

	return Metric{
		Path:      m.Metric[0],
		Labels:    labels,
		Value:     m.Value,
		Timestamp: m.Timestamp,
	}
}

// Prom parses Prometheus metrics
func Prom(
	ctx context.Context,
	in chan PromMetric,
	sendChannels *[]chan Metric,
) {
	for {
		select {
		case <-ctx.Done():
			return
		case s := <-in:
			go promParse(ctx, s, sendChannels)
		}
	}
}

func promParse(
	ctx context.Context,
	message PromMetric,
	sendChannels *[]chan Metric,
) {
	metric := message.Parse()
	for _, sendChan := range *sendChannels {
		sendChan <- metric
	}
}
