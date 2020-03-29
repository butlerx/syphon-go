package syphon

import (
	"context"

	"github.com/butlerx/syphon-go/internal/pkg/config"
	"github.com/butlerx/syphon-go/internal/pkg/parser"
	"github.com/butlerx/syphon-go/internal/pkg/uploader"
)

// Uploader init function
func Uploader(ctx context.Context, cfg *config.Config) []chan parser.Metric {
	var sendChannels []chan parser.Metric

	for _, u := range cfg.Uploader.UDP {
		if u.Enabled {
			sendChan := make(chan parser.Metric)
			sendChannels = append(sendChannels, sendChan)

			if u.Pattern == "" {
				u.Pattern = ".*"
			}

			go uploader.UDP(ctx, u.Host, u.Port, u.Pattern, sendChan)
		}
	}

	for _, u := range cfg.Uploader.File {
		if u.Enabled {
			sendChan := make(chan parser.Metric)
			sendChannels = append(sendChannels, sendChan)

			if u.Pattern == "" {
				u.Pattern = ".*"
			}

			go uploader.File(ctx, u.Path, u.Pattern, sendChan)
		}
	}

	for _, u := range cfg.Uploader.TCP {
		if u.Enabled {
			sendChan := make(chan parser.Metric)
			sendChannels = append(sendChannels, sendChan)

			if u.Pattern == "" {
				u.Pattern = ".*"
			}

			go uploader.TCP(ctx, u.Host, u.Port, u.Pattern, sendChan)
		}
	}

	for _, u := range cfg.Uploader.Grpc {
		if u.Enabled {
			sendChan := make(chan parser.Metric)
			sendChannels = append(sendChannels, sendChan)

			if u.Pattern == "" {
				u.Pattern = ".*"
			}

			go uploader.Grpc(ctx, u.Host, u.Port, u.Pattern, sendChan)
		}
	}

	go registerMetrics(ctx, cfg.Metric.Interval, sendChannels)
	return sendChannels
}
