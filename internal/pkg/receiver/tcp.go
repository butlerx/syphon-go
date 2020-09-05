package receiver

import (
	"context"
	"net"
	"runtime"
	"strings"

	"github.com/butlerx/syphon-go/internal/pkg/handler"
	"github.com/butlerx/syphon-go/internal/pkg/parser"
	"github.com/lomik/zapwriter"
	"go.uber.org/zap"
)

// TCPServer Create a tcp server on a given port
// TODO handle closed connections and EOF more gracefully.
func TCPServer(
	ctx context.Context,
	addr string,
	sendChannels *[]chan parser.Metric,
) {
	logger := zapwriter.Logger("reciever.TCPServer")

	tcpListener, err := net.Listen("tcp4", addr)
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

	parseChan := make(chan string)

	go (func() {
		defer tcpListener.Close()
		for {
			conn, err := tcpListener.Accept()
			if err != nil {
				if strings.Contains(err.Error(), "use of closed network connection") {
					break
				}
				logger.Warn(
					"failed to accept connection",
					zap.String("address", tcpListener.Addr().String()),
					zap.Error(err),
				)

				continue
			}
			go handler.HandleTCPConnection(
				ctx,
				conn,
				parseChan,
			)
		}
	})()

	for i := 0; i < runtime.NumCPU(); i++ {
		go parser.Plain(ctx, parseChan, sendChannels)
	}

	<-ctx.Done()
}
