package goweb_test

import (
	"testing"

	"github.com/twharmon/goweb"
)

func TestLogLevelDebug(t *testing.T) {
	equals(t, goweb.LogLevelDebug.String(), "DEBUG")
}

func TestLogLevelInfo(t *testing.T) {
	equals(t, goweb.LogLevelInfo.String(), "INFO")
}

func TestLogLevelNotice(t *testing.T) {
	equals(t, goweb.LogLevelNotice.String(), "NOTICE")
}

func TestLogLevelWarning(t *testing.T) {
	equals(t, goweb.LogLevelWarning.String(), "WARNING")
}

func TestLogLevelError(t *testing.T) {
	equals(t, goweb.LogLevelError.String(), "ERROR")
}

func TestLogLevelCritical(t *testing.T) {
	equals(t, goweb.LogLevelCritical.String(), "CRITICAL")
}

func TestLogLevelAlert(t *testing.T) {
	equals(t, goweb.LogLevelAlert.String(), "ALERT")
}

func TestLogLevelEmergency(t *testing.T) {
	equals(t, goweb.LogLevelEmergency.String(), "EMERGENCY")
}
