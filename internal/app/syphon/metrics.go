package syphon

import (
	"context"
	"time"

	"github.com/butlerx/syphon-go/internal/pkg/config"
	"github.com/butlerx/syphon-go/internal/pkg/connection"
	"github.com/butlerx/syphon-go/internal/pkg/handler"
	"github.com/butlerx/syphon-go/internal/pkg/parser"
	"github.com/butlerx/syphon-go/internal/pkg/uploader"
	"github.com/lomik/zapwriter"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"
)

func registerMetrics(
	ctx context.Context,
	interval *config.Duration,
	sendChannels []chan parser.Metric,
) {
	logger := zapwriter.Logger("metric_collector")

	handler.RegisterMetrics()
	parser.RegisterMetrics()
	uploader.RegisterMetrics()
	connection.RegisterMetrics()

	gatherer := prometheus.DefaultGatherer

	ticker := time.NewTicker(interval.Value())
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := push(gatherer, sendChannels); err != nil {
				logger.Error("error pushing Internal Metrics", zap.Error(err))
			}
		case <-ctx.Done():
			if err := push(gatherer, sendChannels); err != nil {
				logger.Error("error pushing Internal Metrics", zap.Error(err))
			}

			return
		}
	}
}

// TODO parse prom metrics and send them to channels.
func push(gatherer prometheus.Gatherer, sendChannels []chan parser.Metric) error {
	mfs, err := gatherer.Gather()
	if err != nil || len(mfs) == 0 {
		return err
	}
	//for metric := range parser.ParseProm(mfs) {
	//	for channel := range sendChannels {
	//		channel <- metric
	//	}
	//}
	return nil
}
