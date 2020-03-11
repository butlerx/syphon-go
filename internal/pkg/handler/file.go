package handler

import (
	"bufio"
	"context"
	"io"

	"github.com/lomik/zapwriter"
	"go.uber.org/zap"
)

// HandleFile Loads a file and reads line by line
func HandleFile(
	ctx context.Context,
	file io.Reader,
	parseChan chan string,
) {
	mRec := receiversGauge.WithLabelValues("file")
	mRec.Inc()
	defer mRec.Dec()
	mMsg := messageRecievedCounter.WithLabelValues("file")
	logger := zapwriter.Logger("handler.HandleFile")
	reader := bufio.NewReader(file)

ReceiveLoop:
	for {
		line, _, err := reader.ReadLine()
		if err == io.EOF {
			break ReceiveLoop
		} else if err != nil {
			logger.Error("ReadFrom failed", zap.Error(err))
			continue ReceiveLoop
		}

		s := string(line)
		logger.Debug(
			"line read",
			zap.String("line", s),
		)

		if len(line) > 0 {
			mMsg.Inc()
			parseChan <- s
		}
	}
}
