package handler

import (
	"context"

	"github.com/lomik/zapwriter"
	"go.uber.org/zap"
)

//HandleRawUDPPacket process udp packets in the form of byte arrays
func HandleRawUDPPacket(
	ctx context.Context,
	rawData []byte,
	parseChan chan string,
) {
	logger := zapwriter.Logger("handler.HandleRawUDPConnection")
	mRec := receiversGauge.WithLabelValues("RawUDP")
	mRec.Inc()
	defer mRec.Dec()
	mMsg := messageRecievedCounter.WithLabelValues("RawUDP")

	data := rawData[14:]
	headerLength := (data[0] & 15) * 4
	message := string(data[headerLength:][14:])
	logger.Debug("Message received", zap.String("data", message))
	mMsg.Inc()
	parseChan <- message
}
