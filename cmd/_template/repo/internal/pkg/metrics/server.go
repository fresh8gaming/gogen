package metrics

import (
	"net/http"
	"time"

	"github.com/{{ .Org }}/{{ .Name }}/internal/pkg/logging"
)

func StartServer() {
	logger := logging.GetLogger()

	logger.Debug("setup metrics")

	metricsServer := &http.Server{
		Addr:         "0.0.0.0:9898",
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
