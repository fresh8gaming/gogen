package process

import (
	"context"

	"gitlab.sportradar.ag/ads/{{ .Team }}/{{ .Name }}/internal/pkg/logging"
)

func Work(ctx context.Context) {
	logger := logging.GetLogger()
	logger.Info("do work")
}
