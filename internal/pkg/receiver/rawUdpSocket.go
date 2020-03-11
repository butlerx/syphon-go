package receiver

import (
	"context"
	"fmt"
	"net"
	"os"
	"runtime"
	"strconv"
	"strings"
	"syscall"

	"github.com/butlerx/syphon/internal/pkg/handler"
	"github.com/butlerx/syphon/internal/pkg/parser"
	"github.com/lomik/zapwriter"
	"go.uber.org/zap"
)

// RawUDPServer binds to a port and listens to udp packets
// TODO filter by port currectly
func RawUDPServer(
	ctx context.Context,
	address string,
	sendChannels *[]chan parser.Metric,
) {
	logger := zapwriter.Logger("reciever.RawUDPServer")

	socket, err := bindSocket(address)
	if err != nil {
		logger.Error(
			"Error Listening promiscuously to udp",
			zap.String("address", address),
			zap.Error(err),
		)

		return
	}
	defer syscall.Close(socket)
	logger.Info("server listening", zap.String("address", address))

	file := os.NewFile(uintptr(socket), fmt.Sprintf("fd %d", socket))
	parseChan := make(chan string)

	go (func() {
		for {
			buf := make([]byte, 1024)
			numRead, err := file.Read(buf)

			if err != nil {
				logger.Warn(
					"Error parsing udp buffer",
					zap.String("address", address),
					zap.Error(err),
				)
			}

			go handler.HandleRawUDPPacket(
				ctx,
				buf[:numRead],
				parseChan,
			)
		}
	})()

	for i := 0; i < runtime.NumCPU(); i++ {
		go parser.Plain(ctx, parseChan, sendChannels)
	}

	<-ctx.Done()
}

func bindSocket(address string) (int, error) {
	socket, err := syscall.Socket(
		syscall.AF_INET,
		syscall.SOCK_RAW,
		syscall.IPPROTO_UDP,
	)
	if err != nil {
		return 0, fmt.Errorf("error creating socket: %s", err)
	}

	if err := syscall.SetsockoptInt(
		socket,
		syscall.SOL_SOCKET,
		syscall.SO_REUSEADDR,
		1,
	); err != nil {
		return 0, fmt.Errorf("error setting options on socket: %s", err)
	}

	addr := getAddr(address)

	if err := syscall.Bind(socket, &addr); err != nil {
		return 0, fmt.Errorf("error binding socket: %s", err)
	}

	return socket, nil
}

func getAddr(address string) syscall.SockaddrInet4 {
	hostPort := strings.SplitN(address, ":", 2)

	var host string
	if hostPort[0] != "" {
		host = hostPort[0]
	} else {
		host = "127.0.0.1"
	}

	port, _ := strconv.Atoi(hostPort[1])
	addr := syscall.SockaddrInet4{Port: port}
	copy(addr.Addr[:], net.ParseIP(host).To4())

	return addr
}
