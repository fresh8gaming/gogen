package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/{{ .Org }}/{{ .Name }}/internal/{{ .ServiceName }}/config"
	"github.com/{{ .Org }}/{{ .Name }}/internal/{{ .ServiceName }}/server"
	"github.com/{{ .Org }}/{{ .Name }}/internal/pkg/logging"
	"github.com/{{ .Org }}/{{ .Name }}/internal/pkg/metrics"
	"github.com/{{ .Org }}/{{ .Name }}/internal/pkg/profiling"
	"github.com/{{ .Org }}/{{ .Name }}/internal/pkg/tracing"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_opentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var (
	Version = "dev"
	App     = "{{ .ServiceName }}"
)

func main() {
	logger, tracer := setup()

	opts := getGRPCServerOpts(logger, tracer)
	grpcServer := getGRPCServer(logger, opts)

	metrics.StartServer()

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	defer cancel()

	mux := runtime.NewServeMux(runtime.WithIncomingHeaderMatcher(headerMatcher))
	httpServerEndpoint := fmt.Sprintf("0.0.0.0:%s", config.Get().Port)
	dialOpts := []grpc.DialOption{grpc.WithInsecure()}

	err := server.RegisterHandlers(ctx, mux, httpServerEndpoint, dialOpts)
	if err != nil {
		logger.Fatal(err.Error())
	}

	srv := &http.Server{
		Addr:         httpServerEndpoint,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 30,

		// http.Handler that delegates to grpcServer on incoming gRPC connections or mux otherwise.
		// Copied from https://github.com/philips/grpc-gateway-example/blob/master/cmd/serve.go#L49-L61
		Handler: h2c.NewHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
				grpcServer.ServeHTTP(w, r)
			} else {
				mux.ServeHTTP(w, r)
			}
		}), &http2.Server{}),
	}

	logger.Info(fmt.Sprintf("starting grpc/http up on %s", httpServerEndpoint))

	conn, err := net.Listen("tcp", httpServerEndpoint)
	if err != nil {
		logger.Fatal(err.Error())
	}

	if err := srv.Serve(conn); err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Fatal(err.Error())
	}
}

func setup() (*zap.Logger, opentracing.Tracer) {
	logger, err := logging.NewLogger(os.Getenv("ENV"))
	if err != nil {
		log.Fatalf("Cannot set up logger: %s", err.Error())
	}

	logging.SetLogger(logger)

	config.Initialise()

	if config.Get().ProfilerEnabled {
		err = profiling.NewProfiler(App, Version)
		if err != nil {
			logger.Error(err.Error())
		}
	}

	tracer, err := tracing.NewTracer(App)
	if err != nil {
		logger.Fatal(err.Error())
	}

	tracing.SetTracer(tracer)

	return logger, tracer
}

func getGRPCServerOpts(logger *zap.Logger, tracer opentracing.Tracer) []grpc.ServerOption {
	opts := []grpc.ServerOption{}

	loggerOpts := []grpc_zap.Option{
		grpc_zap.WithLevels(logging.GrpcCodeToLevel),
		grpc_zap.WithDecider(logging.GetHealthMutedDecider(config.Get().MuteHealthLogging)),
	}

	opts = append(opts, grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
		grpc_opentracing.UnaryServerInterceptor(grpc_opentracing.WithTracer(tracer)),
		grpc_prometheus.UnaryServerInterceptor,
		grpc_zap.UnaryServerInterceptor(logger, loggerOpts...),
	)))

	return opts
}

func getGRPCServer(logger *zap.Logger, opts []grpc.ServerOption) *grpc.Server {
	grpcServer := grpc.NewServer(opts...)
	server.RegisterServices(grpcServer)

	// Register reflection service on gRPC server
	reflection.Register(grpcServer)

	// Enable latency histograms
	grpc_prometheus.EnableHandlingTimeHistogram()

	// Register prometheus metrics on gRPC
	grpc_prometheus.Register(grpcServer)

	// Switch gRPC logger for Zap
	grpc_zap.ReplaceGrpcLoggerV2(logger.WithOptions(zap.IncreaseLevel(zap.WarnLevel)))

	return grpcServer
}

func headerMatcher(key string) (string, bool) {
	switch key { // nolint:gocritic
	default:
		return runtime.DefaultHeaderMatcher(key)
	}
}
