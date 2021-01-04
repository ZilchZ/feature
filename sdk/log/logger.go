package log

import (
	"fmt"
	"github.com/kjk/dailyrotate"
	"github.com/op/go-logging"
	"io"
	"os"
	"path/filepath"
)

var (
	ApiLogger *logging.Logger
	Logger    *logging.Logger
	logLevel  = map[string]logging.Level{
		"DEBUG":    logging.DEBUG,
		"INFO":     logging.INFO,
		"WARNING":  logging.WARNING,
		"ERROR":    logging.ERROR,
		"CRITICAL": logging.CRITICAL,
		"NOTICE":   logging.NOTICE,
	}
)

func Init(cfg LoggerConfig) {
	ApiLogger = cfg.getLogger("access")
	Logger = cfg.getLogger("error")
}

type LoggerConfig struct {
	ErrorLogPath  string
	AccessLogPath string
	Level         string
	DevModel      bool
	StdErr        bool
	RotateEnable  bool
}

func (cfg *LoggerConfig) getLevelBackend(logFile io.Writer, level string, format logging.Formatter) logging.LeveledBackend {
	logBackend := logging.NewLogBackend(logFile, "", 0)
	levelBackend := logging.AddModuleLevel(logging.NewBackendFormatter(logBackend, format))
	levelBackend.SetLevel(logLevel[level], "")
	return levelBackend

}

func (cfg *LoggerConfig) getLoggerBackend(logPath, level string, stderr bool) logging.LeveledBackend {
	loggerFmtOut := logging.MustStringFormatter(
		`%{color}%{time:2006-01-02 15:04:05} %{shortfile} %{shortfunc} [%{level}%{color:reset}] > %{message}`)
	loggerFmtFile := logging.MustStringFormatter(
		`%{time:2006-01-02 15:04:05} %{shortfile} %{shortfunc} [%{level:.4s}] > %{message}`)
	var err error
	var logWriter io.Writer
	if !cfg.RotateEnable {
		dir := filepath.Dir(logPath)
		err = os.MkdirAll(dir, 0755)
		if err == nil {
			var file *os.File
			file, err = os.OpenFile(logPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
			if err == nil {
				if _, err = file.Seek(0, io.SeekEnd); err == nil {
					logWriter = file
				} else {
					fmt.Printf("seek file: %s to end err: %s \n", logPath, err)
				}
			} else {
				fmt.Printf("open file: %s err: %s \n", logPath, err)
			}
		} else {
			fmt.Printf("mkdir err: %s \n", err)
		}

	} else {
		logWriter, err = dailyrotate.NewFile(logPath, nil)
		if err != nil {
			fmt.Printf("new file: %s err: %s \n", logPath, err)
		}
	}

	loggerBackend := logging.SetBackend()
	fileLevelBackend := cfg.getLevelBackend(logWriter, level, loggerFmtFile)
	stdErrLevelBackend := cfg.getLevelBackend(os.Stderr, level, loggerFmtOut)
	if err == nil && stderr {
		loggerBackend = logging.SetBackend(fileLevelBackend, stdErrLevelBackend)
	} else if err == nil && !stderr {
		loggerBackend = logging.SetBackend(fileLevelBackend)
	} else {
		loggerBackend = logging.SetBackend(stdErrLevelBackend)
	}
	return loggerBackend
}

func (cfg *LoggerConfig) getLogger(logger string) *logging.Logger {
	switch logger {
	case "access":
		var logger = logging.MustGetLogger("access")
		loggerBackend := cfg.getLoggerBackend(cfg.AccessLogPath, cfg.Level, cfg.StdErr || cfg.DevModel)
		logger.SetBackend(loggerBackend)
		return logger
	default:
		var logger = logging.MustGetLogger("error")
		loggerBackend := cfg.getLoggerBackend(cfg.ErrorLogPath, cfg.Level, cfg.StdErr || cfg.DevModel)
		logger.SetBackend(loggerBackend)
		return logger
	}
}
