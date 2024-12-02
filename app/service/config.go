package service

import (
	"github.com/caarlos0/env/v9"
	"github.com/pkg/errors"

	"github.com/tochemey/cos-go-sample/app/grpconfig"
)

// Config defines the application config
type Config struct {
	CosHost    string           `env:"COS_HOST"` // CosHost is used to connect to ChiefOfState
	CosPort    int              `env:"COS_PORT"` // CosPort is used to connect to ChiefOfState
	GRPCConfig grpconfig.Config // GRPCConfig is used to spawn gRPC service
}

// LoadConfig fetches the Config from env vars
func LoadConfig() *Config {
	config := &Config{}
	// all env vars are required
	opts := env.Options{RequiredIfNoDef: true}
	if err := env.ParseWithOptions(config, opts); err != nil {
		panic(errors.Wrap(err, "unable to load environment variables"))
		return nil
	}

	return config
}
