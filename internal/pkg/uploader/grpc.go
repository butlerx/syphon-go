package uploader

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/butlerx/syphon-go/internal/pkg/connection"
	"github.com/butlerx/syphon-go/internal/pkg/parser"
	"github.com/lomik/zapwriter"
	"go.uber.org/zap"
)

// Grpc outpu
func Grpc(
	ctx context.Context,
	host string,
	port int64,
	pattern string,
	metric chan parser.Metric,
) {
	logger := zapwriter.Logger("uploader.Grpc")
	m := uploaderGuage.WithLabelValues("GRPC", pattern)
	m.Inc()
	defer m.Dec()

	address := fmt.Sprintf("%s:%s", host, strconv.FormatInt(port, 10))

	conn, err := connection.NewGrpcConn(address)
	if err != nil {
		logger.Error(
			"Error opening grpc Connection",
			zap.String("destinationAddress", address),
			zap.Error(err),
		)
		return
	}

	for {
		select {
		case <-ctx.Done():
			conn.Close()
			return
		case <-time.After(bufferSendTimer * time.Second):
			logger.Debug("Timer Executing")
			if err := conn.SendBuffer(ctx); err != nil {
				logger.Error("Error sending grpc message", zap.Error(err))
			}
		case m := <-metric:
			match, _ := regexp.MatchString(pattern, m.Path)
			if match {
				conn.AddMessage(ctx, m)
			}
		}
	}
}
