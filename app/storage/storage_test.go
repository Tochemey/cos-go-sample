package storage

import (
	"context"
	"os"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		// create the text context
		ctx := context.TODO()
		assert.NoError(t, os.Setenv("DB_HOST", testContainer.Host()))
		assert.NoError(t, os.Setenv("DB_PORT", strconv.Itoa(testContainer.Port())))
		assert.NoError(t, os.Setenv("DB_NAME", testDatabase))
		assert.NoError(t, os.Setenv("DB_USER", testUser))
		assert.NoError(t, os.Setenv("DB_PASSWORD", testDatabasePassword))
		assert.NoError(t, os.Setenv("DB_SCHEMA", testContainer.Schema()))

		var dataStore Storage
		fn := func() {
			dataStore = New(ctx)
			assert.NotNil(t, dataStore)
		}

		assert.NotPanics(t, fn)
		assert.NotPanics(t, dataStore.Shutdown(ctx))
		assert.NoError(t, os.Unsetenv("DB_HOST"))
		assert.NoError(t, os.Unsetenv("DB_PORT"))
		assert.NoError(t, os.Unsetenv("DB_NAME"))
		assert.NoError(t, os.Unsetenv("DB_USER"))
		assert.NoError(t, os.Unsetenv("DB_PASSWORD"))
		assert.NoError(t, os.Unsetenv("DB_SCHEMA"))
	})
	t.Run("with env vars parsing failure", func(t *testing.T) {
		// create the text context
		ctx := context.TODO()
		assert.NoError(t, os.Setenv("DB_HOST", testContainer.Host()))
		assert.NoError(t, os.Setenv("DB_PORT", string(rune(testContainer.Port()))))
		assert.NoError(t, os.Setenv("DB_NAME", testDatabase))
		assert.NoError(t, os.Setenv("DB_USER", testUser))
		assert.NoError(t, os.Setenv("DB_PASSWORD", testDatabasePassword))
		assert.NoError(t, os.Setenv("DB_SCHEMA", testContainer.Schema()))

		fn := func() {
			st := New(ctx)
			assert.NotNil(t, st)
		}

		assert.Panics(t, fn)

		assert.NoError(t, os.Unsetenv("DB_HOST"))
		assert.NoError(t, os.Unsetenv("DB_PORT"))
		assert.NoError(t, os.Unsetenv("DB_NAME"))
		assert.NoError(t, os.Unsetenv("DB_USER"))
		assert.NoError(t, os.Unsetenv("DB_PASSWORD"))
		assert.NoError(t, os.Unsetenv("DB_SCHEMA"))
	})
	t.Run("with env vars not set", func(t *testing.T) {
		// create the text context
		ctx := context.TODO()
		assert.NoError(t, os.Setenv("DB_HOST", testContainer.Host()))
		assert.NoError(t, os.Setenv("DB_NAME", testDatabase))
		assert.NoError(t, os.Setenv("DB_USER", testUser))
		assert.NoError(t, os.Setenv("DB_PASSWORD", testDatabasePassword))
		assert.NoError(t, os.Setenv("DB_SCHEMA", testContainer.Schema()))

		fn := func() {
			st := New(ctx)
			assert.NotNil(t, st)
		}

		assert.Panics(t, fn)

		assert.NoError(t, os.Unsetenv("DB_HOST"))
		assert.NoError(t, os.Unsetenv("DB_NAME"))
		assert.NoError(t, os.Unsetenv("DB_USER"))
		assert.NoError(t, os.Unsetenv("DB_PASSWORD"))
		assert.NoError(t, os.Unsetenv("DB_SCHEMA"))
	})
}
