package logging

import (
	"os"
	"strings"

	grpcLogging "github.com/grpc-ecosystem/go-grpc-middleware/logging"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc/codes"
)

var (
	codeToLevel = map[codes.Code]zapcore.Level{
		codes.OK:                 zap.DebugLevel,
		codes.Canceled:           zap.DebugLevel,
		codes.Unknown:            zap.ErrorLevel,
		codes.InvalidArgument:    zap.DebugLevel,
		codes.DeadlineExceeded:   zap.WarnLevel,
		codes.NotFound:           zap.DebugLevel,
		codes.AlreadyExists:      zap.DebugLevel,
		codes.PermissionDenied:   zap.WarnLevel,
		codes.Unauthenticated:    zap.DebugLevel, // unauthenticated requests can happen
		codes.ResourceExhausted:  zap.WarnLevel,
		codes.FailedPrecondition: zap.WarnLevel,
		codes.Aborted:            zap.WarnLevel,
		codes.OutOfRange:         zap.WarnLevel,
		codes.Unimplemented:      zap.ErrorLevel,
		codes.Internal:           zap.ErrorLevel,
		codes.Unavailable:        zap.WarnLevel,
		codes.DataLoss:           zap.ErrorLevel,
	}
)

// GrpcCodeToLevel is the default implementation of gRPC return codes and interceptor log level for server side.
func GrpcCodeToLevel(code codes.Code) zapcore.Level {
	val, ok := codeToLevel[code]

	if ok {
		return val
	}

	return zap.ErrorLevel
}

func GetHealthMutedDecider(muted bool) grpcLogging.Decider {
	return func(fullMethodName string, err error) bool {
		if muted && fullMethodName == "/grpc.health.v1.Health/Check" {
			return false
		}

		return true
	}
}

func GetGrpLogDefaultVerbosity() int {
	// By default GRPC will log at warn
	switch strings.ToLower(os.Getenv("GRPC_LOG_LEVEL")) {
	case DebugLevel:
		return int(zap.DebugLevel)
	case InfoLevel:
		return int(zap.InfoLevel)
	case ErrorLevel:
		return int(zap.ErrorLevel)
	case DPanicLevel:
		return int(zap.DPanicLevel)
	case PanicLevel:
		return int(zap.PanicLevel)
	case FatalLevel:
		return int(zap.FatalLevel)
	default:
		return int(zap.WarnLevel)
	}
}
