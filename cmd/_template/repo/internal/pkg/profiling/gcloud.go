package profiling

import (
	"os"

	"cloud.google.com/go/profiler"
)

// NewProfiler starts a new Google Go profiler using the given service name and version.
func NewProfiler(serviceName, serviceVersion string) error {
	return profiler.Start(profiler.Config{
		Service:        serviceName,
		ServiceVersion: serviceVersion,
		DebugLogging:   os.Getenv("PROFILER_DEBUG") == "true",
	})
}
