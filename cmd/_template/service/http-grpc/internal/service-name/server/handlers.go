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

func RegisterServices(grpcServer grpc.ServiceRegistrar) {
	healthServer := health.NewServer()
	healthPB.RegisterHealthServer(grpcServer, healthServer)
	// should possibly use full "service" name here but health will do
	healthServer.SetServingStatus("health", healthPB.HealthCheckResponse_SERVING)
}
