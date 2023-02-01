package tonic_test

import (
	"crypto/rand"
	"testing"
	"time"

	"bitbucket.org/idomdavis/tonic"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
)

type ErrorReader struct{}

func (ErrorReader) Read([]byte) (int, error) {
	return 0, assert.AnError
}

func TestGenerateSecret(t *testing.T) {
	t.Run("GenerateSecret returns a non empty string", func(t *testing.T) {
		t.Parallel()

		s := tonic.GenerateSecret(rand.Reader)

		assert.NotEmpty(t, s)
	})

	t.Run("GenerateSecret will panic on error", func(t *testing.T) {
		t.Parallel()

		assert.Panics(t, func() { tonic.GenerateSecret(ErrorReader{}) })
	})
}

func TestSignatory_Sign(t *testing.T) {
	t.Run("A zero TTL will error", func(t *testing.T) {
		t.Parallel()

		s := &tonic.Signatory{}
		_, err := s.Sign(&gin.Context{})

		assert.ErrorIs(t, err, tonic.ErrInvalidTTL)
	})

	t.Run("A negative TTL will error", func(t *testing.T) {
		t.Parallel()

		s := &tonic.Signatory{TTL: -1}
		_, err := s.Sign(&gin.Context{})

		assert.ErrorIs(t, err, tonic.ErrInvalidTTL)
	})

	t.Run("An invalid signing method will error", func(t *testing.T) {
		t.Parallel()

		s := &tonic.Signatory{TTL: time.Second, Method: &jwt.SigningMethodHMAC{}}
		_, err := s.Sign(&gin.Context{})

		assert.ErrorIs(t, err, jwt.ErrHashUnavailable)
	})
}

func TestSignatory_Validate(t *testing.T) {
	t.Run("An error parsing will cause validate to fail", func(t *testing.T) {
		t.Parallel()

		s := &tonic.Signatory{}
		r := s.Validate(&gin.Context{}, "")

		assert.False(t, r)
	})
}
