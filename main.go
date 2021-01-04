package main

import (
	"flag"
	"fmt"
	"github.com/feature/handler"
	"github.com/feature/sdk/log"
	"github.com/feature/sdk/svc"
	"os"
	"runtime"
	"syscall"
)

var (
	flagSet = flag.NewFlagSet("feature", flag.ExitOnError)
	cfgPath = flagSet.String("config", "./conf/conf.toml", "Path of Config files")
	version = flagSet.Bool("version", false, "show relate version info")
)

var (
	Version  string
	CommitId string
	Built    string
)

type program struct {
	worker *handler.Worker
}

func (p *program) Init() error {
	_ = flagSet.Parse(os.Args[1:])
	if *version {
		fmt.Printf("commit id:%s\n", CommitId)
		fmt.Printf("built by %s %s/%s at %s\n", runtime.Version(), runtime.GOOS, runtime.GOARCH, Built)
		os.Exit(2)
	}
	daemon := &handler.Worker{
		Config: *cfgPath,
		Server: &handler.Server{},
	}
	if err := daemon.Init(); err != nil {
		return err
	}
	p.worker = daemon
	return nil
}

func (p *program) Start() error {
	log.Logger.Infof("starting")
	if p.worker != nil {
		p.worker.Main()
	} else {
		log.Logger.Warning("worker was nil")
	}
	return nil
}

func (p *program) Stop() error {
	log.Logger.Warning("stopping")
	if p.worker != nil {
		if err := p.worker.Exit(); err != nil {
			return err
		} else {
			log.Logger.Warning("worker was nil")
		}
	}
	return nil
}

func (p *program) Reload(signal os.Signal) {
	log.Logger.Infof("got signal:%s", signal.String())
	switch signal {
	case syscall.SIGHUP:
		p.worker.Reload()

	}
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	pg := &program{}
	svc.Notify(syscall.SIGHUP, pg.Reload)
	if err := svc.Run(pg); err != nil {
		fmt.Printf("run exit with err:%s", err)
		os.Exit(2)
	} else {
		log.Logger.Info("bye!")
	}
}
