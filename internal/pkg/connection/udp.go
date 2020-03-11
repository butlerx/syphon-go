package connection

import (
	"bytes"
	"fmt"
	"net"

	"github.com/butlerx/syphon/internal/pkg/parser"
	"github.com/lomik/zapwriter"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"
)

// UDPConn manages udp connection and buffer
type UDPConn struct {
	logger    *zap.Logger
	conn      *net.UDPConn
	buffer    *bytes.Buffer
	address   string
	collector prometheus.Counter
}

// SendBuffer to udp connection
func (u *UDPConn) SendBuffer() {
	if u.buffer.Len() == 0 {
		u.logger.Debug("buffer empty", zap.String("destinationAddress", u.address))

		return
	}
	if _, err := u.conn.Write(u.buffer.Bytes()); err != nil {
		u.logger.Info(
			"error sending message",
			zap.String("destinationAddress", u.address),
			zap.Error(err),
		)

		return
	}
	u.buffer.Reset()
	u.collector.Inc()
	u.logger.Debug("message sent via UDP",
		zap.String("destinationAddress", u.address),
	)
}

// Close ends udp connection and empties buffer
func (u *UDPConn) Close() {
	u.conn.Close()
	u.buffer.Reset()
}

// AddMessage add a message to the udp buffer
func (u *UDPConn) AddMessage(m parser.Metric) {
	u.buffer.WriteString(m.String())
	u.logger.Debug("buffer size", zap.Int("bytes", u.buffer.Len()))
	if u.buffer.Len() > bufferSendSize {
		u.SendBuffer()
	}
}

// NewUDPConn creates new udp connection
func NewUDPConn(address string) (*UDPConn, error) {
	logger := zapwriter.Logger("connection.UDPConn")

	destinationAddress, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		return nil, fmt.Errorf("error ResolvingUDPAddr: %s", err)
	}

	conn, err := net.DialUDP("udp", nil, destinationAddress)
	if err != nil {
		return nil, fmt.Errorf("error connecting to udp: %s", err)
	}

	buffer := bytes.NewBuffer(make([]byte, 1024))
	buffer.Reset()

	return &UDPConn{
		logger,
		conn,
		buffer,
		address,
		messageSentCounter.WithLabelValues("UDP", address),
	}, nil
}
