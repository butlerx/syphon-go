package handler

import (
	"bufio"
	"context"
	"net"
	"strings"
	"time"

	"github.com/lomik/zapwriter"
	"go.uber.org/zap"
)

// HandleTCPConnection process messages from tcp connection
func HandleTCPConnection(
	ctx context.Context,
	conn net.Conn,
	parseChan chan string,
) {
	logger := zapwriter.Logger("handler.HandleTCPConnection")
	mRec := receiversGauge.WithLabelValues("TCP")
	mRec.Inc()
	defer mRec.Dec()
	mMsg := messageRecievedCounter.WithLabelValues("TCP")
	logger.Debug("Serving Connection", zap.String("peer", conn.RemoteAddr().String()))

	defer conn.Close()

	finished := make(chan bool)
	defer close(finished)

	go (func() {
		select {
		case <-finished:
			conn.Close()
			return
		case <-ctx.Done():
			conn.Close()
			return
		}
	})()

	for {
		if err := conn.SetReadDeadline(time.Time{}); err != nil {
			logger.Warn(
				"Failed to set read deadline for TCP connection",
				zap.String("peer", conn.RemoteAddr().String()),
				zap.Error(err),
			)
		}

		netData, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			logger.Error(
				"Error reading from TCP connection",
				zap.String("peer", conn.RemoteAddr().String()),
				zap.Error(err),
			)
			finished <- true
		}

		data := strings.TrimSpace(netData)
		logger.Debug(
			"Message Received",
			zap.String("data", data),
			zap.String("peer", conn.RemoteAddr().String()),
		)
		mMsg.Inc()
		parseChan <- data
	}
}
