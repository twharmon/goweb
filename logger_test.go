package goweb_test

import (
	"testing"

	"github.com/twharmon/goweb"
)

func TestLogLevelDebug(t *testing.T) {
	equals(t, goweb.LogLevelDebug.String(), "Debug")
}

func TestLogLevelInfo(t *testing.T) {
	equals(t, goweb.LogLevelInfo.String(), "Info")
}

func TestLogLevelNotice(t *testing.T) {
	equals(t, goweb.LogLevelNotice.String(), "Notice")
}

func TestLogLevelWarning(t *testing.T) {
	equals(t, goweb.LogLevelWarning.String(), "Warning")
}

func TestLogLevelError(t *testing.T) {
	equals(t, goweb.LogLevelError.String(), "Error")
}

func TestLogLevelCritical(t *testing.T) {
	equals(t, goweb.LogLevelCritical.String(), "Critical")
}

func TestLogLevelAlert(t *testing.T) {
	equals(t, goweb.LogLevelAlert.String(), "Alert")
}

func TestLogLevelEmergency(t *testing.T) {
	equals(t, goweb.LogLevelEmergency.String(), "Emergency")
}
