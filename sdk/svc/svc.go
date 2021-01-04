package svc

import (
	"os"
	"os/signal"
	"syscall"
)

var (
	signalChan     = make(chan os.Signal, 1)
	signalNotifier map[os.Signal]SignalHandle
)

type Service interface {
	Init() error
	Start() error
	Stop() error
}

type SignalHandle func(sig os.Signal)

func Notify(sig os.Signal, handle SignalHandle) {
	if signalNotifier == nil {
		signalNotifier = make(map[os.Signal]SignalHandle)
	}
	signalNotifier[sig] = handle
}

func Run(s Service) error {
	if err := s.Init(); err != nil {
		return err
	}
	if err := s.Start(); err != nil {
		return err
	}
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGQUIT)
	for {
		sig := <-signalChan
		if handle, ok := signalNotifier[sig]; ok {
			handle(sig)
		} else {
			return s.Stop()
		}
	}
}
