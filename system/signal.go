package system

import (
	"os"
	"os/signal"
	"syscall"
)

type SignalQuitFunc func()

func InitSignal(fn SignalQuitFunc) {

NEW_SIGNAL:
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL)
	for {
		select {
		case sig := <-ch:
			{
				switch sig {
				case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL:
					if fn != nil {
						fn()
					}
					close(ch)
					return
				case syscall.SIGHUP:
					close(ch)
					goto NEW_SIGNAL
				default:
					close(ch)
					goto NEW_SIGNAL
				}
			}
		}
	}
}
