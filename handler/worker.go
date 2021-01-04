package handler

import (
	"github.com/feature/conf"
	"github.com/feature/handler/dao"
	"github.com/feature/sdk/log"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
)

type Worker struct {
	Config  string
	HasPid  bool
	Server  *Server
	Crontab *Crontab
}

func (w *Worker) Init() error {
	if err := conf.Init(w.Config); err != nil {
		return err
	}
	log.Init(conf.Config.Log)
	dao.Init()
	w.Server.Init()
	if !w.HasPid {
		dir := filepath.Dir(conf.Config.Server.Pid)
		if err := os.MkdirAll(dir, 755); err != nil {
			return err
		}
		pid := strconv.Itoa(os.Getpid())
		if err := ioutil.WriteFile(conf.Config.Server.Pid, []byte(pid), 644); err != nil {
			return err
		}
		w.HasPid = true
	}
	return nil
}

func (w *Worker) Reload() {
	log.Logger.Infof("reload with config:%s", w.Config)
	if err := w.Init(); err != nil {
		log.Logger.Error(err)
	}
	log.Logger.Info("reload config to end")

}

func (w *Worker) Main() {
	w.Crontab.Start()
	w.Server.Launch()

}

func (w *Worker) Exit() error {
	w.Crontab.Stop()
	if err := w.Server.Stop(); err != nil {
		return nil
	}
	if err := os.Remove(conf.Config.Server.Pid); err != nil {
		log.Logger.Error(err)
		return err
	}
	return nil
}
