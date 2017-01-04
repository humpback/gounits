package system

import (
	"os"
	"os/signal"
	"syscall"
)

type QuitFunc func()

func InitSignal(fn QuitFunc) {

NEW_SIGNAL:
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL)
	for {
		select {
		case s := <-c:
			{
				switch s {
				case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL:
					if fn != nil {
						fn()
					}
					os.Exit(0) //退出
				case syscall.SIGHUP:
					goto NEW_SIGNAL
				default:
					goto NEW_SIGNAL
				}
			}
		}
	}
}
