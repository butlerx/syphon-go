package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"sort"

	"github.com/butlerx/syphon-go/internal/app/syphon"
	"github.com/butlerx/syphon-go/internal/pkg/config"
	"github.com/lomik/zapwriter"
	"github.com/urfave/cli/v2"
)

var version = "development"

func main() {
	app := &cli.App{
		Name:      "syphon",
		Usage:     "Versatile metrics processor, proxy and forwarder",
		UsageText: "syphon [options] COMMAND",
		Description: `
			syphon is designed to accept and route metrics traffic.
			Metrics can be received from socket, snooped from live traffic or read from file or kafka.
			Metrics can be exportered via file, kafka or udp/tcp
		`,
		Authors: []*cli.Author{
			{
				Name:  "Cian Butler",
				Email: "butlerx@notthe.cloud",
			},
		},
		Version:              version,
		EnableBashCompletion: true,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "config, c",
				Usage: "Config file to use",
				Value: "",
			},
			&cli.BoolFlag{
				Name:  "print",
				Usage: "Print default config",
			},
		},
		Action: func(c *cli.Context) error {
			if c.Bool("print") {
				return config.PrintDefaultConfig()
			}

			cfg, err := config.ReadConfig(c.String("config"))
			if err != nil {
				return fmt.Errorf("can't load config: %v", err)
			}
			if err = zapwriter.ApplyConfig(cfg.Logging); err != nil {
				return fmt.Errorf("can't start logger: %v", err)
			}

			mainLogger := zapwriter.Logger("main")
			ctx := context.Background()
			listenChan := syphon.Uploader(ctx, cfg)
			syphon.Server(ctx, cfg, listenChan)
			mainLogger.Info("app started")

			<-ctx.Done()

			return nil
		},
	}

	sort.Sort(cli.FlagsByName(app.Flags))

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
