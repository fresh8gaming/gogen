package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/{{ .Org }}/{{ .Name }}/internal/pkg/logging"
	"github.com/{{ .Org }}/{{ .Name }}/internal/pkg/metrics"
	"github.com/{{ .Org }}/{{ .Name }}/internal/pkg/profiling"
	"github.com/{{ .Org }}/{{ .Name }}/internal/{{ .ServiceName }}/config"
	"github.com/{{ .Org }}/{{ .Name }}/internal/{{ .ServiceName }}/process"
)

var (
	Version = "dev"
	App     = "{{ .ServiceName }}"
)

func main() {
	setup()

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	start := time.Now()
	defer func() {
		metrics.WriteDurationToPushGateway(start)
		cancel()
	}()

	process.Work(ctx)
}

func setup() {
	logger, err := logging.NewLogger(os.Getenv("ENV"))
	if err != nil {
		log.Fatalf("Cannot set up logger: %s", err.Error())
	}

	logging.SetLogger(logger)

	config.Initialise()

	if config.Get().ProfilerEnabled {
		err = profiling.NewProfiler(App, Version)
		if err != nil {
			logger.Error(err.Error())
		}
	}

	if config.Get().PushMetricsEnabled {
		metrics.SetPusher(config.Get().PushGatewayAddress, App)
	}
}
