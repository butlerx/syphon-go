package receiver

import (
	"context"
	"net"
	"net/http"
	"runtime"
	"time"

	"github.com/butlerx/syphon-go/internal/pkg/handler"
	"github.com/butlerx/syphon-go/internal/pkg/parser"
	"github.com/lomik/zapwriter"
	"go.uber.org/zap"
)

// Prometheus http endpoint for Remote Write
func Prometheus(
	ctx context.Context,
	addr string,
	sendChannels *[]chan parser.Metric,
) {
	logger := zapwriter.Logger("reciever.Prometheus")

	tcpListener, err := net.Listen("tcp", addr)
	if err != nil {
		logger.Error(
			"Error creating TCP server",
			zap.String("address", addr),
			zap.Error(err),
		)

		return
	}

	logger.Info("server listening", zap.String("address", tcpListener.Addr().String()))
	defer tcpListener.Close()

	parseChan := make(chan parser.PromMetric)

	h := handler.HandlePrometheusConnection(parseChan)
	server := &http.Server{
		Handler:        h,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	if err := server.Serve(tcpListener); err != nil {
		logger.Error(
			"Failed to Start HTTP server",
			zap.String("address", tcpListener.Addr().String()),
			zap.Error(err),
		)
		return
	}

	for i := 0; i < runtime.NumCPU(); i++ {
		go parser.Prom(ctx, parseChan, sendChannels)
	}

	<-ctx.Done()
}
