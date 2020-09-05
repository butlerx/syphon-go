package receiver

import (
	"context"
	"net"
	"runtime"

	"github.com/butlerx/syphon-go/internal/pkg/handler"
	"github.com/butlerx/syphon-go/internal/pkg/parser"
	"github.com/lomik/zapwriter"
	"go.uber.org/zap"
)

// UDPServer Create a udp server on a given port.
func UDPServer(
	ctx context.Context,
	address string,
	sendChannels *[]chan parser.Metric,
) {
	logger := zapwriter.Logger("reciever.UDPServer")

	conn, err := net.ListenPacket("udp", address)
	if err != nil {
		logger.Error(
			"Error creating UDP server",
			zap.String("address", address),
			zap.Error(err),
		)

		return
	}

	logger.Info("server listening", zap.String("address", conn.LocalAddr().String()))
	defer conn.Close()

	parseChan := make(chan string)

	for i := 0; i < runtime.NumCPU(); i++ {
		go handler.HandleUDPConnection(
			ctx,
			conn,
			parseChan,
		)
	}

	for i := 0; i < runtime.NumCPU(); i++ {
		go parser.Plain(ctx, parseChan, sendChannels)
	}

	<-ctx.Done()
}
