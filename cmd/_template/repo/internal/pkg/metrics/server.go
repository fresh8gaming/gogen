package metrics

import (
	"net/http"
	"time"

	"github.com/{{ .Org }}/{{ .Name }}/internal/pkg/logging"

	"go.uber.org/zap"
)

func StartServer() {
	logger := logging.GetLogger()

	metricsServer := &http.Server{
		Addr:         fmt.Sprintf("0.0.0.0:%s", config.Get().MPort),
		WriteTimeout: time.Second * 5,
		ReadTimeout:  time.Second * 5,
		IdleTimeout:  time.Second * 60,
		Handler:      GetRouter(),
	}

	logger.Info("starting metrics on 9898", zap.String(config.Get().MPort))

	go func() {
		if err := metricsServer.ListenAndServe(); err != nil {
			logger.Error(err.Error())
		}
	}()
}

