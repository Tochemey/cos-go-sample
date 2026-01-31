package storage

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tochemey/gopack/postgres"
)

func TestConfig(t *testing.T) {
	t.Run("With all environment variables properly set", func(t *testing.T) {
		// set the required env vars
		assert.NoError(t, os.Setenv("DB_HOST", "localhost"))
		assert.NoError(t, os.Setenv("DB_PORT", "5432"))
		assert.NoError(t, os.Setenv("DB_NAME", "postgres"))
		assert.NoError(t, os.Setenv("DB_USER", "postgres"))
		assert.NoError(t, os.Setenv("DB_PASSWORD", "postgres"))
		assert.NoError(t, os.Setenv("DB_SCHEMA", "public"))
		// load the config
		cfg := LoadConfig()
		assert.NotNil(t, cfg)
		assert.Equal(t, "localhost", cfg.DBHost)
		assert.Equal(t, 5432, cfg.DBPort)
		assert.Equal(t, "postgres", cfg.DBUser)
		assert.Equal(t, "postgres", cfg.DBPassword)
		assert.Equal(t, "postgres", cfg.DBName)
		assert.Equal(t, "public", cfg.DBSchema)
		// unset the en vars previously set
		assert.NoError(t, os.Unsetenv("DB_HOST"))
		assert.NoError(t, os.Unsetenv("DB_PORT"))
		assert.NoError(t, os.Unsetenv("DB_NAME"))
		assert.NoError(t, os.Unsetenv("DB_USER"))
		assert.NoError(t, os.Unsetenv("DB_PASSWORD"))
		assert.NoError(t, os.Unsetenv("DB_SCHEMA"))
	})
	t.Run("With the service name is not set", func(t *testing.T) {
		// fetch the actual config. This will panic
		var actual *postgres.Config
		assert.Panics(t, func() {
			actual = LoadConfig()
		})
		assert.Nil(t, actual)
	})
}
