package service

import (
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/tochemey/cos-go-sample/app/grpconfig"
)

func TestGetConfig(t *testing.T) {
	t.Run("With all environment variables properly set", func(t *testing.T) {
		// set the service name and make use of the default env vars values
		assert.NoError(t, os.Setenv("SERVICE_NAME", "accounts"))
		assert.NoError(t, os.Setenv("COS_HOST", "localhost"))
		assert.NoError(t, os.Setenv("COS_PORT", "9000"))

		// let us defined the expected value
		expected := &Config{
			CosHost: "localhost",
			CosPort: 9000,
			GRPCConfig: grpconfig.Config{
				ServiceName:      "accounts",
				GrpcPort:         50051,
				TraceEnabled:     false,
				TraceURL:         "",
				EnableReflection: false,
				MetricsEnabled:   false,
				MetricsPort:      9102,
			},
		}

		// fetch the actual config
		actual := LoadConfig()
		require.NotNil(t, actual)
		assert.True(t, cmp.Equal(expected, actual))
		// free resources
		assert.NoError(t, os.Unsetenv("SERVICE_NAME"))
		assert.NoError(t, os.Unsetenv("COS_HOST"))
		assert.NoError(t, os.Unsetenv("COS_PORT"))
	})

	t.Run("With environment variables not set", func(t *testing.T) {
		// fetch the actual config. This will panic
		var actual *Config
		assert.Panics(t, func() {
			actual = LoadConfig()
		})
		assert.Nil(t, actual)
	})
}
