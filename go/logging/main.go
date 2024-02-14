package main

import (
	"log"
	"log/slog"
	"os"

	"go.uber.org/zap"
)

const (
	// Simulate information from an HTTP request.
	method = "GET"
	path   = "/api/v1/users"
)

func main() {
	// log
	stdLog := log.New(os.Stdout, "", log.LstdFlags)
	stdLog.Printf("requests received: method=%s path=%s", method, path)

	// log/slog
	stdSlog := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{}))
	stdSlog.Info("request received:",
		"method", method,
		"path", path,
	)

	// go.uber.org/zap
	zapLogger, _ := zap.NewProduction()
	defer zapLogger.Sync()
	zapLogger.Info("request received:",
		zap.String("method", method),
		zap.String("path", path),
	)
}
