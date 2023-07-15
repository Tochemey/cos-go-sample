package grpconfig

import (
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetGrpcConfig(t *testing.T) {
	t.Run("With all environment variables properly set", func(t *testing.T) {
		// set the service name and make use of the default env vars values
		assert.NoError(t, os.Setenv("SERVICE_NAME", "accounts"))

		// let us defined the expected value
		expected := &Config{
			ServiceName:      "accounts",
			GrpcPort:         50051,
			TraceEnabled:     false,
			TraceURL:         "",
			EnableReflection: false,
			MetricsEnabled:   false,
			MetricsPort:      9102,
		}

		// fetch the actual config
		actual := LoadConfig()
		require.NotNil(t, actual)
		assert.True(t, cmp.Equal(expected, actual))
		// free resources
		assert.NoError(t, os.Unsetenv("SERVICE_NAME"))
	})
	t.Run("With the service name is not set", func(t *testing.T) {
		// fetch the actual config. This will panic
		var actual *Config
		assert.Panics(t, func() {
			actual = LoadConfig()
		})
		assert.Nil(t, actual)
	})
}
