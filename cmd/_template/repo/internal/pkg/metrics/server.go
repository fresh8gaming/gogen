package metrics

import (
	"fmt"
	"net/http"
	"time"

	"gitlab.sportradar.ag/ads/{{ .Team }}/{{ .Name }}/internal/pkg/logging"

	"go.uber.org/zap"
)

func StartServer(port string) {
	logger := logging.GetLogger()

	metricsServer := &http.Server{
		Addr:              fmt.Sprintf("0.0.0.0:%s", port),
		WriteTimeout:      time.Second * 5,
		ReadTimeout:       time.Second * 5,
		IdleTimeout:       time.Second * 60,
		ReadHeaderTimeout: time.Second * 2,
		Handler:           GetRouter(),
	}

	logger.Info("starting metrics on 9898", zap.String("port", port))

	go func() {
		if err := metricsServer.ListenAndServe(); err != nil {
			logger.Error(err.Error())
		}
	}()
}
