package config

import (
	"bytes"
	"fmt"

	"github.com/BurntSushi/toml"
	"github.com/lomik/zapwriter"
)

// PrintDefaultConfig  to terminal.
func PrintDefaultConfig() error {
	cfg := newConfig()

	if cfg.Logging == nil {
		cfg.Logging = make([]zapwriter.Config, 0)
	}

	if len(cfg.Logging) == 0 {
		cfg.Logging = append(cfg.Logging, zapwriter.NewConfig())
	}

	cfg.Uploader.File[0].Pattern = ".*"
	buf := new(bytes.Buffer)
	encoder := toml.NewEncoder(buf)
	encoder.Indent = ""

	if err := encoder.Encode(cfg); err != nil {
		return err
	}

	fmt.Print(buf.String())

	return nil
}
