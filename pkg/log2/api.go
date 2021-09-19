// Package log2 knows how to log
package log2

import (
	"fmt"
	"os"
	"strings"

	"go.uber.org/zap"
)

// New returns a new zap.SugaredLogger
func New() (*zap.SugaredLogger, error) {
	logType, ok := os.LookupEnv("LOG_TYPE")
	if !ok {
		logType = "JSON"
	}

	if strings.ToLower(logType) == "simple" {
		l, err := zap.NewDevelopment()
		if err != nil {
			return nil, fmt.Errorf("could not get dev logger: %w", err)
		}

		return l.Sugar(), nil
	}

	l, err := zap.NewProduction()
	if err != nil {
		return nil, fmt.Errorf("could not get prod logger: %w", err)
	}

	return l.Sugar(), nil
}
