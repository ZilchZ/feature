package conf

import (
	"github.com/BurntSushi/toml"
	"github.com/feature/sdk/dal"
	"github.com/feature/sdk/log"
)

var Config = new(Conf)

type Conf struct {
	Title     string
	Server    Server              `toml:"Server"`
	Log       log.LoggerConfig    `toml:"LogConfig"`
	Database  dal.DatabaseConfig  `toml:"Database"`
	Scheduler CrontabJobBasicInfo `toml:"Scheduler"`
}

type Server struct {
	Host            string
	Port            int
	BaseUUID        string
	Pid             string
	Env             string
	ShutdownTimeout int
}

type CrontabJobBasicInfo struct {
	IsScheduled    bool
	Crontab        string
	RunImmediately bool
	Timeout        int
}

func Init(cfgPath string) error {
	if _, err := toml.DecodeFile(cfgPath, &Config); err != nil {
		return err
	}
	if Config.Server.ShutdownTimeout == 0 {
		Config.Server.ShutdownTimeout = 3
	}
	if Config.Server.Pid == "" {
		Config.Server.Pid = "/var/run/feature.pid"
	}
	return nil
}
