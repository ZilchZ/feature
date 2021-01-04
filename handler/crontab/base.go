package crontab

import (
	"fmt"
	"github.com/feature/conf"
	"github.com/feature/sdk/log"
	"github.com/robfig/cron/v3"
	"time"
)

var (
	ManCronCtx = &CronContext{
		CronPointer: cron.New(
			cron.WithSeconds(),
			cron.WithChain(
				cron.SkipIfStillRunning(cron.DefaultLogger),
				cron.Recover(cron.DefaultLogger)),
		)}
)

type CronContext struct {
	CronPointer *cron.Cron
}

func AddScheduler() {
	if err := SchedulerDemo(); err != nil {
		log.Logger.Error(err)
	}
}

func SchedulerDemo() error {
	config := conf.Config.Scheduler
	if !config.IsScheduled {
		return nil
	}
	if config.RunImmediately {
		go TestDemo()
	}
	if _, err := ManCronCtx.CronPointer.AddFunc(config.Crontab, TestDemo); err != nil {
		log.Logger.Error(err)
		return err
	}
	return nil
}

func TestDemo() {
	fmt.Println("hello world!\n", time.Now())
}
