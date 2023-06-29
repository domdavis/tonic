package tonic

import (
	"crypto/tls"
	"fmt"
	"net/http"

	"github.com/domdavis/tonic/config"
	"github.com/sirupsen/logrus"
)

// Listen on the given port, using TLS if this is configured in the settings.
// Listen is non-blocking, the returned Done channel can be used to block.
// Call Server.Cancel to cancel Listening.
func Listen(settings config.Server, router http.Handler) (*http.Server, <-chan error) {
	var protocol string

	done := make(chan error)
	server := NewServer(settings, router)

	if settings.TLS() {
		protocol = "HTTPS"

		go func() {
			done <- server.ListenAndServeTLS(settings.Certificate, settings.Key)
		}()
	} else {
		protocol = "HTTP"

		go func() {
			done <- server.ListenAndServe()
		}()
	}

	logrus.WithFields(logrus.Fields{
		"port":     settings.Port,
		"protocol": protocol,
	}).Info("Listening")

	return server, done
}

// NewServer will return a configured http.Server for the given router.
func NewServer(settings config.Server, router http.Handler) *http.Server {
	addressFormat := ":%d"

	if settings.Development {
		addressFormat = "localhost:%d"
	}

	server := &http.Server{
		Handler:      router,
		Addr:         fmt.Sprintf(addressFormat, settings.Port),
		WriteTimeout: settings.WriteTimeout,
		ReadTimeout:  settings.ReadTimeout,
	}

	if settings.TLS() {
		server.TLSConfig = &tls.Config{
			MinVersion: tls.VersionTLS12,
			CurvePreferences: []tls.CurveID{
				tls.CurveP521,
				tls.CurveP384,
				tls.CurveP256},
			CipherSuites: []uint16{
				tls.TLS_AES_128_GCM_SHA256,
				tls.TLS_AES_256_GCM_SHA384,
				tls.TLS_CHACHA20_POLY1305_SHA256,
			},
		}

		server.TLSNextProto = make(
			map[string]func(*http.Server, *tls.Conn, http.Handler))
	}

	return server
}
