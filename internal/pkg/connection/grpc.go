package connection

import (
	"context"
	"fmt"
	"reflect"
	"syscall"
	"time"

	"github.com/butlerx/syphon/internal/pkg/parser"
	api "github.com/lomik/carbon-clickhouse/grpc"
	"github.com/lomik/zapwriter"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

// GrpcConn manages grpc connection and buffer
type GrpcConn struct {
	logger    *zap.Logger
	conn      *grpc.ClientConn
	buffer    []*api.Metric
	address   string
	carbon    api.CarbonClient
	collector prometheus.Counter
}

// SendBuffer to grpc connection
func (u *GrpcConn) SendBuffer(ctx context.Context) error {
	if len(u.buffer) == 0 {
		u.logger.Debug("buffer empty", zap.String("destinationAddress", u.address))
		return nil
	}

	errChan := make(chan error)

	if _, err := u.carbon.Store(ctx, &api.Payload{Metrics: u.buffer}); err == syscall.EPIPE {
		go u.reconnect(errChan)
	} else if err == nil {
		u.buffer = make([]*api.Metric, 20)

		u.logger.Debug("message sent via grpc",
			zap.String("destinationAddress", u.address),
		)
		return nil
	}
	err := <-errChan
	if err != nil {
		return err
	}
	u.collector.Inc()
	return u.SendBuffer(ctx)
}

// Close ends grpc connection and empties buffer
func (u *GrpcConn) Close() {
	u.conn.Close()
	u.buffer = make([]*api.Metric, 20)
}

// AddMessage add a message to the grpc buffer
func (u *GrpcConn) AddMessage(ctx context.Context, m parser.Metric) {
	u.buffer = append(u.buffer, m.Grpc())
	bufferSize := int(uintptr(len(u.buffer)) * reflect.TypeOf(u.buffer).Elem().Size())
	u.logger.Debug("buffer size", zap.Int("bytes", bufferSize))
	if bufferSize > bufferSendSize {
		if err := u.SendBuffer(ctx); err != nil {
			u.logger.Info(
				"error sending message",
				zap.Int("bytes", bufferSize),
				zap.String("destinationAddress", u.address),
				zap.Error(err),
			)
		}
	}
}

func (u *GrpcConn) reconnect(errChan chan error) {
	u.conn.Close()
	conn, err := connectGrpc(u.address)
	if err != nil {
		errChan <- err
		return
	}
	u.conn = conn
	u.carbon = api.NewCarbonClient(conn)

	errChan <- nil
}

func connectGrpc(address string) (*grpc.ClientConn, error) {
	var err error
	for range []int{1, 2, 3, 4, 5} {
		conn, err := grpc.Dial(address)
		if err == nil {
			return conn, nil
		}

		time.Sleep(15 * time.Second)
	}

	return nil, fmt.Errorf("error connecting to grpc: %s", err)
}

// NewGrpcConn creates new grpc connection
func NewGrpcConn(address string) (*GrpcConn, error) {
	logger := zapwriter.Logger("connection.GrpcConn")

	conn, err := connectGrpc(address)
	if err != nil {
		return nil, err
	}

	carbon := api.NewCarbonClient(conn)
	buffer := make([]*api.Metric, 20)

	return &GrpcConn{
		logger,
		conn,
		buffer,
		address,
		carbon,
		messageSentCounter.WithLabelValues("GRPC", address),
	}, nil
}
