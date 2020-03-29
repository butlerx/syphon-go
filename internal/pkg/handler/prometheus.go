package handler

import (
	"context"
	"io/ioutil"
	"math"
	"net/http"

	"github.com/butlerx/syphon-go/internal/pkg/parser"
	"github.com/golang/snappy"
	"github.com/lomik/carbon-clickhouse/helper/pb"
	"github.com/prometheus/client_golang/prometheus"
)

// PrometheusRemoteWrite https handler
type PrometheusRemoteWrite struct {
	http.Handler
	parseChan chan parser.PromMetric
	counter   prometheus.Counter
}

// HandlePrometheusConnection returns http hander for prometheus remote writes
func HandlePrometheusConnection(
	parseChan chan parser.PromMetric,
) *PrometheusRemoteWrite {
	mRec := receiversGauge.WithLabelValues("prometheus")
	mRec.Inc()
	mMsg := messageRecievedCounter.WithLabelValues("prometheus")
	return &PrometheusRemoteWrite{parseChan: parseChan, counter: mMsg}
}

func (rcv *PrometheusRemoteWrite) unpackFast(ctx context.Context, bufBody []byte) error {

	b := bufBody
	var err error
	var ts []byte
	var sample []byte

	metricBuffer := newPrometheusMetricBuffer()

	var metric []string
	var samplesOffset int

	var value float64
	var timestamp int64

TimeSeriesLoop:
	for len(b) > 0 {
		if b[0] != 0x0a { // repeated prometheus.TimeSeries timeseries = 1;
			if b, err = pb.Skip(b); err != nil {
				break TimeSeriesLoop
			}
			continue TimeSeriesLoop
		}

		if ts, b, err = pb.Bytes(b[1:]); err != nil {
			break TimeSeriesLoop
		}

		if metric, samplesOffset, err = metricBuffer.timeSeries(ts); err != nil {
			break TimeSeriesLoop
		}

		ts = ts[samplesOffset:]
	SamplesLoop:
		for len(ts) > 0 {
			if ts[0] != 0x12 { // repeated Sample samples = 2;
				if ts, err = pb.Skip(ts); err != nil {
					break TimeSeriesLoop
				}
				continue SamplesLoop
			}

			if sample, ts, err = pb.Bytes(ts[1:]); err != nil {
				break TimeSeriesLoop
			}

			timestamp = 0
			value = 0

			for len(sample) > 0 {
				switch sample[0] {
				case 0x09: // double value    = 1;
					if value, sample, err = pb.Double(sample[1:]); err != nil {
						break TimeSeriesLoop
					}
				case 0x10: // int64 timestamp = 2;
					if timestamp, sample, err = pb.Int64(sample[1:]); err != nil {
						break TimeSeriesLoop
					}
				default:
					if sample, err = pb.Skip(sample); err != nil {
						break TimeSeriesLoop
					}
				}
			}

			if math.IsNaN(value) {
				continue SamplesLoop
			}

			rcv.counter.Inc()
			rcv.parseChan <- parser.PromMetric{metric, value, timestamp / 1000}
		}
	}

	return err
}

func (rcv *PrometheusRemoteWrite) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	compressed, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	reqBuf, err := snappy.Decode(nil, compressed)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = rcv.unpackFast(r.Context(), reqBuf)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
