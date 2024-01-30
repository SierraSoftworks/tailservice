package proxy

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseListener(t *testing.T) {
	t.Run("local", func(t *testing.T) {
		l, err := ParseListener("8080:80", "tcp", false)
		require.NoError(t, err)

		assert.Equal(t, "tcp", l.Proto)
		assert.Equal(t, 8080, l.Port)
		assert.Equal(t, false, l.Secure)
		assert.Equal(t, "127.0.0.1:80", l.Target)
	})

	t.Run("remote", func(t *testing.T) {
		l, err := ParseListener("8080:example.com:80", "tcp", false)
		require.NoError(t, err)

		assert.Equal(t, "tcp", l.Proto)
		assert.Equal(t, 8080, l.Port)
		assert.Equal(t, false, l.Secure)
		assert.Equal(t, "example.com:80", l.Target)
	})
}
