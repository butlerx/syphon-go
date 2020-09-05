package uploader

import (
	"context"
	"os"
	"regexp"

	"github.com/butlerx/syphon-go/internal/pkg/parser"
	"github.com/lomik/zapwriter"
	"go.uber.org/zap"
)

// File output
func File(ctx context.Context, path string, pattern string, metric chan parser.Metric) {
	logger := zapwriter.Logger("uploader.File")
	m := uploaderGuage.WithLabelValues("file", pattern)
	m.Inc()
	defer m.Dec()

	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		file.Close()
		logger.Error("Error opening file", zap.String("path", path), zap.Error(err))
	}

	reg := regexp.MustCompile(pattern)

	for {
		select {
		case <-ctx.Done():
			file.Close()
			return
		case m := <-metric:
			match := reg.MatchString(m.Path)
			if match {
				_, err := file.WriteString(m.String())
				if err != nil {
					logger.Error("error writing to file", zap.String("path", path), zap.Error(err))
				} else {
					logger.Debug("message written to file", zap.String("path", path), zap.String("metric", m.String()))
				}
			}
		}
	}
}
