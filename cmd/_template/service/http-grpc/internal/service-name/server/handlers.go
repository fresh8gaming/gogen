package server

import (
	"context"

	healthPB "github.com/{{ .Org }}/{{ .Name }}/proto/grpc/health/v1"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

func RegisterHandlers(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) (err error) {
	return healthPB.RegisterHealthHandlerFromEndpoint(ctx, mux, endpoint, opts)
}

func RegisterServices(grpcServer grpc.ServiceRegistrar) {
	healthPB.RegisterHealthServer(grpcServer, HealthServer{})
}
