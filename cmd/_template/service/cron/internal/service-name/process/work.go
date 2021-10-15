package process

import (
	"context"

	"github.com/{{ .Org }}/{{ .Name }}/internal/pkg/logging"
)

func Work(ctx context.Context) {
	logger := logging.GetLogger()
	logger.Info("do work")
}
