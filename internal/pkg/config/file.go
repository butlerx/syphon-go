package config

import (
	"fmt"
	"io/ioutil"
	"net/url"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/lomik/zapwriter"
)

// ReadConfig from file
func ReadConfig(filename string) (*Config, error) {
	cfg := newConfig()

	if filename != "" {
		b, err := ioutil.ReadFile(filename)
		if err != nil {
			return nil, err
		}

		body := string(b)

		// @TODO: fix for config starts with [logging]
		body = strings.Replace(body, "\n[logging]\n", "\n[[logging]]\n", -1)

		if _, err := toml.Decode(body, cfg); err != nil {
			return nil, err
		}
	}

	if cfg.Logging == nil {
		cfg.Logging = make([]zapwriter.Config, 0)
	}

	if len(cfg.Logging) == 0 {
		cfg.Logging = append(cfg.Logging, zapwriter.NewConfig())
	}

	if err := zapwriter.CheckConfig(cfg.Logging, nil); err != nil {
		return nil, err
	}

	if cfg.Metric.Endpoint == "" {
		cfg.Metric.Endpoint = MetricEndpointLocal
	}

	if cfg.Metric.Endpoint != MetricEndpointLocal {
		u, err := url.Parse(cfg.Metric.Endpoint)
		if err != nil {
			return nil, fmt.Errorf("common.metric-endpoint parse error: %s", err.Error())
		}

		if u.Scheme != "tcp" && u.Scheme != "udp" {
			return nil, fmt.Errorf("common.metric-endpoint supports only tcp and udp protocols. %#v is unsupported", u.Scheme)
		}
	}
	return cfg, nil
}
