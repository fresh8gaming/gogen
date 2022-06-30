package metrics

import (
	"net/http"
	"time"

	"github.com/{{ .Org }}/{{ .Name }}/internal/pkg/logging"
)

const defaultMetricsPort = "9898"

func StartServer(ports ...string) {
	logger := logging.GetLogger()
	port := getMetricsPort(ports)
	logger.Debug("setup metrics", zap.String("port", port))

	metricsServer := &http.Server{
		Addr:         fmt.Sprintf("0.0.0.0:%s", port),
		WriteTimeout: time.Second * 5,
		ReadTimeout:  time.Second * 5,
		IdleTimeout:  time.Second * 60,
		Handler:      GetRouter(),
	}

	logger.Info("starting metrics on 9898")

	go func() {
		if err := metricsServer.ListenAndServe(); err != nil {
			logger.Error(err.Error())
		}
	}()
}

func getMetricsPort(in []string) string {
	for _, s := range in {
		return s
	}
	return defaultMetricsPort
}
