package grpconfig

import (
	"github.com/caarlos0/env/v9"
	"github.com/pkg/errors"
	"github.com/tochemey/gopack/grpc"
	"github.com/tochemey/gopack/log/zapl"
)

// Config represents the gRPC service configuration
type Config struct {
	ServiceName      string `env:"SERVICE_NAME"`                                 // ServiceName is the name given that will show in the traces
	GrpcPort         int    `env:"GRPC_PORT" envDefault:"50051"`                 // GrpcPort is the gRPC port used to received and handle gRPC requests
	MetricsEnabled   bool   `env:"METRICS_ENABLED" envDefault:"false"`           // MetricsEnabled checks whether metrics should be enabled or not
	MetricsPort      int    `env:"METRICS_PORT" envDefault:"9102"`               // MetricsPort is used to send gRPC server metrics to the prometheus server
	TraceEnabled     bool   `env:"TRACE_ENABLED" envDefault:"false"`             // TraceEnabled checks whether tracing should be enabled or not
	TraceURL         string `env:"TRACE_URL" envDefault:""`                      // TraceURL is the OTLP collector url.
	EnableReflection bool   `env:"SERVER_REFLECTION_ENABLED" envDefault:"false"` // EnableReflection this is useful or local dev testing
}

// GetGrpcConfig returns a grpc config from the config object
func (c Config) GetGrpcConfig() *grpc.Config {
	return &grpc.Config{
		ServiceName:      c.ServiceName,
		GrpcHost:         "",
		GrpcPort:         int32(c.GrpcPort),
		TraceEnabled:     c.TraceEnabled,
		TraceURL:         c.TraceURL,
		EnableReflection: c.EnableReflection,
		MetricsEnabled:   c.MetricsEnabled,
		MetricsPort:      c.MetricsPort,
	}
}

// LoadConfig return the Config from env vars or panic in case of error
// we panic here because this call is usually and must be done during application start
func LoadConfig() *grpc.Config {
	config := &Config{}
	// all env vars are required
	opts := env.Options{RequiredIfNoDef: true}

	// parse the environment variables and panic in case of error
	// we panic here because this call is usually and must be done during application start
	if err := env.ParseWithOptions(config, opts); err != nil {
		zapl.Panic(errors.Wrap(err, "unable to load environment variables"))
	}

	// create the actual grpc config and return the object
	return &grpc.Config{
		ServiceName:      config.ServiceName,
		GrpcHost:         "",
		GrpcPort:         int32(config.GrpcPort),
		TraceEnabled:     config.TraceEnabled,
		TraceURL:         config.TraceURL,
		EnableReflection: config.EnableReflection,
		MetricsEnabled:   config.MetricsEnabled,
		MetricsPort:      config.MetricsPort,
	}
}
