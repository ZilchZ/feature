package handler

import (
	"github.com/feature/handler/crontab"
	"github.com/feature/sdk/log"
)

type Crontab struct {
}

func (c *Crontab) Start() {
	crontab.AddScheduler()
	crontab.ManCronCtx.CronPointer.Start()
}

func (c *Crontab) Stop() {
	log.Logger.Warning("cron stop now ...")
	crontab.ManCronCtx.CronPointer.Stop()

}
