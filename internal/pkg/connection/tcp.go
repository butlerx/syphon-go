package connection

import (
	"bytes"
	"fmt"
	"net"
	"syscall"
	"time"

	"github.com/butlerx/syphon-go/internal/pkg/parser"
	"github.com/lomik/zapwriter"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"
)

// TCPConn manages tcp connection and buffer
type TCPConn struct {
	logger    *zap.Logger
	conn      *net.TCPConn
	buffer    *bytes.Buffer
	address   string
	collector prometheus.Counter
}

// SendBuffer to tcp connection
func (u *TCPConn) SendBuffer() error {
	if u.buffer.Len() == 0 {
		u.logger.Debug("buffer empty", zap.String("destinationAddress", u.address))
		return nil
	}

	errChan := make(chan error)

	if _, err := u.conn.Write(u.buffer.Bytes()); err == syscall.EPIPE {
		go u.reconnect(errChan)
	} else if err == nil {
		u.buffer.Reset()
		u.logger.Debug("message sent via TCP",
			zap.String("destinationAddress", u.address),
		)
		return nil
	}
	err := <-errChan
	if err != nil {
		return err
	}
	u.collector.Inc()
	return u.SendBuffer()
}

// Close ends tcp connection and empties buffer
func (u *TCPConn) Close() {
	u.conn.Close()
	u.buffer.Reset()
}

// AddMessage add a message to the tcp buffer
func (u *TCPConn) AddMessage(m parser.Metric) {
	u.buffer.WriteString(m.String())
	u.logger.Debug("buffer size", zap.Int("bytes", u.buffer.Len()))
	if u.buffer.Len() > bufferSendSize {
		if err := u.SendBuffer(); err != nil {
			u.logger.Info(
				"error sending message",
				zap.String("destinationAddress", u.address),
				zap.Error(err),
			)
		}
	}
}

func (u *TCPConn) reconnect(errChan chan error) {
	u.conn.Close()
	conn, err := connectTCP(u.address)
	if err != nil {
		errChan <- err
		return
	}
	u.conn = conn
	errChan <- nil
}

func connectTCP(address string) (*net.TCPConn, error) {
	destinationAddress, err := net.ResolveTCPAddr("tcp", address)
	if err != nil {
		return nil, fmt.Errorf("error ResolvingTCPAddr: %s", err)
	}

	for range []int{1, 2, 3, 4, 5} {
		conn, err := net.DialTCP("tcp", nil, destinationAddress)
		if err == nil {
			return conn, nil
		}

		time.Sleep(15 * time.Second)
	}

	return nil, fmt.Errorf("error connecting to tcp: %s", err)
}

// NewTCPConn creates new tcp connection
func NewTCPConn(address string) (*TCPConn, error) {
	logger := zapwriter.Logger("connection.TCPConn")

	conn, err := connectTCP(address)
	if err != nil {
		return nil, err
	}

	buffer := bytes.NewBuffer(make([]byte, 1024))
	buffer.Reset() // Sometimes this buffer has not been empty

	return &TCPConn{
		logger,
		conn,
		buffer,
		address,
		messageSentCounter.WithLabelValues("TCP", address),
	}, nil
}
