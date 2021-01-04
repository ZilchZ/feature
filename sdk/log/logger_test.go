package log

import (
	"testing"
)

func TestLog(t *testing.T) {
	l := LoggerConfig{
		ErrorLogPath:  "/var/log/feature/error.log",
		AccessLogPath: "/var/log/feature/error.log",
		Level:         "DEBUG",
		DevModel:      true,
		StdErr:        true,
		RotateEnable:  false,
	}

	Init(l)
	ApiLogger.Info("xx")
	Logger.Error("xxx")
}
