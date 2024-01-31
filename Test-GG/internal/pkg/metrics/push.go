package metrics

import (
	"time"

	"gitlab.sportradar.ag/ads/dmp/Test-GG/internal/pkg/logging"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/push"
	"go.uber.org/zap"
)

var (
	duration = promauto.NewHistogram(prometheus.HistogramOpts{
		Name:    "cronjob_duration",
		Help:    "Duration of cronjob",
		Buckets: prometheus.DefBuckets,
	})

	globalPusher *push.Pusher
)

func SetPusher(pushGatewayURL, service string) {
	globalPusher = push.New(pushGatewayURL, service)
}

func GetPusher() *push.Pusher {
	return globalPusher
}

func WriteDurationToPushGateway(start time.Time) {
	if globalPusher == nil {
		return
	}

	logger := logging.GetLogger()

	duration.Observe(float64(time.Since(start) / time.Millisecond))

	logger.Debug("writing data to push gateway",
		zap.String("duration", time.Since(start).String()))

	err := globalPusher.Collector(duration).Push()
	if err != nil {
		logger.Error("failed to write data to push gateway",
			zap.Error(err),
			zap.String("duration", time.Since(start).String()))
	}
}
