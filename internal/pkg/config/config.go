package config

import (
	"time"

	"github.com/lomik/zapwriter"
)

// MetricEndpointLocal text for local metric endpoint
const MetricEndpointLocal = "local"

// Config for syphon server
type Config struct {
	Metric     metricConfig       `toml:"metric"`
	Logging    []zapwriter.Config `toml:"logging"`
	File       fileConfig         `toml:"file"`
	Prometheus promConfig         `toml:"prometheus"`
	TCP        tcpConfig          `toml:"tcp"`
	UDP        udpConfig          `toml:"udp"`
	Uploader   uploaderConfig     `toml:"uploader"`
}

type metricConfig struct {
	Endpoint string    `toml:"endpoint"`
	Interval *Duration `toml:"interval"`
}

func newConfig() *Config {
	return &Config{
		Metric: metricConfig{
			Interval: &Duration{
				Duration: time.Minute,
			},
			Endpoint: MetricEndpointLocal,
		},
		Logging: nil,
		Uploader: uploaderConfig{
			File: []fileUploadConfig{{
				Enabled: true,
				Path:    "metrics_recieved.txt",
			}},
		},
		UDP: udpConfig{
			Listen:  ":2003",
			Mode:    "normal",
			Enabled: true,
		},
		TCP: tcpConfig{
			Listen:  ":2003",
			Enabled: true,
		},
		Prometheus: promConfig{
			Listen:  ":2006",
			Enabled: false,
		},
	}
}
