package syphon

import (
	"context"

	"github.com/butlerx/syphon/internal/pkg/config"
	"github.com/butlerx/syphon/internal/pkg/parser"
	"github.com/butlerx/syphon/internal/pkg/receiver"
)

// Server init Function
func Server(
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
