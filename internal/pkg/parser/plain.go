package parser

import (
	"context"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/lomik/zapwriter"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"
)

var parsedCounter = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name:      "messages_parsed",
		Namespace: "syphon",
		Help:      "How many messages have been Succesfully parsed, partitioned by if they where tagged and if a timestamp was provided.",
	},
	[]string{"tagged", "timestamp"},
)

// RegisterMetrics for collection
func RegisterMetrics() {
	prometheus.MustRegister(parsedCounter)
}

// Plain plain text metrics parser
func Plain(ctx context.Context, in chan string, sendChannels *[]chan Metric) {
	for {
		select {
		case <-ctx.Done():
			return
		case s := <-in:
			go plainParse(ctx, s, sendChannels)
		}
	}
}

func removeDoubleDot(path string) string {
	return strings.ReplaceAll(path, "..", ".")
}

func plainParse(
	ctx context.Context,
	message string,
	sendChannels *[]chan Metric,
) {
	logger := zapwriter.Logger("parser.Plain")
	lines := strings.Split(message, "\n")
	for _, line := range lines {
		if err := parseLine(line, sendChannels); err != nil {
			logger.Error("error parsing line", zap.Error(err))
		}
	}
}

func parseLine(line string, sendChannels *[]chan Metric) error {
	fields := strings.Fields(line)
	if len(fields) < 2 {
		m, _ := strconv.ParseInt(line, 16, 64)
		if m == 0 {
			return nil
		}
		return fmt.Errorf("message to short: %#v", line)
	}

	value, err := strconv.ParseFloat(fields[1], 64)
	if err != nil || math.IsNaN(value) {
		return fmt.Errorf("value not parsable: %#v", line)
	}

	path, labels, err := graphiteTags(removeDoubleDot(fields[0]))
	if err == nil {
		timestamp, parsed := parseTimestamp(fields)
		metric := Metric{
			Path:      path,
			Labels:    labels,
			Value:     value,
			Timestamp: timestamp,
		}
		parsedCounter.WithLabelValues(
			strconv.FormatBool(len(labels) != 0),
			strconv.FormatBool(parsed),
		).Inc()
		for _, sendChan := range *sendChannels {
			sendChan <- metric
		}
	}
	return err
}

// parseTimestamp from message.
// If a timestamp cant be parsed used when message recieved
func parseTimestamp(fields []string) (int64, bool) {
	if len(fields) == 2 {
		return time.Now().Unix(), false
	}
	tsf, err := strconv.ParseFloat(fields[2], 64)
	if err != nil || math.IsNaN(tsf) {
		return time.Now().Unix(), false
	}
	return int64(tsf), true
}
