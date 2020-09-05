package receiver

import (
	"context"
	"os"
	"runtime"

	"github.com/butlerx/syphon-go/internal/pkg/handler"
	"github.com/butlerx/syphon-go/internal/pkg/parser"
	"github.com/lomik/zapwriter"
	"go.uber.org/zap"
)

// FileReader loads from file.
func FileReader(
	ctx context.Context,
	path string,
	sendChannels *[]chan parser.Metric,
) {
	logger := zapwriter.Logger("reciever.FileReader")

	file, err := os.Open(path)
	if err != nil {
		logger.Error(
			"Error opening file",
			zap.String("file", path),
			zap.Error(err),
		)

		return
	}

	logger.Info("Reading from file", zap.String("file", path))

	defer file.Close()

	parseChan := make(chan string)

	for i := 0; i < runtime.NumCPU(); i++ {
		go handler.HandleFile(ctx, file, parseChan)
	}

	for i := 0; i < runtime.NumCPU(); i++ {
		go parser.Plain(ctx, parseChan, sendChannels)
	}

	<-ctx.Done()
}
