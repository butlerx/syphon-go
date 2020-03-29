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

// UDP output
func UDP(
	ctx context.Context,
	host string,
	port int64,
	pattern string,
	metric chan parser.Metric,
) {
	logger := zapwriter.Logger("uploader.UDP")
	m := uploaderGuage.WithLabelValues("UDP", pattern)
	m.Inc()
	defer m.Dec()

	address := fmt.Sprintf("%s:%s", host, strconv.FormatInt(port, 10))

	conn, err := connection.NewUDPConn(address)
	if err != nil {
		logger.Error(
			"Error opening UDP Connection",
			zap.String("destinationAddress", address),
			zap.Error(err),
		)
	}

	for {
		select {
		case <-ctx.Done():
			conn.Close()
			return
		case <-time.After(bufferSendTimer * time.Second):
			logger.Debug("Timer Executing")
			conn.SendBuffer()
		case m := <-metric:
			match, _ := regexp.MatchString(pattern, m.Path)
			if match {
				conn.AddMessage(m)
			}
		}
	}
}
