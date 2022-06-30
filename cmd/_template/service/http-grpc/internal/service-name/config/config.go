package config

import (
	"flag"

	"github.com/{{ .Org }}/{{ .Name }}/internal/pkg/env"
	"github.com/{{ .Org }}/{{ .Name }}/internal/pkg/logging"

	"go.uber.org/zap"
)

// Config stores all the values required throughout various parts of the system.
type Config struct {
	GcpProject        string
	Port              string
	MPort             string
	MuteHealthLogging bool
	ProfilerEnabled   bool
}

var (
	config Config
)

// Initialise sets all the values required for a service config, overwriting with
// envvars and CLI flags.
func Initialise() {
	logger := logging.GetLogger()

	var (
		muteHealthLogging, profilerEnabled bool
		gcpProject, port, mPort            string
	)

	flag.StringVar(&gcpProject, "gcpProject", env.GetenvString("GCP_PROJECT", "local-project"), "gcp project used for pubsub")
	flag.StringVar(&port, "port", env.GetenvString("PORT", "8080"), "port for http/grpc service")
	flag.StringVar(&mPort, "mPort", env.GetenvString("MPORT", "9898"), "metrics port for http/grpc service")
	flag.BoolVar(&muteHealthLogging, "muteHealthLogging", env.GetenvBool("MUTE_HEATLH_LOGGING", true), "mute grpc logging for health checks")
	flag.BoolVar(&profilerEnabled, "profilerEnabled", env.GetenvBool("PROFILER_ENABLED", false), "enable the google cloud profiler")
	flag.Parse()

	logger.Info("configuration",
		zap.String("gcpProject", gcpProject),
		zap.String("port", port),
		zap.String("mPort", mPort),
		zap.Bool("muteHealthLogging", muteHealthLogging),
		zap.Bool("profilerEnabled", profilerEnabled),
	)

	config.GcpProject = gcpProject
	config.Port = port
	config.MPort = mPort
	config.MuteHealthLogging = muteHealthLogging
	config.ProfilerEnabled = profilerEnabled
}

func Get() Config {
	return config
}
