package server

import (
	"context"

	healthPB "github.com/{{ .Org }}/{{ .Name }}/proto/grpc/health/v1"
)

type HealthServer struct {
	healthPB.UnimplementedHealthServer
}

func (hs HealthServer) Check(
	ctx context.Context,
	request *healthPB.HealthCheckRequest,
) (*healthPB.HealthCheckResponse, error) {
	return &healthPB.HealthCheckResponse{
		Status: healthPB.HealthCheckResponse_SERVING,
	}, nil
}

func (hs HealthServer) Watch(*healthPB.HealthCheckRequest, healthPB.Health_WatchServer) error {
	return nil
}
