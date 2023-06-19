package tonic_test

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"testing"

	"bitbucket.org/idomdavis/tonic"
	"bitbucket.org/idomdavis/tonic/config"
	"github.com/stretchr/testify/assert"
)

const ping = "%s://localhost:%d/ping"

func TestListen(t *testing.T) {
	t.Run("Listen works for HTTP", func(t *testing.T) {
		t.Parallel()

		port := Port()
		settings := config.Server{Port: port, Development: true}
		r, err := tonic.New(settings)

		assert.NoError(t, err)

		s, _ := tonic.Listen(settings, r)

		req, err := http.NewRequest(http.MethodGet, fmt.Sprintf(ping, "http", port), nil)

		assert.NoError(t, err)

		res, err := http.DefaultClient.Do(req)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, res.StatusCode)

		_ = s.Close()
	})

	t.Run("Listen works for HTTPS", func(t *testing.T) {
		t.Parallel()

		port := Port()
		settings := config.Server{
			Port:        port,
			Key:         "testdata/localhost.key",
			Certificate: "testdata/localhost.crt",
			Development: true,
		}

		r, err := tonic.New(settings)

		assert.NoError(t, err)

		s, _ := tonic.Listen(settings, r)

		req, err := http.NewRequest(http.MethodGet, fmt.Sprintf(ping, "https", port), nil)

		assert.NoError(t, err)

		tr := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
		client := &http.Client{Transport: tr}
		res, err := client.Do(req)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, res.StatusCode)

		_ = s.Close()
	})
}

// Port returns a usable, free port from the OS. Port will return 0 on error.
func Port() int {
	listener, _ := net.Listen("tcp", "localhost:0")

	defer func() { _ = listener.Close() }()

	//nolint:forcetypeassert // We know it's one of these.
	return listener.Addr().(*net.TCPAddr).Port
}
