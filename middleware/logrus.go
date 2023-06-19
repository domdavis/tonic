package middleware

import (
	"github.com/sirupsen/logrus"
)

// Logrus Logger type for use with Gin.
type Logrus struct {
	Instance *logrus.Logger
}

// LogrusReporter will provide a Reporter that logs via Logrus.
func LogrusReporter(instance *logrus.Logger) *Reporter {
	return NewReporter(&Logrus{Instance: instance})
}

// Log request using Logrus.
func (l *Logrus) Log(level Level, fields map[string]any, message string) {
	if l.Instance == nil {
		l.Instance = logrus.StandardLogger()
	}

	switch level {
	case InfoLevel:
		l.Instance.WithFields(fields).Info(message)
	case WarnLevel:
		l.Instance.WithFields(fields).Warn(message)
	default:
		l.Instance.WithFields(fields).Error(message)
	}
}
