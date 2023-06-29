package tonic_test

import (
	"testing"

	"github.com/domdavis/tonic"
	"github.com/domdavis/tonic/config"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	t.Run("Release mode is set if configured", func(t *testing.T) {
		t.Parallel()

		_, err := tonic.New(config.Server{Development: false})

		assert.NoError(t, err)

		assert.Equal(t, gin.ReleaseMode, gin.Mode())
	})

	t.Run("Invalid static path will error", func(t *testing.T) {
		t.Parallel()

		_, err := tonic.New(config.Server{Static: "./html"})

		assert.Error(t, err)
	})
}
