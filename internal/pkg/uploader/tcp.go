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

// TCP output
func TCP(
	ctx context.Context,
	host string,
	port int64,
	pattern string,
	metric chan parser.Metric,
) {
	logger := zapwriter.Logger("uploader.TCP")
	m := uploaderGuage.WithLabelValues("TCP", pattern)
	m.Inc()
	defer m.Dec()

	address := fmt.Sprintf("%s:%s", host, strconv.FormatInt(port, 10))

	conn, err := connection.NewTCPConn(address)
	if err != nil {
		logger.Error(
			"Error opening TCP Connection",
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
			if err := conn.SendBuffer(); err != nil {
				logger.Error("Error sending tcp message", zap.Error(err))
			}
		case m := <-metric:
			match, _ := regexp.MatchString(pattern, m.Path)
			if match {
				conn.AddMessage(m)
			}
		}
	}
}
