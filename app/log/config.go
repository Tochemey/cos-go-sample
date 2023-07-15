package log

import (
	"github.com/caarlos0/env/v9"
	"github.com/pkg/errors"
)

// Config represents the gRPC service configuration
type Config struct {
	LogLevel string `env:"LOG_LEVEL" envDefault:"INFO"` // LogLevel define the application log level
}

// loadConfig return the Config from env vars or panic in case of error
// we panic here because this call is usually and must be done during application start
func loadConfig() *Config {
	config := &Config{}
	// all env vars are required
	opts := env.Options{RequiredIfNoDef: true}

	// parse the environment variables and panic in case of error
	// we panic here because this call is usually and must be done during application start
	if err := env.ParseWithOptions(config, opts); err != nil {
		panic(errors.Wrap(err, "unable to load environment variables"))
	}

	return config
}
