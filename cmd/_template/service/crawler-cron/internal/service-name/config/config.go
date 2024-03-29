package config

import (
	"flag"

	"gitlab.sportradar.ag/ads/{{ .Team }}/{{ .Name }}/internal/pkg/env"
	"gitlab.sportradar.ag/ads/{{ .Team }}/{{ .Name }}/internal/pkg/logging"

	"go.uber.org/zap"
)

// Config stores all the values required throughout various parts of the system.
type Config struct {
	GcpProject         string
	PushGatewayAddress string
	ProfilerEnabled    bool
	PushMetricsEnabled bool
	Inplay             bool
}

var (
	config Config
)

// Initialise sets all the values required for a service config, overwriting with
// envvars and CLI flags.
func Initialise() {
	logger := logging.GetLogger()

	var (
		profilerEnabled, pushMetricsEnabled, inplay bool
		gcpProject, pushGatewayAddress              string
	)

	flag.StringVar(&gcpProject, "gcpProject", env.GetenvString("GCP_PROJECT", "local-project"), "gcp project")
	flag.StringVar(
		&pushGatewayAddress,
		"pushGatewayAddress",
		env.GetenvString("PUSH_GATEWAY_ADDRESS", "prometheus-push-prometheus-pushgateway.monitoring.svc.cluster.local:9091"),
		"",
	)
	flag.BoolVar(
		&profilerEnabled,
		"profilerEnabled",
		env.GetenvBool("PROFILER_ENABLED", false),
		"enable the google cloud profiler",
	)
	flag.BoolVar(
		&pushMetricsEnabled,
		"pushMetricsEnabled",
		env.GetenvBool("PUSH_METRICS_ENABLED", false),
		"enable the prometheus push metrics",
	)
	flag.BoolVar(
		&inplay,
		"inplay",
		env.GetenvBool("INPLAY", false),
		"inplay service",
	)
	flag.Parse()

	logger.Info("configuration",
		zap.String("gcpProject", gcpProject),
		zap.String("pushGatewayAddress", pushGatewayAddress),
		zap.Bool("profilerEnabled", profilerEnabled),
		zap.Bool("pushMetricsEnabled", pushMetricsEnabled),
	)

	config.GcpProject = gcpProject
	config.PushGatewayAddress = pushGatewayAddress
	config.ProfilerEnabled = profilerEnabled
	config.PushMetricsEnabled = pushMetricsEnabled
	config.Inplay = inplay
}

func Get() Config {
	return config
}
