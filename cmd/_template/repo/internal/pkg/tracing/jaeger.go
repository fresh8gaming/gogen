package tracing

import (
	"context"
	"log"

	"github.com/{{ .Org }}/{{ .Name }}/internal/pkg/env"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
)

var (
	globalTracer opentracing.Tracer
)

// GetTracer returns the set global tracer, or fatals if it has not been set up.
func GetTracer() opentracing.Tracer {
	if globalTracer == nil {
		log.Fatal("no tracer instantiated")
	}

	return globalTracer
}

// SetTracer sets the global tracer to the one passed to the function.
func SetTracer(tracer opentracing.Tracer) {
	globalTracer = tracer
}

// NewTracer returns a new tracer with the given service name, returning an error
// if there is an issue with creation.
func NewTracer(serviceName string) (opentracing.Tracer, error) {
	tracerConfig := config.Configuration{
		Sampler: &config.SamplerConfig{
			Type:  env.GetenvString("JAEGER_SAMPLER_TYPE", jaeger.SamplerTypeProbabilistic),
			Param: 0,
		},
		Reporter: &config.ReporterConfig{
			LogSpans: true,
		},
	}

	overrideConfig, err := tracerConfig.FromEnv()
	if err != nil {
		log.Fatal(err.Error())
	}

	overrideConfig.ServiceName = serviceName

	tracer, _, err := overrideConfig.NewTracer()
	if err != nil {
		return tracer, err
	}

	return tracer, err
}

func ChildSpanFromContext(ctx context.Context, name string) opentracing.Span {
	parentSpan := opentracing.SpanFromContext(ctx)

	return parentSpan.Tracer().StartSpan(
		name,
		opentracing.ChildOf(parentSpan.Context()),
	)
}
