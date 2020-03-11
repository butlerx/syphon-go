package handler

import (
	"bytes"
	"context"
	"net"
	"strings"

	"github.com/lomik/zapwriter"
	"go.uber.org/zap"
)

const maxBufferSize = 1024

// HandleUDPConnection processes messages from udp connections
func HandleUDPConnection(
	ctx context.Context,
	conn net.PacketConn,
	parseChan chan string,
) {
	defer conn.Close()
	mRec := receiversGauge.WithLabelValues("UDP")
	mRec.Inc()
	defer mRec.Dec()
	mMsg := messageRecievedCounter.WithLabelValues("UDP")

	logger := zapwriter.Logger("handler.HandleUDPConnection")
	buffer := make([]byte, maxBufferSize)

ReceiveLoop:
	for {
		n, peer, err := conn.ReadFrom(buffer)
		if err != nil {
			if strings.Contains(err.Error(), "use of closed network connection") {
				break ReceiveLoop
			}
			logger.Error("ReadFrom failed", zap.Error(err), zap.String("peer", peer.String()))
			continue ReceiveLoop
		}

		logger.Debug(
			"Message Received",
			zap.String("data", string(buffer[:n])),
			zap.String("peer", peer.String()),
		)

		if n > 0 {
			if chunkSize := bytes.LastIndexByte(buffer[:n], '\n') + 1; chunkSize < n {
			} else if chunkSize > 0 {
				mMsg.Inc()
				parseChan <- string(buffer)
				buffer = make([]byte, maxBufferSize)
			}
		}
	}
}
