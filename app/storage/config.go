package storage

import (
	"time"

	"github.com/caarlos0/env/v9"
	"github.com/pkg/errors"
	"github.com/tochemey/gopack/log/zapl"
	"github.com/tochemey/gopack/postgres"
)

// Config holds the storage configuration
type Config struct {
	DBHost                string        `env:"DB_HOST"`                                   // DBHost represents the database host
	DBPort                int           `env:"DB_PORT"`                                   // DBPort is the database port
	DBName                string        `env:"DB_NAME"`                                   // DBName is the database name
	DBUser                string        `env:"DB_USER"`                                   // DBUser is the database user used to connect
	DBPassword            string        `env:"DB_PASSWORD"`                               // DBPassword is the database password
	DBSchema              string        `env:"DB_SCHEMA"`                                 // DBSchema represents the database schema
	MaxOpenConnections    int           `env:"MAX_OPEN_CONNECTIONS" envDefault:"25"`      // MaxOpenConnections represents the number of open connections in the pool
	MaxIdleConnections    int           `env:"MAX_IDLE_CONNECTIONS" envDefault:"25"`      // MaxIdleConnections represents the number of idle connections in the pool
	ConnectionMaxLifetime time.Duration `env:"CONNECTION_MAX_LIFETIME" envDefault:"5m0s"` // ConnectionMaxLifetime represents the connection max life time
}

// LoadConfig read the Postgres config from environment variables
func LoadConfig() *postgres.Config {
	config := &Config{}
	// all env vars are required
	opts := env.Options{RequiredIfNoDef: true}

	// parse the env vars into the configuration
	if err := env.ParseWithOptions(config, opts); err != nil {
		zapl.Panic(errors.Wrap(err, "unable to load environment variables"))
	}

	// return the postgres config
	return &postgres.Config{
		DBHost:                config.DBHost,
		DBPort:                config.DBPort,
		DBName:                config.DBName,
		DBUser:                config.DBUser,
		DBPassword:            config.DBPassword,
		DBSchema:              config.DBSchema,
		MaxConnections:        config.MaxOpenConnections,
		MaxConnectionLifetime: config.ConnectionMaxLifetime,
	}
}
