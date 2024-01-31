package logging

import (
	"log"
	"os"
	"strings"
	"sync"

	"go.uber.org/zap"
)

var (
	globalLogger *zap.Logger
	once         sync.Once
)

const (
	DebugLevel  = "debug"
	WarnLevel   = "warn"
	ErrorLevel  = "error"
	DPanicLevel = "dpanic"
	PanicLevel  = "panic"
	FatalLevel  = "fatal"
	InfoLevel   = "info"
)

// GetLogger returns the global set logger, or fatals if one is not set.
func GetLogger() *zap.Logger {
	if globalLogger == nil {
		log.Fatal("no logger instantiated")
	}

	return globalLogger
}

// SetLogger sets the global logger to the one passed to the function.
func SetLogger(logger *zap.Logger) {
	globalLogger = logger
}

// NewLogger returns a new Zap logger, changing the config based on the environment
// provided, and returning an error if there is an issue with creation.
func NewLogger(env string) (*zap.Logger, error) {
	var config zap.Config

	switch env {
	case "production":
		config = zap.NewProductionConfig()
	default:
		config = zap.NewDevelopmentConfig()
	}

	// Only change the log level if it's defined
	if os.Getenv("LOG_LEVEL") != "" {
		switch strings.ToLower(os.Getenv("LOG_LEVEL")) {
		case DebugLevel:
			config.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
		case WarnLevel:
			config.Level = zap.NewAtomicLevelAt(zap.WarnLevel)
		case ErrorLevel:
			config.Level = zap.NewAtomicLevelAt(zap.ErrorLevel)
		case DPanicLevel:
			config.Level = zap.NewAtomicLevelAt(zap.DPanicLevel)
		case PanicLevel:
			config.Level = zap.NewAtomicLevelAt(zap.PanicLevel)
		case FatalLevel:
			config.Level = zap.NewAtomicLevelAt(zap.FatalLevel)
		default:
			config.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
		}
	}

	logger, err := config.Build()
	if err != nil {
		return logger, err
	}

	_ = logger.Sync()

	return logger, err
}

func SetupLogger() {
	once.Do(func() {
		logger, err := NewLogger(os.Getenv("ENV"))
		if err != nil {
			log.Fatalf("Cannot set up logger: %s", err.Error())
		}
		SetLogger(logger)
	})
}
