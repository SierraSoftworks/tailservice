package proxy

import (
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestResolveDataDir(t *testing.T) {
	configDir, err := os.UserConfigDir()
	require.NoError(t, err)
	require.NotEmpty(t, configDir)

	t.Run("with DataDir", func(t *testing.T) {
		c := Config{
			Hostname: "test",
			DataDir:  "./config",
		}

		assert.Equal(t, "./config", c.resolveDataDir())
	})

	t.Run("without DataDir", func(t *testing.T) {
		c := Config{
			Hostname: "test",
		}

		assert.Equal(t, path.Join(configDir, "tailservice", "test"), c.resolveDataDir())
	})
}
