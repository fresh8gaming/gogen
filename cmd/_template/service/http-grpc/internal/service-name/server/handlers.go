package server

import (
	"context"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthPB "google.golang.org/grpc/health/grpc_health_v1"
)

func RegisterHandlers(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) (err error) {
	return nil
}

func RegisterServices(grpcServer grpc.ServiceRegistrar, healthServer *health.Server) {
	healthPB.RegisterHealthServer(grpcServer, healthServer)
	healthServer.SetServingStatus("ready", healthPB.HealthCheckResponse_NOT_SERVING)
}
