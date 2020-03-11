package main

import (
	"context"

	"github.com/butlerx/syphon/internal/pkg/config"
	"github.com/butlerx/syphon/internal/pkg/parser"
	"github.com/butlerx/syphon/internal/pkg/receiver"
	"github.com/butlerx/syphon/internal/pkg/uploader"
)

func startServer(
	ctx context.Context,
	cfg *config.Config,
	sendChannels []chan parser.Metric,
) {
	if cfg.TCP.Enabled {
		go receiver.TCPServer(ctx, cfg.TCP.Listen, &sendChannels)
	}

	if cfg.UDP.Enabled {
		if cfg.UDP.Mode == "promiscuous" {
			go receiver.RawUDPServer(ctx, cfg.UDP.Listen, &sendChannels)
		} else {
			go receiver.UDPServer(ctx, cfg.UDP.Listen, &sendChannels)
		}
	}

	if cfg.File.Enabled {
		go receiver.FileReader(ctx, cfg.File.Path, &sendChannels)
	}

	if cfg.Prometheus.Enabled {
		go receiver.Prometheus(ctx, cfg.Prometheus.Listen, &sendChannels)
	}
}

func startUploaders(ctx context.Context, cfg *config.Config) []chan parser.Metric {
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

	return sendChannels
}
